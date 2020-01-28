package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

//用于删除 服务器中outputs/目录指定的文件 :bookname
func Del(c *gin.Context) {

	bookname := c.Param("bookname") //定义为小说名，或者是或者杂项；如果出现的为del,就执行删除目录操作
	ebpath := c.Param("ebpath")     //小说对应目录
	ext := filepath.Ext(bookname)
	//删除目录操作
	if bookname == "del" {
		fullpath := "outputs/" + ebpath + "/"
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
		fullpath := "outputs/" + ebpath + "/" + bookname
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
		c.JSON(http.StatusBadRequest, gin.H{"error": bookname + " is not support ext file"})
		return
	}
}

//用于删除 服务器中outputs/目录指定的目录 :ebpath, *action == del时删除，其它的时候，选择无视
func DelFolder(c *gin.Context) {

	ebpath := c.Param("ebpath") //小说对应目录
	action := c.Param("action") //要执行的操作，目前只有指定 del
	if action == "del" {

		fullpath := "outputs/" + ebpath
		if !com.IsExist(fullpath) {
			c.JSON(http.StatusOK, gin.H{"Status": fullpath + " is not exist in the serve"})
			return
		}
		err := os.Remove(fullpath)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Status": ebpath + " has been delete"})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "actionName '" + action + "' is not support"})
		return
	}
}
