package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	edl "github.com/sndnvaps/ebookdownloader"
	"github.com/sndnvaps/ebookdownloader/http-server/middleware"
	"gopkg.in/urfave/cli.v1"
)

var (
	//Version 版本信息
	Version string = "dev"
	//Commit git commit信息
	Commit string = "b40f73c79"
	//BuildTime 编译时间
	BuildTime string = "2020-02-16 16:34"
)

var (
	bookinfo      edl.BookInfo         //初始化变量
	ebdlInterface edl.EBookDLInterface //初始化接口
)

var (
	lock sync.Mutex
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

var identityKey = "id"

//HTTPStat 系统信息
func HTTPStat(c *gin.Context) {
	c.JSON(200, gin.H{
		"ebookdownloader_Version": Version,
		"HashCommit":              Commit,
		"SystemBuildTime":         BuildTime,
		"hostinfo":                conf,
	})
	c.String(http.StatusOK, "ok")
}

//CheckUpdate 检查是否需要更新
func CheckUpdate(c *gin.Context) {
	result, err := edl.UpdateCheck()
	if err == nil {
		compareResult := result.Compare(Version)
		c.JSON(200, gin.H{
			"CurrentVersion": Version,
			"update_check":   compareResult,
		})
		c.String(http.StatusOK, "ok")
		return
	}
}

// ParseEbhostAndBookIDPost 处理下载小说的请求
func ParseEbhostAndBookIDPost(c *gin.Context) {

	bookid := c.Query("bookid")
	ebhost := c.DefaultQuery("ebhost", "xsbiquge.com") //设置默认值为 xsbiquge.com

	isTxtStr := c.DefaultQuery("istxt", "false")   //需要传入bool值 , 0,1,true,false
	isMobiStr := c.DefaultQuery("ismobi", "false") //需要传入bool值, 0,1,true,false
	isEpubStr := c.DefaultQuery("isepub", "false") //需要传入bool值，0,1,true,false

	var metainfo edl.Meta //用于保存小说的meta信息
	var cmdArgs []string  //定义命令用到的参数

	isTxt, errTxt := strconv.ParseBool(isTxtStr)
	if errTxt != nil {
		isTxt = false
	}
	isMobi, errMobi := strconv.ParseBool(isMobiStr)
	if errMobi != nil {
		isMobi = false
	}

	isEpub, errEpub := strconv.ParseBool(isEpubStr)
	if errEpub != nil {
		isEpub = false
	}

	//当 bookid没有设置的时候，返回错误
	if bookid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bookid没有设置"})
		return
	}
	switch ebhost {
	case "xsbiquge.com":
		cmdArgs = append(cmdArgs, "--ebhost=xsbiquge.com")
		xsbiquge := edl.NewXSBiquge()
		ebdlInterface = xsbiquge //实例化接口
	case "999xs.com":
		cmdArgs = append(cmdArgs, "--ebhost=999xs.com")
		xs999 := edl.New999XS()
		ebdlInterface = xs999 //实例化接口
	case "23us.la":
		cmdArgs = append(cmdArgs, "--ebhost=23us.la")
		xs23 := edl.New23US()
		ebdlInterface = xs23 //实例化接口
	}

	//add --bookid={{.bookid}}
	cmdArgs = append(cmdArgs, fmt.Sprintf("--bookid=%s", bookid))

	bookinfo = ebdlInterface.GetBookBriefInfo(bookid, "")

	if isTxt {
		cmdArgs = append(cmdArgs, "--txt")
	}
	if isMobi {
		cmdArgs = append(cmdArgs, "--mobi")
	}

	if isEpub {
		cmdArgs = append(cmdArgs, "--epub")
	}

	//添加生成meta.json参数
	cmdArgs = append(cmdArgs, "--meta")

	cmd := EbookdownloaderCliCmd(cmdArgs...)
	lock.Lock()
	err := cmd.Run()
	lock.Unlock()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  "下载小说失败",
		})
	}
	metainfo, _ = edl.GetMetaData("./outputs/" + bookinfo.Name + "-" + bookinfo.Author + "/meta.json")

	c.JSON(http.StatusOK, gin.H{
		"isTxt":    isTxtStr,
		"isMobi":   isMobiStr,
		"metainfo": metainfo,
	})

}

//ebook_http_server 启动ebookdownloader服务器后台程序
func ebookHTTPServer(c *cli.Context) error {

	//从参数中获取配置文件的路径
	CFGPATH = c.String("conf")
	//初始化配置文件
	ConfInit()

	// Creates a router with Default
	router := gin.Default()

	//使用中间件，处理跨域问题
	router.Use(middleware.Cors())

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &User{
					UserName:  userID,
					LastName:  "Jimes",
					FirstName: "Yang",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.UserName == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	//登陆，并生成 token
	router.POST("/login", authMiddleware.LoginHandler)
	//退出，并删除cookie中的 token
	router.GET("/logout", authMiddleware.LogoutHandler)

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := router.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// $ curl -X GET -v --form istxt=true --form ismobi=false "http://localhost:8080/auth/post?ebhost=23us.la&bookid=0_062&istxt=true&ismobi=true"
		auth.GET("/post", ParseEbhostAndBookIDPost)
		//列举./public目录所有的文件
		auth.GET("/get_list", List)
		//删除 服务器上面已经下载的小说
		// $ curl -X GET "http://localhost:8080/auth/del/我是谁-sndnvaps/我是谁-sndnvaps.mobi"
		// $ curl -X GET "http://localhost:8080/auth/del/我真不是作者菌-sndnvaps/我真不是作者菌-sndnvaps.txt"
		auth.DELETE("/del/:ebpath/:bookname", Del)
		//系统状态信息
		// http://localhost:8080/auth/stat
		auth.GET("/stat", HTTPStat)

		//检查是否需要更新
		// http://localhost:8080/auth/check_update
		auth.GET("/check_update", CheckUpdate)

	}

	apiV1 := router.Group("/api/v1")
	apiV1.Use(authMiddleware.MiddlewareFunc())
	{
		apiV1.POST("/job", EbookDLCreateJob)
		apiV1.GET("/job/*id", EbookDLGetJob)
	}

	//简单文件服务器
	// http://localhost:8080/file
	//public存放着要显示的文件
	router.StaticFS("/public", http.Dir("outputs"))

	router.Run(conf.InerHost + ":" + conf.Port) // 监听并在 0.0.0.0:8080 上启动服务

	return nil

}

func main() {

	app := cli.NewApp()
	app.Name = "golang EBookDownloader http-server"
	app.Version = Version + "-" + Commit + "-" + BuildTime
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jimes Yang",
			Email: "sndnvaps@gmail.com",
		},
	}
	app.Copyright = "© 2019 - 2020 Jimes Yang<sndnvaps@gmail.com>"
	app.Usage = "用于下载 笔趣阁(https://www.xsbiquge.com),999小说网(https://www.999xs.com/) ,顶点小说网(https://www.23us.la) 上面的电子书，并保存为txt格式或者(mobi格式,awz3格式)的电子书"
	app.Action = ebookHTTPServer
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf,c",
			Value: "./conf/ebdl_conf.ini",
			Usage: "定义http-server的配置文件路径",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
