package main

import (
	"strings"

	"io/ioutil"

	"github.com/Aiicy/htmlquery"
	"github.com/Unknwon/com"
)

type BookInfo struct {
	Name     string
	Author   string
	Chapters []Chapter
}

type Chapter struct {
	Title   string
	Content string
	Link    string
}

//读取文件内容，并存入string,最终返回
func ReadAllString(filename string) string {
	tmp, _ := ioutil.ReadFile(filename)
	return string(tmp)
}

//生成txt电子书
func (this BookInfo) GenerateTxt() {
	chapters := this.Chapters
	content := "" //用于存放章节内容
	outfpath := "./outputs/"
	outfname := outfpath + this.Name + ".txt"
	for index := 0; index < len(chapters); index++ {
		content = content + chapters[index].Title + "\n\n"
		content = content + chapters[index].Content + "\n\n"
	}

	com.WriteFile(outfname, []byte(content))
}

/*
func (this BookInfo) GenerateMobi() {

}
*/
func GetBookInfo(bookid string) BookInfo {

	var bi BookInfo
	var chapters []Chapter
	pollURL := "https://www.xbiquge6.com/" + bookid + "/"
	doc, _ := htmlquery.LoadURL(pollURL)

	//获取书名字
	bookNameMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:book_name']")
	bookName := htmlquery.SelectAttr(bookNameMeta, "content")
	//fmt.Println("bookName = ", bookName)

	//获取书作者
	AuthorMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:author']")
	author := htmlquery.SelectAttr(AuthorMeta, "content")
	//fmt.Println("author = ", author)

	//获取书章节列表
	ddNode, _ := htmlquery.Find(doc, "//div[@id='list']//dl//dd")
	for i := 0; i < len(ddNode); i++ {
		var tmp Chapter
		aNode, _ := htmlquery.Find(ddNode[i], "//a")
		tmp.Link = "https://www.xsbiquge.com" + htmlquery.SelectAttr(aNode[0], "href")
		tmp.Title = htmlquery.InnerText(aNode[0])
		chapters = append(chapters, tmp)
	}

	//导入信息
	bi = BookInfo{
		Name:     bookName,
		Author:   author,
		Chapters: chapters,
	}
	return bi
}

func GetChapterContent(C Chapter) Chapter {

	pollURL := C.Link
	doc, _ := htmlquery.LoadURL(pollURL)
	contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
	contentText := htmlquery.InnerText(contentNode)

	//替换字符串中的特殊字符 \xA0\XC2 为换行符 \n
	tmp := strings.Replace(contentText, "\xA0\xC2", "\r\n", -1)
	tmp = strings.Replace(tmp, "\xA0", "", -1)
	tmp = strings.Replace(tmp, "\xC2", "", -1)
	tmp = strings.Replace(tmp, "\xE6\xA1", "", -1)
	tmp = strings.Replace(tmp, "\xE7\xAB", "", -1)
	//把 readx(); 替换成 ""
	tmp = strings.Replace(tmp, "readx();", "", -1)
	tmp = tmp + "\r\n"
	//返回数据，填写Content内容
	result := Chapter{
		Title:   C.Title,
		Link:    C.Link,
		Content: tmp,
	}

	return result
}

//根据每个章节的 url连接，下载每章对应的内容Content当中
func (this BookInfo) DownloadChapters() BookInfo {
	chapters := this.Chapters

	for index := 0; index < len(chapters); index++ {
		temp := GetChapterContent(chapters[index])
		chapters[index].Content = temp.Content
	}

	result := BookInfo{
		Name:     this.Name,
		Author:   this.Author,
		Chapters: chapters,
	}

	return result
}
func main() {
	bookid := "0_642"
	bookinfo := GetBookInfo(bookid)
	//下载章节内容
	bookinfo.DownloadChapters()
	//生成txt文件
	bookinfo.GenerateTxt()
	//fmt.Println(bookinfo)
	//chapterLink := "https://www.xsbiquge.com//0_642/5776123.html"
	//contentText := GetChapterContent(chapterLink)
	//fmt.Println("chapterContent =", contentText)

}
