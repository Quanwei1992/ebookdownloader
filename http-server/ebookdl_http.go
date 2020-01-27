package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	edl "github.com/sndnvaps/ebookdownloader"
)

var (
	Version   string = "1.7.0"
	Commit    string = "b40f73c79"
	BuildTime string = "2020-01-25 16:19"
)

var (
	bookinfo      edl.BookInfo         //初始化变量
	EBDLInterface edl.EBookDLInterface //初始化接口
)

var (
	lock sync.Mutex
)

//系统信息
func HttpStat(c *gin.Context) {
	// gin设置响应头，设置跨域
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, gin.H{
		"ebookdownloader_Version": Version,
		"HashCommit":              Commit,
		"SystemBuildTime":         BuildTime,
		"hostinfo":                conf,
	})
	c.String(http.StatusOK, "ok")
}

func ParseEbhostAndBookIdPost(c *gin.Context) {

	bookid := c.Query("bookid")
	ebhost := c.DefaultQuery("ebhost", "xsbiquge.com") //设置默认值为 xsbiquge.com

	isTxtStr := c.DefaultQuery("istxt", "false")   //需要传入bool值 , 0,1,true,false
	isMobiStr := c.DefaultQuery("ismobi", "false") //需要传入bool值, 0,1,true,false

	txtfilepath := ""  //定义 txt下载后，获取得到的 地址
	mobifilepath := "" //定义 mobi下载后，获取得到的 地址
	cover_url_path := "" //定义下载小说后，封面的url地址
	var metainfo Meta  //用于保存小说的meta信息

	isTxt, errTxt := strconv.ParseBool(isTxtStr)
	if errTxt != nil {
		isTxt = false
	}
	isMobi, errMobi := strconv.ParseBool(isMobiStr)
	if errMobi != nil {
		isMobi = false
	}

	//当 bookid没有设置的时候，返回错误
	if bookid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bookid没有设置"})
		return
	}
	switch ebhost {
	case "xsbiquge.com":
		xsbiquge := edl.NewXSBiquge()
		EBDLInterface = xsbiquge //实例化接口
	case "999xs.com":
		xs999 := edl.New999XS()
		EBDLInterface = xs999 //实例化接口
	case "23us.la":
		xs23 := edl.New23US()
		EBDLInterface = xs23 //实例化接口
	}

	bookinfo = EBDLInterface.GetBookInfo(bookid, "")

	bookinfo = EBDLInterface.DownloadChapters(bookinfo, "")

	if isTxt {
		bookinfo.GenerateTxt()
		txtfilepath = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + bookinfo.Name + "-" + bookinfo.Author + ".txt"
	}
	if isMobi {
		bookinfo.SetKindleEbookType(true, false)
		lock.Lock()
		bookinfo.GenerateMobi()
		lock.Unlock()
		mobifilepath = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + bookinfo.Name + "-" + bookinfo.Author + ".mobi"
		cover_url_path = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + "cover.jpg"

	}

	metainfo = Meta{
		Ebhost:      ebhost,
		Bookid:      bookid,
		BookName:    bookinfo.Name,
		Author:      bookinfo.Author,
		CoverUrl: cover_url_path,
		Description: bookinfo.Description,
		TxtUrlPath:  txtfilepath,
		MobiUrlPath: mobifilepath,
	}

	metainfo.WriteFile("./outputs/" + bookinfo.Name + "-" + bookinfo.Author + "/meta.json")

	c.JSON(http.StatusOK, gin.H{
		"isTxt":    isTxtStr,
		"isMobi":   isMobiStr,
		"metainfo": metainfo,
	})

}

//用于上传文件，并保存到服务器的 public目录里面
func Upload(c *gin.Context) {

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename
	out, err := os.Create("outputs/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	filepath := conf.URL_BASE + "/public/" + filename
	c.JSON(http.StatusOK, gin.H{"filepath": filepath})
}

func main() {
	// Creates a router with Default
	router := gin.Default()

	//使用中间件，处理跨域问题
	router.Use(AccessCROSMiddleware())

	// $ curl -X GET -v --form istxt=true --form ismobi=false "http://localhost:8080/post?ebhost=23us.la&bookid=0_062&istxt=true&ismobi=true"
	router.GET("/post", ParseEbhostAndBookIdPost)

	// $ curl -X POST --form "file=@./hello.txt" http://localhost:8080/upload
	//router.POST("/upload", Upload)

	//列举./public目录所有的文件
	router.GET("/get_list", List)

	//删除 服务器上面已经下载的小说
	// $ curl -X GET "http://localhost:8080/del/我是谁-sndnvaps/我是谁-sndnvaps.mobi"
	// $ curl -X GET "http://localhost:8080/del/我真不是作者菌-sndnvaps/我真不是作者菌-sndnvaps.txt"
	del := router.Group("/del")
	{
		del.GET("/:ebpath/:bookname",Del)
	}

	//简单文件服务器
	// http://localhost:8080/file
	//public存放着要显示的文件
	router.StaticFS("/public", http.Dir("outputs"))

	//系统状态信息
	// http://localhost:8080/stat
	router.GET("/stat", HttpStat)

	router.Run(conf.InerHost + ":" + conf.Port) // 监听并在 0.0.0.0:8080 上启动服务
}
