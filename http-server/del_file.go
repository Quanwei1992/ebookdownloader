package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// Del 用于删除 服务器中outputs/目录指定的文件 :bookname
func Del(c *gin.Context) {

	bookname := c.Param("bookname") //定义为小说名，或者是或者杂项；如果出现的为del,就执行删除目录操作
	ebpath := c.Param("ebpath")     //小说对应目录
	ext := filepath.Ext(bookname)
	//删除目录操作
	if bookname == "del" {
		fullpath := "./outputs/" + ebpath + "/"
		fullpath, _ = filepath.Abs(fullpath)
		if !com.IsExist(fullpath) {
			c.JSON(http.StatusOK, gin.H{"Status": fullpath + " is not exist in the serve"})
			return
		}
		err := os.RemoveAll(fullpath)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Status": ebpath + " has been remove"})
		return
	}

	//删除文件操作
	if ext == ".txt" || ext == ".mobi" ||
		ext == ".azw3" || ext == ".json" ||
		ext == ".jpg" {
		fullpath := "./outputs/" + ebpath + "/" + bookname
		fullpath, _ = filepath.Abs(fullpath)
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

	}

	c.JSON(http.StatusBadRequest, gin.H{"error": bookname + " is not support ext file"})
	return

}
