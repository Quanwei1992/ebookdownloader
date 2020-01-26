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
	Version   string = "1.6.9"
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
	c.JSON(200, gin.H{
		"ebookdownloader_Version": Version,
		"HashCommit":              Commit,
		"SystemBuildTime":         BuildTime,
	})
	c.String(http.StatusOK, "ok")
}

//The request responds to a url matching:  ?ebhost=xsbiquge.com&bookid=0_062
func ParseEbhostAndBookIdGet(c *gin.Context) {
	ebhost := c.DefaultQuery("ebhost", "xsbiquge.com") //设置默认值为 xsbiquge.com
	bookid := c.Query("bookid")
	c.JSON(200, gin.H{
		"status": "geted",
		"ebhost": ebhost,
		"bookid": bookid,
	})
}

func ParseEbhostAndBookIdPost(c *gin.Context) {
	bookid := c.Query("bookid")
	ebhost := c.DefaultQuery("ebhost", "xsbiquge.com") //设置默认值为 xsbiquge.com

	isTxtStr := c.PostForm("istxt")   //需要传入bool值 , 0,1,true,false
	isMobiStr := c.PostForm("ismobi") //需要传入bool值, 0,1,true,false

	txtfilepath := ""  //定义 txt下载后，获取得到的 地址
	mobifilepath := "" //定义 txt下载后，获取得到的 地址

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

	author := bookinfo.Author
	description := bookinfo.Description

	bookinfo = EBDLInterface.DownloadChapters(bookinfo, "")

	if isTxt {
		bookinfo.GenerateTxt()
		txtfilepath = "public/" + bookinfo.Name + "-" + bookinfo.Author + ".txt"
	}
	if isMobi {
		bookinfo.SetKindleEbookType(true, false)
		lock.Lock()
		bookinfo.GenerateMobi()
		lock.Unlock()
		mobifilepath = "public/" + bookinfo.Name + "-" + bookinfo.Author + ".mobi"

	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "post",
		"ebhost":       ebhost,
		"bookid":       bookid,
		"isTxt":        isTxtStr,
		"isMobi":       isMobiStr,
		"author":       author,
		"description":  description,
		"txtfilepath":  txtfilepath,
		"mobifilepath": mobifilepath,
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
	filepath := "http://localhost:8080/public" + filename
	c.JSON(http.StatusOK, gin.H{"filepath": filepath})
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//The request responds to a url matching:  /get?ebhost=xsbiquge.com&bookid=0_062
	// //get?ebhost=23us.la&bookid=0_062
	// /get?bookid=0_062
	// $ curl -X GET -v http://localhost:8080/get?ebhost=999xs.com&bookid=0_062"

	//router.GET("/get", ParseEbhostAndBookIdGet)

	// $ curl -X POST -v --form istxt=true --form ismobi=false "http://localhost:8080/post?ebhost=23us.la&bookid=0_062"
	router.POST("/post", ParseEbhostAndBookIdPost)

	// $ curl -X POST --form "file=@./hello.txt" http://localhost:8080/upload
	router.POST("/upload", Upload)

	//列举./public目录所有的文件
	router.GET("/list", List)

	//简单文件服务器
	// http://localhost:8080/file
	//public存放着要显示的文件
	router.StaticFS("/public", http.Dir("outputs"))

	//系统状态信息
	// http://localhost:8080/stat
	router.GET("/stat", HttpStat)

	router.Run(Host + ":" + Port) // 监听并在 0.0.0.0:8080 上启动服务
}
