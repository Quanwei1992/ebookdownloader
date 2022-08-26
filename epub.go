package ebookdownloader

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Chain-Zhang/pinyin"
	"github.com/sndnvaps/go-epub"
	"github.com/unknwon/com"
)

// GenerateEPUB 生成ebpub小说
func (this BookInfo) GenerateEPUB() error {
	//将文件名转换成拼音
	strPinyin, _ := pinyin.New(this.Name).Split("-").Mode(pinyin.WithoutTone).Convert()
	savepath := "./tmp/" + strPinyin
	savepath, _ = filepath.Abs(savepath) //使用绝对路径
	if com.IsExist(savepath) {
		os.RemoveAll(savepath)
	}
	os.MkdirAll(path.Dir(savepath)+string(os.PathSeparator), os.ModePerm)

	//bookISBN 设置小说的urn码
	bookISBNStr := this.ISBN()

	//设置生成mobi的输出目录
	outputpath := "./outputs/" + this.Name + "-" + this.Author + "/"
	outputpath, _ = filepath.Abs(outputpath)
	outputpath = outputpath + string(os.PathSeparator) //使用绝对路径
	//fmt.Println(outputpath)
	if !com.IsExist(outputpath) {
		os.MkdirAll(outputpath, os.ModePerm)
	}
	// 生成封面
	//GenerateCover(this)
	//下载封面

	err := this.GetCover()
	if err != nil {
		fmt.Println(err.Error())
	}

	//把封面复制到 tmp 目录当中
	coverPath, _ := filepath.Abs("./cover.jpg")

	//把封面复制到 outputs/小说名-作者/cover.jpg
	err = com.Copy(coverPath, outputpath+string(os.PathSeparator)+"cover.jpg")
	if err != nil {
		fmt.Println(err.Error())
	}
	//删除当前目前的cover.jpg文件
	os.RemoveAll(coverPath)

	//创建epub小说信息
	e := epub.NewEpub(this.Name)
	e.SetAuthor(this.Author)
	epubCover, _ := e.AddImage("./outputs/"+this.Name+"-"+this.Author+"/"+"cover.jpg", "cover.jpg")
	//epubCoverCSS, _ := e.AddCSS("./tpls/epub_cover.css", "cover.css")
	e.SetCover(epubCover, "")                  //设置封面,使用默认的 cover.css, 不做自定义
	e.SetDescription(this.Description)         //设置小说简介
	e.SetIdentifier("urn:isbn:" + bookISBNStr) //设置小说的urn:isbn编码
	e.SetLang("zh-CN")                         //设置小说的语言为中文

	//设置章节信息
	chapters := this.Chapters
	for index := 0; index < len(chapters); index++ {
		sectionBody := "<h1>" + chapters[index].Title + "</h1>\n<p></p>\n"
		contentTmp := chapters[index].Content
		contentLines := strings.Split(contentTmp, "\r")
		for _, line := range contentLines {
			line = strings.TrimSpace(line) //删除每行所有的空格

			if line != "" {
				sectionBody += "<p>" + line + "</p>\n"
			}
		}
		sectionName := fmt.Sprintf("section%04d.xhtml", index)
		//fmt.Println(sectionName)
		if _, err := e.AddSection(sectionBody, chapters[index].Title, sectionName, ""); err != nil {
			return err
		}

	}

	if err = e.Write(outputpath + this.Name + "-" + this.Author + ".epub"); err != nil {
		return err
	}

	return nil
}
