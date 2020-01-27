package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileList struct {
	FileName string `json:"filename"`
	Url string `json:"url"`
}

//用于显示 public目录所有的文件
func List(c *gin.Context) {
    var filelist []FileList
	files, _ := filepath.Glob("./outputs/*")
	//fmt.Println(files) // contains a list of all files in the current directory
	for index := 0; index < len(files); index++ {
		files[index] = conf.URL_BASE + "/" + "public/" + filepath.Base(files[index])
		tmp := FileList{
			FileName: filepath.Base(files[index]),
			Url: files[index],
		}
		filelist = append(filelist,tmp)
	}

	c.JSON(http.StatusOK, gin.H{"files": filelist})
}
