package ebook

import (
	"context"
	"fmt"
	"sync"

	"strings"

	"github.com/Aiicy/htmlquery"
	progressbar "github.com/schollz/progressbar/v2"
	edl "github.com/sndnvaps/ebookdownloader"
)

var _ edl.EBookDLInterface = Biqugei{}

// Biqugei http://www.biqugei.net/小说网
type Biqugei struct {
	URL  string
	Lock *sync.Mutex
}

// NewBiqugei 初始化
func NewBiqugei() Biqugei {
	return Biqugei{
		URL:  "http://www.biqugei.net",
		Lock: new(sync.Mutex),
	}
}

// GetBookDownloadLinkPages 获取小说每个部分的下载页面，以50章为一个下载页面
func (this Biqugei) GetBookDownloadLinkPages(bookid string, proxy string) []string {
	var chaptersPages []string
	pollURL := this.URL + "/book/" + bookid + ".html"
	doc, err := htmlquery.LoadURL(pollURL)
	if err != nil {
		fmt.Println(err.Error())
	}
	optionRootNode, _ := htmlquery.Find(doc, "//div[@class='index-container']//select[@id='indexselect']")
	//ulNode[0] 为最新章节的下载链接
	//ulNode[1] 为正文的下载链接
	optionNodes, _ := htmlquery.Find(optionRootNode[0], "//option")
	for i := 0; i < len(optionNodes)-1; i++ {
		tmpUrl := this.URL + htmlquery.SelectAttr(optionNodes[i], "value")
		fmt.Println(tmpUrl)
		chaptersPages = append(chaptersPages, tmpUrl)
	}
	return chaptersPages

}

// GetBookBriefInfo 获取小说的信息
func (this Biqugei) GetBookBriefInfo(bookid string, proxy string) edl.BookInfo {
	var bi edl.BookInfo
	pollURL := this.URL + "/book/" + bookid + ".html"

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

// GetBookInfo 获取小说的信息
func (this Biqugei) GetBookInfo(ctx context.Context, bookid string, proxy string) edl.BookInfo {

	bi := this.GetBookBriefInfo(bookid, proxy)
	var chapters []edl.Chapter
	//for test
	chaptersPages := this.GetBookDownloadLinkPages(bookid, proxy)
	for i := 0; i < len(chaptersPages)-1; i++ {

		//当 proxy 不为空的时候，表示设置代理
		doc, err := htmlquery.LoadURL(chaptersPages[i])
		if err != nil {
			fmt.Println(err.Error())
		}

		//获取书章节列表
		ulNode, _ := htmlquery.Find(doc, "//div[@class='layout layout-col1']//div[@class='section-box']//ul[@class='section-list fix']")
		//ulNode[0] 为最新章节的下载链接
		//ulNode[1] 为正文的下载链接
		fmt.Println(ulNode)
		liNode, _ := htmlquery.Find(ulNode[1], "//li")
		for i := 0; i < len(liNode)-1; i++ {
			var tmp edl.Chapter
			aNode, _ := htmlquery.Find(liNode[i], "//a")
			tmp.Link = this.URL + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])
			chapters = append(chapters, tmp)
		}
	}

	//导入信息
	bookinfo := edl.BookInfo{
		EBHost:      this.URL,
		EBookID:     bookid,
		Author:      bi.Author,
		Name:        bi.Name,
		CoverURL:    bi.CoverURL,
		Description: bi.Description,
		Chapters:    chapters,
	}

	//生成ISBN码
	bookinfo.GenerateISBN()
	//生成UUID
	bookinfo.GenerateUUID()
	return bookinfo
}

// DownloadChapters 下载所有章节
func (this Biqugei) DownloadChapters(Bi edl.BookInfo, proxy string) edl.BookInfo {
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

// 根据每个章节的 URL连接，下载每章对应的内容Content当中
func (this Biqugei) downloadChapters(Bi edl.BookInfo, proxy string) edl.BookInfo {
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

// DownloaderChapter 一个章节一个章节得下载
func (this Biqugei) DownloaderChapter(ResultChan chan chan edl.Chapter, pc edl.ProxyChapter, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	c := make(chan edl.Chapter)
	ResultChan <- c

	go func(pc edl.ProxyChapter) {
		pollURL := pc.C.Link
		proxy := pc.Proxy
		var result edl.Chapter

		if proxy != "" {
			doc, _ := htmlquery.LoadURLWithProxy(pollURL, proxy)
			contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
			contentText := htmlquery.OutputHTML(contentNode, false)

			//删除章节开始的红色文字
			rpstr := fmt.Sprintf("%s", `<div class="posterror"><a href="javascript:report();" class="red">章节错误,点此举报(免注册)</a>,举报后维护人员会在两分钟内校正章节内容,请耐心等待,并刷新页面。</div>`)
			tmp := strings.Replace(contentText, rpstr, "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(contentText, "<p>", "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(tmp, "</p>", "\r\n", -1)

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

			rpstr := fmt.Sprintf("%s", `<div class="posterror"><a href="javascript:report();" class="red">章节错误,点此举报(免注册)<\/a>,举报后维护人员会在两分钟内校正章节内容,请耐心等待,并刷新页面。<\/div>`)
			tmp := strings.Replace(contentText, rpstr, "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(contentText, "<p>", "		", -1)
			//替换一个 html换行
			tmp = strings.Replace(tmp, "</p>", "\r\n", -1)
			//tmp = tmp + "\r\n"
			//返回数据，填写Content内容
			result = edl.Chapter{
				Title:   pc.C.Title,
				Link:    pc.C.Link,
				Content: tmp,
			}
		}
		c <- result
	}(pc)
}
