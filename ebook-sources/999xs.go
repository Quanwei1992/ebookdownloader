package ebook

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Aiicy/htmlquery"
	progressbar "github.com/schollz/progressbar/v2"
	edl "github.com/sndnvaps/ebookdownloader"
)

// https://www.899zw.com/

/*
bookid 规则



 https://www.899zw.com/files/article/html/0/591/ -> bookid = 591
 591 -> {0,591}

 https://www.899zw.com/files/article/html/1/1599/ -> bookid = 1599
 1599 -> {1,599}

https://www.899zw.com/files/article/html/75/75842/ -> bookid = 75842
 75842 - > {75,842}


 https://www.899zw.com/files/article/html/113/113582/ -> bookid = 113582
 113582 -> {113,582}
*/

var _ edl.EBookDLInterface = XS999{}

//999小说网 899zw.com
type XS999 struct {
	URL  string
	Lock *sync.Mutex
}

func New999XS() XS999 {
	return XS999{
		URL:  "https://www.899zw.com",
		Lock: new(sync.Mutex),
	}
}

/*
input: 113582 ,output: 113/113582
input: 75842, output: 75/75842
input: 1599, output: 1/1599
input: 591, output: 0/591
*/
func handleBookid(bookid string) string {
	tmp := []rune(bookid)
	if len(tmp) == 3 {
		return "0/" + string(tmp)
	}
	//最后结果不算怎样，都保留最后三个数字，前面n个数字需要分离出来
	return (string(tmp[:len(tmp)-3]) + "/" + string(tmp))
}

//GetBookBriefInfo 获取小说的信息
func (this XS999) GetBookBriefInfo(bookid string, proxy string) edl.BookInfo {

	var bi edl.BookInfo
	pollURL := this.URL + "/" + "files/article/html/" + handleBookid(bookid) + "/"

	//当 proxy 不为空的时候，表示设置代理
	if proxy != "" {
		doc, err := htmlquery.LoadURLWithProxy(pollURL, proxy)
		if err != nil {
			fmt.Println(err.Error())
		}

		//获取书名字
		bookNameMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:book_name']")
		bookName := htmlquery.SelectAttr(bookNameMeta, "content")
		bookName = SanitizeName(bookName)
		fmt.Println("书名 = ", bookName)

		//获取书作者
		AuthorMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:author']")
		author := htmlquery.SelectAttr(AuthorMeta, "content")
		author = SanitizeName(author)
		fmt.Println("作者 = ", author)

		//获取书的描述信息
		DescriptionMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:description']")
		description := htmlquery.SelectAttr(DescriptionMeta, "content")
		fmt.Println("简介 = ", description)

		//获取书的封面下载地址
		CoverURLMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:image']")
		CoverURL := htmlquery.SelectAttr(CoverURLMeta, "content")
		//fmt.Println("封面下载地址 = ", CoverURL)

		//导入信息
		bi = edl.BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
			CoverURL:    CoverURL,
		}
	} else { //没有设置代理
		doc, err := htmlquery.LoadURL(pollURL)
		if err != nil {
			fmt.Println(err.Error())
		}

		//获取书名字
		bookNameMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:book_name']")
		bookName := htmlquery.SelectAttr(bookNameMeta, "content")
		bookName = SanitizeName(bookName)
		fmt.Println("书名 = ", bookName)

		//获取书作者
		AuthorMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:author']")
		author := htmlquery.SelectAttr(AuthorMeta, "content")
		author = SanitizeName(author)
		fmt.Println("作者 = ", author)

		//获取书的描述信息
		DescriptionMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:description']")
		description := htmlquery.SelectAttr(DescriptionMeta, "content")
		fmt.Println("简介 = ", description)

		//获取书的封面下载地址
		CoverURLMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:image']")
		CoverURL := htmlquery.SelectAttr(CoverURLMeta, "content")
		//fmt.Println("封面下载地址 = ", CoverURL)

		//导入信息
		bi = edl.BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
			CoverURL:    CoverURL,
		}
	}
	return bi
}

func (this XS999) GetBookInfo(bookid string, proxy string) edl.BookInfo {

	var bi edl.BookInfo
	var chapters []edl.Chapter
	pollURL := this.URL + "/" + "files/article/html/" + handleBookid(bookid) + "/"

	//当 proxy 不为空的时候，表示设置代理
	if proxy != "" {
		doc, err := htmlquery.LoadURLWithProxy(pollURL, proxy)
		if err != nil {
			fmt.Println(err.Error())
		}

		//获取书名字
		bookNameMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:book_name']")
		bookName := htmlquery.SelectAttr(bookNameMeta, "content")
		bookName = SanitizeName(bookName)
		fmt.Println("书名 = ", bookName)

		//获取书作者
		AuthorMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:author']")
		author := htmlquery.SelectAttr(AuthorMeta, "content")
		author = SanitizeName(author)
		fmt.Println("作者 = ", author)

		//获取书的描述信息
		DescriptionMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:description']")
		description := htmlquery.SelectAttr(DescriptionMeta, "content")
		fmt.Println("简介 = ", description)

		//获取书的封面下载地址
		CoverURLMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:image']")
		CoverURL := htmlquery.SelectAttr(CoverURLMeta, "content")
		//fmt.Println("封面下载地址 = ", CoverURL)

		//获取书章节列表
		ddNode, _ := htmlquery.Find(doc, "//div[@id='list']//dl//dd")
		for i := 0; i < len(ddNode); i++ {
			var tmp edl.Chapter
			aNode, _ := htmlquery.Find(ddNode[i], "//a")
			tmp.Link = this.URL + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])
			chapters = append(chapters, tmp)
		}

		//导入信息
		bi = edl.BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
			CoverURL:    CoverURL,
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
		bookName = SanitizeName(bookName)
		fmt.Println("书名 = ", bookName)

		//获取书作者
		AuthorMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:author']")
		author := htmlquery.SelectAttr(AuthorMeta, "content")
		author = SanitizeName(author)
		fmt.Println("作者 = ", author)

		//获取书的描述信息
		DescriptionMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:description']")
		description := htmlquery.SelectAttr(DescriptionMeta, "content")
		fmt.Println("简介 = ", description)

		//获取书的封面下载地址
		CoverURLMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:image']")
		CoverURL := htmlquery.SelectAttr(CoverURLMeta, "content")
		//fmt.Println("封面下载地址 = ", CoverURL)

		//获取书章节列表
		ddNode, _ := htmlquery.Find(doc, "//div[@id='list']//dl//dd")
		for i := 0; i < len(ddNode); i++ {
			var tmp edl.Chapter
			aNode, _ := htmlquery.Find(ddNode[i], "//a")
			tmp.Link = this.URL + "/" + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])
			chapters = append(chapters, tmp)
		}
		//导入信息
		bi = edl.BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
			CoverURL:    CoverURL,
			Chapters:    chapters,
		}

	}
	//生成ISBN码
	bi.GenerateISBN()
	//生成UUID
	bi.GenerateUUID()
	return bi
}
func (this XS999) DownloadChapters(Bi edl.BookInfo, proxy string) edl.BookInfo {
	result := Bi //先进行赋值，把数据
	var chapters []edl.Chapter
	result.Chapters = chapters //把原来的数据清空
	bis := Bi.Split()

	for index := 0; index < len(bis); index++ {
		this.Lock.Lock()
		bookinfo := bis[index]
		rec := this.downloadChapters(bookinfo, "")
		chapters = append(chapters, rec.Chapters...)
		//fmt.Printf("Get into this.Lock.Unlock() time: %d\n", index+1)
		this.Lock.Unlock()
	}
	result.Chapters = chapters

	return result
}

//根据每个章节的 URL连接，下载每章对应的内容Content当中
func (this XS999) downloadChapters(Bi edl.BookInfo, proxy string) edl.BookInfo {
	chapters := Bi.Chapters

	NumChapter := len(chapters)
	tmpChapter := make(chan edl.Chapter, NumChapter)
	ResultCh := make(chan chan edl.Chapter, NumChapter)
	wg := sync.WaitGroup{}
	var c []edl.Chapter
	var bar *progressbar.ProgressBar
	go AsycChapter(ResultCh, tmpChapter)
	for index := 0; index < NumChapter; index++ {
		tmp := edl.ProxyChapter{
			Proxy: proxy,
			C:     chapters[index],
		}
		this.DownloaderChapter(ResultCh, tmp, &wg)
	}

	wg.Wait()

	//下载章节的时候显示进度条
	bar = progressbar.NewOptions(
		NumChapter,
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowIts(),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "#", SaucerPadding: "-", BarStart: ">", BarEnd: "<"}),
	)

	for index := 0; index <= NumChapter; {
		select {
		case tmp := <-tmpChapter:
			//fmt.Printf("tmp.Title = %s\n", tmp.Title)
			//fmt.Printf("tmp.Content= %s\n", tmp.Content)
			c = append(c, tmp)
			index++
			bar.Add(1)
			if index == NumChapter {
				goto ForEnd
			}
		}

	}
ForEnd:

	result := edl.BookInfo{
		EBHost:      Bi.EBHost,
		EBookID:     Bi.EBookID,
		BookISBN:    Bi.ISBN(),
		BookUUID:    Bi.UUID(),
		Name:        Bi.Name,
		Author:      Bi.Author,
		Description: Bi.Description,
		Volumes:     Bi.Volumes,       //小说分卷信息在 GetBookInfo()的时候已经下载完成
		HasVolume:   Bi.VolumeState(), //小说分卷信息在 GetBookInfo()的时候已经定义
		Chapters:    c,
	}

	return result
}

//DownloaderChapter 下载小说
func (this XS999) DownloaderChapter(ResultChan chan chan edl.Chapter, pc edl.ProxyChapter, wg *sync.WaitGroup) {
	c := make(chan edl.Chapter)
	ResultChan <- c
	wg.Add(1)
	go func(pc edl.ProxyChapter) {
		pollURL := pc.C.Link
		proxy := pc.Proxy
		var result edl.Chapter

		if proxy != "" {
			doc, _ := htmlquery.LoadURLWithProxy(pollURL, proxy)
			contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
			contentText := htmlquery.OutputHTML(contentNode, false)

			//替换两个 html换行
			tmp := strings.Replace(contentText, "<br/><br/>", "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(tmp, "<br/>", "\r\n", -1)

			//把 readx(); 替换成 ""
			tmp = strings.Replace(tmp, "999小说更新最快 电脑端:https://www.999xs.com/", "", -1)
			tmp = strings.Replace(tmp, "ωωω.九九九xs.com", "", -1)
			tmp = strings.Replace(tmp, "999小说首发 https://www.999xs.com https://m.999xs.com", "", -1)
			tmp = strings.Replace(tmp, "手机\\端 一秒記住『www.999xs.com』為您提\\供精彩小說\\閱讀", "", -1)

			//tmp = tmp + "\r\n"
			//返回数据，填写Content内容
			result = edl.Chapter{
				Title:   pc.C.Title,
				Link:    pc.C.Link,
				Content: tmp,
			}
		} else {
			doc, _ := htmlquery.LoadURL(pollURL)
			contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
			contentText := htmlquery.OutputHTML(contentNode, false)

			//替换两个 html换行
			tmp := strings.Replace(contentText, "<br/><br/>", "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(tmp, "<br/>", "\r\n", -1)

			//把 readx(); 替换成 ""
			tmp = strings.Replace(tmp, "999小说更新最快 电脑端:https://www.999xs.com/", "", -1)
			tmp = strings.Replace(tmp, "ωωω.九九九xs.com", "", -1)
			tmp = strings.Replace(tmp, "999小说首发 https://www.999xs.com https://m.999xs.com", "", -1)
			tmp = strings.Replace(tmp, "手机\\端 一秒記住『www.999xs.com』為您提\\供精彩小說\\閱讀", "", -1)

			//tmp = tmp + "\r\n"
			//返回数据，填写Content内容
			result = edl.Chapter{
				Title:   pc.C.Title,
				Link:    pc.C.Link,
				Content: tmp,
			}
		}
		c <- result
		wg.Done()
	}(pc)
}
