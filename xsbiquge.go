package main

import (
	"fmt"
	"sync"

	"strings"

	"github.com/Aiicy/htmlquery"
	pool "github.com/dgrr/goslaves"
	"gopkg.in/schollz/progressbar.v2"
)

//新笔趣阁 xsbiquge.com
type EbookXSBiquge struct {
	Url string
}

func NewXSBiquge() EbookXSBiquge {
	return EbookXSBiquge{
		Url: "https://www.xsbiquge.com",
	}
}

func (this EbookXSBiquge) GetBookInfo(bookid string, proxy string) BookInfo {

	var bi BookInfo
	var chapters []Chapter
	pollURL := this.Url + "/" + bookid + "/"

	//当 proxy 不为空的时候，表示设置代理
	if proxy != "" {
		doc, err := htmlquery.LoadURLWithProxy(pollURL, proxy)
		if err != nil {
			fmt.Println(err.Error())
		}

		//获取书名字
		bookNameMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:book_name']")
		bookName := htmlquery.SelectAttr(bookNameMeta, "content")
		fmt.Println("书名 = ", bookName)

		//获取书作者
		AuthorMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:author']")
		author := htmlquery.SelectAttr(AuthorMeta, "content")
		fmt.Println("作者 = ", author)

		//获取书的描述信息
		DescriptionMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:description']")
		description := htmlquery.SelectAttr(DescriptionMeta, "content")
		fmt.Println("简介 = ", description)

		//获取书章节列表
		ddNode, _ := htmlquery.Find(doc, "//div[@id='list']//dl//dd")
		for i := 0; i < len(ddNode); i++ {
			var tmp Chapter
			aNode, _ := htmlquery.Find(ddNode[i], "//a")
			tmp.Link = this.Url + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])
			chapters = append(chapters, tmp)
		}

		//导入信息
		bi = BookInfo{
			Name:        bookName,
			Author:      author,
			Description: description,
			Chapters:    chapters,
		}
	} else { //没有设置代理
		doc, err := htmlquery.LoadURL(pollURL)
		if err != nil {
			fmt.Println(err.Error())
		}

		//获取书名字
		bookNameMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:book_name']")
		bookName := htmlquery.SelectAttr(bookNameMeta, "content")
		fmt.Println("书名 = ", bookName)

		//获取书作者
		AuthorMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:author']")
		author := htmlquery.SelectAttr(AuthorMeta, "content")
		fmt.Println("作者 = ", author)

		//获取书的描述信息
		DescriptionMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:description']")
		description := htmlquery.SelectAttr(DescriptionMeta, "content")
		fmt.Println("简介 = ", description)

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
			Name:        bookName,
			Author:      author,
			Description: description,
			Chapters:    chapters,
		}
	}
	return bi
}

func (this EbookXSBiquge) GetChapterContent(pc ProxyChapter) Chapter {
	pollURL := pc.C.Link
	proxy := pc.Proxy
	var result Chapter

	if proxy != "" {
		doc, _ := htmlquery.LoadURLWithProxy(pollURL, proxy)
		contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
		contentText := htmlquery.InnerText(contentNode)

		//替换字符串中的特殊字符 \xC2\xA0 为换行符 \n
		tmp := strings.Replace(contentText, "\xC2\xA0", "\r\n", -1)

		//把 readx(); 替换成 ""
		tmp = strings.Replace(tmp, "readx();", "", -1)
		//tmp = tmp + "\r\n"
		//返回数据，填写Content内容
		result = Chapter{
			Title:   pc.C.Title,
			Link:    pc.C.Link,
			Content: tmp,
		}
	} else {
		doc, _ := htmlquery.LoadURL(pollURL)
		contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
		contentText := htmlquery.InnerText(contentNode)

		//替换字符串中的特殊字符 \xC2\xA0 为换行符 \n
		tmp := strings.Replace(contentText, "\xC2\xA0", "\r\n", -1)

		//把 readx(); 替换成 ""
		tmp = strings.Replace(tmp, "readx();", "", -1)
		//tmp = tmp + "\r\n"
		//返回数据，填写Content内容
		result = Chapter{
			Title:   pc.C.Title,
			Link:    pc.C.Link,
			Content: tmp,
		}
	}

	return result
}

func excuteServe(p *pool.Pool, chapters []Chapter, proxy string) {
	for i := 0; i < len(chapters); i++ {
		tmp := ProxyChapter{
			Proxy: proxy,
			C:     chapters[i],
		}
		p.Serve(tmp)
	}
}

//根据每个章节的 url连接，下载每章对应的内容Content当中
func (this EbookXSBiquge) DownloadChapters(Bi BookInfo, proxy string) BookInfo {
	chapters := Bi.Chapters

	NumChapter := len(chapters)
	tmpChapter := make(chan Chapter, NumChapter)
	ResultCh := make(chan chan Chapter, NumChapter)
	wg := sync.WaitGroup{}
	var c []Chapter
	var bar *progressbar.ProgressBar
	go AsycChapter(ResultCh, tmpChapter)
	for index := 0; index < NumChapter; index++ {
		tmp := ProxyChapter{
			Proxy: proxy,
			C:     chapters[index],
		}
		this.DownloaderChapter(ResultCh, tmp, &wg)
	}

	wg.Wait()

	//下载章节的时候显示进度条
	bar = progressbar.New(NumChapter)
	bar.RenderBlank()

	for index := 0; index < NumChapter; {
		select {
		case tmp := <-tmpChapter:
			//fmt.Printf("tmp.Title = %s\n", tmp.Title)
			//fmt.Printf("tmp.Content= %s\n", tmp.Content)
			c = append(c, tmp)
			index++
			if index == (NumChapter - 1) {
				goto ForEnd
			}
		}
		bar.Add(1)

	}
ForEnd:

	result := BookInfo{
		Name:        Bi.Name,
		Author:      Bi.Author,
		Description: Bi.Description,
		Volumes:     Bi.Volumes,       //小说分卷信息在 GetBookInfo()的时候已经下载完成
		HasVolume:   Bi.VolumeState(), //小说分卷信息在 GetBookInfo()的时候已经定义
		Chapters:    c,
	}

	return result
}

//func DownloaderChapter(ResultChan chan chan Chapter)
func (this EbookXSBiquge) DownloaderChapter(ResultChan chan chan Chapter, pc ProxyChapter, wg *sync.WaitGroup) {
	c := make(chan Chapter)
	ResultChan <- c
	wg.Add(1)
	go func(pc ProxyChapter) {
		pollURL := pc.C.Link
		proxy := pc.Proxy
		var result Chapter

		if proxy != "" {
			doc, _ := htmlquery.LoadURLWithProxy(pollURL, proxy)
			contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
			contentText := htmlquery.InnerText(contentNode)

			//替换字符串中的特殊字符 \xC2\xA0 为换行符 \n
			tmp := strings.Replace(contentText, "\xC2\xA0", "\r\n", -1)

			//把 readx(); 替换成 ""
			tmp = strings.Replace(tmp, "readx();", "", -1)
			//tmp = tmp + "\r\n"
			//返回数据，填写Content内容
			result = Chapter{
				Title:   pc.C.Title,
				Link:    pc.C.Link,
				Content: tmp,
			}
		} else {
			doc, _ := htmlquery.LoadURL(pollURL)
			contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
			contentText := htmlquery.InnerText(contentNode)

			//替换字符串中的特殊字符 \xC2\xA0 为换行符 \n
			tmp := strings.Replace(contentText, "\xC2\xA0", "\r\n", -1)

			//把 readx(); 替换成 ""
			tmp = strings.Replace(tmp, "readx();", "", -1)
			//tmp = tmp + "\r\n"
			//返回数据，填写Content内容
			result = Chapter{
				Title:   pc.C.Title,
				Link:    pc.C.Link,
				Content: tmp,
			}
		}
		c <- result
		wg.Done()
	}(pc)
}
