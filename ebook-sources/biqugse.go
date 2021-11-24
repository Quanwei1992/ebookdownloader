package ebook

import (
	"fmt"
	"sync"

	"strings"

	"github.com/Aiicy/htmlquery"
	progressbar "github.com/schollz/progressbar/v2"
	edl "github.com/sndnvaps/ebookdownloader"
)

var _ edl.EBookDLInterface = Biqugse{}

// biqugse http://www.biqugse.com 小说网
type Biqugse struct {
	URL  string
	Lock *sync.Mutex
}

// NewBiqugse 初始化
func NewBiqugse() Biqugse {
	return Biqugse{
		URL:  "http://www.Biqugse.com/",
		Lock: new(sync.Mutex),
	}
}

//GetBookBriefInfo 获取小说的信息
func (this Biqugse) GetBookBriefInfo(bookid string, proxy string) edl.BookInfo {
	var bi edl.BookInfo
	pollURL := this.URL + bookid + "/"

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

//GetBookInfo 获取小说的信息
func (this Biqugse) GetBookInfo(bookid string, proxy string) edl.BookInfo {

	var bi edl.BookInfo
	var chapters []edl.Chapter
	pollURL := this.URL + bookid + "/"

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
		//i := 9，因为最前面的9章是：显示最新章节信息，需要忽略掉
		for i := 9; i < len(ddNode); i++ {
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
		//i := 9，因为最前面的9章是：显示最新章节信息，需要忽略掉
		for i := 9; i < len(ddNode); i++ {
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

	}
	//生成ISBN码
	bi.GenerateISBN()
	//生成UUID
	bi.GenerateUUID()
	return bi
}

//DownloadChapters 下载所有章节
func (this Biqugse) DownloadChapters(Bi edl.BookInfo, proxy string) edl.BookInfo {
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
func (this Biqugse) downloadChapters(Bi edl.BookInfo, proxy string) edl.BookInfo {
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

//DownloaderChapter 一个章节一个章节得下载
func (this Biqugse) DownloaderChapter(ResultChan chan chan edl.Chapter, pc edl.ProxyChapter, wg *sync.WaitGroup) {
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

			//替换两个 html换行
			tmp := strings.Replace(contentText, "<br/><br/>", "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(tmp, "<br/>", "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(tmp, "<br>", "\r\n", -1)

			//把 readx(); 替换成 ""
			tmp = strings.Replace(tmp, "readx();", "", -1)
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
			//替换一个 html换行
			tmp = strings.Replace(tmp, "<br>", "\r\n", -1)
			//替换一个 html换行
			//&lt;/p&gt; -> </p>
			tmp = strings.Replace(tmp, "&lt;/p&gt;", "", -1)
			// &nbsp; -> 代表一个html空格
			tmp = strings.Replace(tmp, "&nbsp;&nbsp;&nbsp;&nbsp;", "    ", -1)

			//把 readx(); 替换成 ""
			tmp = strings.Replace(tmp, "readx();", "", -1)
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
