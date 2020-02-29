package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	edl "github.com/sndnvaps/ebookdownloader"
)

//FileList 文件信息
type FileList struct {
	Metainfo edl.Meta `json:"metainfo"`
}

//List 用于显示 public目录所有的文件
func List(c *gin.Context) {

	var filelist []FileList
	var tmp FileList
	var err error

	path := "./outputs/"
	path, _ = filepath.Abs(path)

	//fmt.Println(path)
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
			metapath := path + string(os.PathSeparator) + info.Name() + "/meta.json"

			tmp.Metainfo, err = edl.GetMetaData(metapath)
			if err == nil {
				if tmp.Metainfo.CoverURL != "" {
					tmp.Metainfo.CoverURL = conf.URLBase + "/" + tmp.Metainfo.CoverURL
				}
				if tmp.Metainfo.TxtURLPath != "" {
					tmp.Metainfo.TxtURLPath = conf.URLBase + "/" + tmp.Metainfo.TxtURLPath
				}
				if tmp.Metainfo.MobiURLPath != "" {
					tmp.Metainfo.MobiURLPath = conf.URLBase + "/" + tmp.Metainfo.MobiURLPath
				}
				filelist = append(filelist, tmp)
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"files": filelist})
}
