package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//用于显示 public目录所有的文件
func List(c *gin.Context) {

	files, _ := filepath.Glob("./outputs/*")
	//fmt.Println(files) // contains a list of all files in the current directory
	for index := 0; index < len(files); index++ {
		files[index] = "public/" + filepath.Base(files[index])
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}
