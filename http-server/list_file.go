package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type FileList struct {
	Metainfo Meta `json:"metainfo"`
}

//用于显示 public目录所有的文件
func List(c *gin.Context) {

	var filelist []FileList
	var tmp FileList
	var err error

	path := "./outputs/"
	//以只读的方式打开目录
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	//延迟关闭目录
	defer f.Close()
	fileInfo, _ := f.Readdir(-1)

	for _, info := range fileInfo {
		//判断是否是目录,当前只有目录，不存在文件
		if info.IsDir() {
			metapath := path + info.Name() + "/meta.json"

			tmp.Metainfo, err = GetMetaData(metapath)
			if err == nil {
				filelist = append(filelist, tmp)
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"files": filelist})
}
