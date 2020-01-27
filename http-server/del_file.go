package main

import (
	"net/http"
	"path/filepath"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

//用于删除 服务器中outputs/目录指定的文件 :bookname
func Del(c *gin.Context) {
	// gin设置响应头，设置跨域
	c.Header("Access-Control-Allow-Origin", "*")

	bookname := c.Param("bookname")
	ext := filepath.Ext(bookname)
	if ext == ".txt" || ext == ".mobi" || ext == ".azw3" {
	 fullpath := "outputs/" + bookname
	 if !com.IsExist(fullpath) {
		c.JSON(http.StatusOK, gin.H{"Status": fullpath + " is not exist in the serve"})
		return
	 }
	  err := os.Remove(fullpath)
	  if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	  }
	  c.JSON(http.StatusOK, gin.H{"Status": bookname + " has been delete"})
	  return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error":  bookname+" is not support ext file"})
		return
	}
}
