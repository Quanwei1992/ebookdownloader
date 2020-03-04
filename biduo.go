package ebookdownloader

import (
	"fmt"
	"sync"

	"strings"

	"github.com/Aiicy/htmlquery"
	"github.com/schollz/progressbar/v2"
)

var _ EBookDLInterface = EbookBiDuo{}

// EbookBiDuo www.biduo.cc小说网
type EbookBiDuo struct {
	URL  string
	Lock *sync.Mutex
}

// NewBiDuo 初始化
func NewBiDuo() EbookBiDuo {
	return EbookBiDuo{
		URL:  "https://www.biduo.cc",
		Lock: new(sync.Mutex),
	}
}

//GetBookBriefInfo 获取小说的信息
func (this EbookBiDuo) GetBookBriefInfo(bookid string, proxy string) BookInfo {
	var bi BookInfo
	pollURL := this.URL + "/biquge/" + bookid + "/"

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

		//导入信息
		bi = BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
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

		//导入信息
		bi = BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
		}
	}
	return bi
}

//GetBookInfo 获取小说的信息
func (this EbookBiDuo) GetBookInfo(bookid string, proxy string) BookInfo {

	var bi BookInfo
	var chapters []Chapter
	pollURL := this.URL + "/biquge/" + bookid + "/"

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
			tmp.Link = this.URL + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])
			chapters = append(chapters, tmp)
		}

		//导入信息
		bi = BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
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
			tmp.Link = this.URL + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])
			chapters = append(chapters, tmp)
		}

		//导入信息
		bi = BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
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
func (this EbookBiDuo) DownloadChapters(Bi BookInfo, proxy string) BookInfo {
	result := Bi //先进行赋值，把数据
	var chapters []Chapter
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
func (this EbookBiDuo) downloadChapters(Bi BookInfo, proxy string) BookInfo {
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
			if index == NumChapter {
				goto ForEnd
			}
		}
		bar.Add(1)

	}
ForEnd:

	result := BookInfo{
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
func (this EbookBiDuo) DownloaderChapter(ResultChan chan chan Chapter, pc ProxyChapter, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	c := make(chan Chapter)
	ResultChan <- c

	go func(pc ProxyChapter) {
		pollURL := pc.C.Link
		proxy := pc.Proxy
		var result Chapter

		if proxy != "" {
			doc, _ := htmlquery.LoadURLWithProxy(pollURL, proxy)
			contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
			contentText := htmlquery.OutputHTML(contentNode, false)

			//替换两个 html换行
			tmp := strings.Replace(contentText, "<br/><br/>", "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(tmp, "<br/>", "\r\n", -1)

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
			contentText := htmlquery.OutputHTML(contentNode, false)

			//替换两个 html换行
			tmp := strings.Replace(contentText, "<br/><br/>", "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(tmp, "<br/>", "\r\n", -1)

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
	}(pc)
}
