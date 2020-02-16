package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ajvb/kala/job"
	"github.com/gin-gonic/gin"
	edl "github.com/sndnvaps/ebookdownloader"
)

//EbookDLCreateJob 基于kala Job Schedule创建的下载任务
func EbookDLCreateJob(c *gin.Context) {
	bookid := c.Query("bookid")
	ebhost := c.DefaultQuery("ebhost", "xsbiquge.com") //设置默认值为 xsbiquge.com

	isTxtStr := c.DefaultQuery("istxt", "false")   //需要传入bool值 , 0,1,true,false
	isMobiStr := c.DefaultQuery("ismobi", "false") //需要传入bool值, 0,1,true,false
	var cmdArgs []string

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

	bookinfo = ebdlInterface.GetBookInfo(bookid, "")

	if isTxt {
		cmdArgs = append(cmdArgs, "--txt")
	}
	if isMobi {
		cmdArgs = append(cmdArgs, "--mobi")
	}

	//添加生成meta.json参数
	cmdArgs = append(cmdArgs, "--meta")

	cmdArgsStr := strings.Join(cmdArgs, " ")

	//构建kala 运行命令
	cmd := ebdBinPathConf.Path + " " + cmdArgsStr

	//因为在windows系统测试，需要做一些替换
	cmd = strings.Replace(cmd, "\\", "/", -1)

	//Build kala job schedule info
	// our job just run once,after 5 minute
	schedule := fmt.Sprintf("R0/%s/", time.Now().Add(time.Minute*10).Format(time.RFC3339))
	jobName := fmt.Sprintf("Downloader ebook %s-%s", bookinfo.Name, bookinfo.Author)

	downloadJobBody := BuildKalaJobInfo(jobName, cmd, schedule)
	kalaClient := NewKalaClient()
	jobID, err := kalaClient.CreateJob(downloadJobBody)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   201,
			"msg":    "创建下载下载失败",
			"errMsg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{

		"code":        200,
		"msg":         "创建下载任务完成",
		"JobID":       jobID,
		"JobName":     jobName,
		"JobCMD":      cmd,
		"JobSchedule": schedule,
	})
}

//JobInfoMeta 用于显示相应的job信息
type JobInfoMeta struct {
	JobID string   `json:"id"`
	Job   *job.Job `json:"job"`
}

//EbookDLGetJob 列举kala Job Schedule中有的任务
func EbookDLGetJob(c *gin.Context) {
	jobID := c.Param("id")
	var jobInfoMetas []JobInfoMeta

	//把 '/' 替换成 空白符
	jobID = strings.Replace(jobID, "/", "", -1)

	kc := NewKalaClient()
	if jobID == "" {
		jobInfos, _ := kc.GetAllJobs()
		for k, v := range jobInfos {
			tmp := JobInfoMeta{
				JobID: k,
				Job:   v,
			}
			jobInfoMetas = append(jobInfoMetas, tmp)

		}

		c.JSON(http.StatusOK, gin.H{
			"jobinfos": jobInfoMetas,
		})
		return
	}
	if kc.VerifyJobID(jobID) {
		jobInfo, _ := kc.GetJob(jobID)
		c.JSON(http.StatusOK, gin.H{
			"jobinfo": jobInfo,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobid": jobID,
		"msg":   "输入的jobid无法识别",
	})
	return

}
