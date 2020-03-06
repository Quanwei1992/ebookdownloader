package ebook

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Aiicy/htmlquery"
	"github.com/schollz/progressbar/v2"
	edl "github.com/sndnvaps/ebookdownloader"
)

var _ edl.EBookDLInterface = BookTXT{}

//BookTXT 顶点小说网 www.booktxt.net
type BookTXT struct {
	URL  string
	Lock *sync.Mutex
}

//NewBookTXT 初始化
func NewBookTXT() BookTXT {
	return BookTXT{
		URL:  "https://www.booktxt.net/",
		Lock: new(sync.Mutex),
	}
}

//GetBookBriefInfo 获取小说的信息
func (this BookTXT) GetBookBriefInfo(bookid string, proxy string) edl.BookInfo {

	var bi edl.BookInfo
	pollURL := this.URL + bookid
	fmt.Printf("pollURL = %s", pollURL)

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
		bi = edl.BookInfo{
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
		bi = edl.BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
		}
	}
	return bi
}

//GetBookInfo 下载小说信息
func (this BookTXT) GetBookInfo(bookid string, proxy string) edl.BookInfo {

	var bi edl.BookInfo
	var volumes []edl.Volume
	var chapters []edl.Chapter
	pollURL := this.URL + bookid

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

		//替换掉 volume是最前面的 作品名字
		replaceStr := fmt.Sprintf("《%s》", bookName)
		//获取书分卷信息
		dtNode, _ := htmlquery.Find(doc, "//div[@id='list']//dl/dd") //获取书分卷信息
		testVolStr := htmlquery.InnerText(dtNode[1])

		if TestContainVolume(testVolStr) {
			bi.ChangeVolumeState(true)
			if len(dtNode) == 2 { //就是说刚好两个节点，我们要去除第一个，只保留第二个
				var tmp edl.Volume
				tmp.CurrentVolume = htmlquery.InnerText(dtNode[1])
				volumes = append(volumes, tmp)
			} else { //当len(dtNode) >= 3
				for index := 0; index < len(dtNode); index++ { //因为第一个为 最新章节部分，需要去掉
					var tmp edl.Volume
					// 根据当前节点，查找上一个dd节点
					PrevChapter, _ := htmlquery.FindOne(dtNode[index], "//preceding-sibling::dd[1]")
					aNode, _ := htmlquery.Find(PrevChapter, "//a")
					tmp.PrevChapter.Link = this.URL + bookid + "/" + htmlquery.SelectAttr(aNode[0], "href")
					tmp.PrevChapter.Title = htmlquery.InnerText(aNode[0])

					//根据当前节点，查找下一个dd节点
					NextChapter, _ := htmlquery.FindOne(dtNode[index], "//following-sibling::dd[1]")
					aNode, _ = htmlquery.Find(NextChapter, "//a")
					tmp.NextChapter.Link = this.URL + bookid + "/" + htmlquery.SelectAttr(aNode[0], "href")
					CurrentVolume := htmlquery.InnerText(dtNode[index])
					tmp.CurrentVolume = strings.Replace(CurrentVolume, replaceStr, "", -1)
					tmp.NextChapter.Title = htmlquery.InnerText(aNode[0])
					volumes = append(volumes, tmp)
				}
			}
			volumes[0].PrevChapterID = 0      //第一分卷，前面的章节，设置为0
			volumes[0].PrevChapter.Link = ""  //第一分卷，前面的章节，连接设置为空
			volumes[0].PrevChapter.Title = "" //第一分卷，前面的章节，标题设置为空
		}
		//获取书章节列表
		ddNode, _ := htmlquery.Find(doc, "//div[@id='list']//dl/dd")
		//i := 5，因为最前面的6章是：显示最新章节信息，需要忽略掉
		for i := 5; i < len(ddNode); i++ {
			var tmp edl.Chapter
			aNode, _ := htmlquery.Find(ddNode[i], "//a")
			tmp.Link = this.URL + bookid + "/" + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])
			if bi.VolumeState() && len(volumes) >= 2 { //正式写入 PrevChapterID
				for index := 1; index < len(volumes); index++ { //第二个分卷开始，前面就有章节内容了
					if volumes[index].PrevChapter.Link == tmp.Link {
						volumes[index].PrevChapterID = i
					}
				}
			}
			chapters = append(chapters, tmp)
		}
		HasVolume := bi.VolumeState() //先赋值给 HasVolume,再把值导入到结构体中，用于数据返回

		//导入信息
		bi = edl.BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
			Volumes:     volumes,
			HasVolume:   HasVolume,
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

		//替换掉 volume是最前面的 作品名字
		replaceStr := fmt.Sprintf("《%s》", bookName)

		//获取书分卷信息
		dtNode, _ := htmlquery.Find(doc, "//div[@id='list']//dl/dt") //获取书分卷信息
		testVolStr := htmlquery.InnerText(dtNode[1])
		if TestContainVolume(testVolStr) {
			bi.ChangeVolumeState(true)
			if len(dtNode) == 2 { //就是说刚好两个节点，我们要去除第一个，只保留第二个
				var tmp edl.Volume
				tmp.CurrentVolume = htmlquery.InnerText(dtNode[1])
				volumes = append(volumes, tmp)
			} else { //当len(dtNode) >= 3
				for index := 1; index < len(dtNode); index++ { //因为第一个为 最新章节部分，需要去掉
					var tmp edl.Volume
					// 根据当前节点，查找上一个dd节点
					PrevChapter, _ := htmlquery.FindOne(dtNode[index], "//preceding-sibling::dd[1]")
					aNode, _ := htmlquery.Find(PrevChapter, "//a")
					tmp.PrevChapter.Link = this.URL + bookid + "/" + htmlquery.SelectAttr(aNode[0], "href")
					tmp.PrevChapter.Title = htmlquery.InnerText(aNode[0])

					//根据当前节点，查找下一个dd节点
					NextChapter, _ := htmlquery.FindOne(dtNode[index], "//following-sibling::dd[1]")
					aNode, _ = htmlquery.Find(NextChapter, "//a")
					tmp.NextChapter.Link = this.URL + bookid + "/" + htmlquery.SelectAttr(aNode[0], "href")
					tmp.NextChapter.Title = htmlquery.InnerText(aNode[0])
					CurrentVolume := htmlquery.InnerText(dtNode[index])
					tmp.CurrentVolume = strings.Replace(CurrentVolume, replaceStr, "", -1)
					volumes = append(volumes, tmp)
				}
			}
			volumes[0].PrevChapterID = 0      //第一分卷，前面的章节，设置为0
			volumes[0].PrevChapter.Link = ""  //第一分卷，前面的章节，连接设置为空
			volumes[0].PrevChapter.Title = "" //第一分卷，前面的章节，标题设置为空
		}
		//获取书章节列表
		ddNode, _ := htmlquery.Find(doc, "//div[@id='list']//dl/dd")
		for i := 6; i < len(ddNode); i++ { //因为前面的6个ddNode为显示最新的12章，与后面的会重复，所以直接Drop
			var tmp edl.Chapter
			aNode, _ := htmlquery.Find(ddNode[i], "//a")
			tmp.Link = this.URL + bookid + "/" + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])
			//fmt.Printf("tmp.Link = %s\n", tmp.Link)   //用于测试
			//fmt.Printf("tmp.Title = %s\n", tmp.Title) //用于测试

			if bi.VolumeState() && len(volumes) >= 2 { //正式写入 PrevChapterID && NextChapterID
				for index := 0; index < len(volumes); index++ {
					if volumes[index].PrevChapter.Link == tmp.Link {
						volumes[index].PrevChapterID = (i - 12) + 1 //表示 设置 第一个章节为0
					}
					if volumes[index].NextChapter.Link == tmp.Link {
						volumes[index].NextChapterID = (i - 12) + 1 //表示 设置 第一个章节为0
					}
				}
			}
			chapters = append(chapters, tmp)
		}
		HasVolume := bi.VolumeState() //先赋值给 HasVolume,再把值导入到结构体中，用于数据返回

		//导入信息
		bi = edl.BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
			Volumes:     volumes,
			HasVolume:   HasVolume,
			Chapters:    chapters,
		}
	}
	//生成ISBN码
	bi.GenerateISBN()
	//生成UUID
	bi.GenerateUUID()
	return bi
}

//DownloadChapters 下载小说章节
func (this BookTXT) DownloadChapters(Bi edl.BookInfo, proxy string) edl.BookInfo {
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

//根据每个章节的 url连接，下载每章对应的内容Content当中
func (this BookTXT) downloadChapters(Bi edl.BookInfo, proxy string) edl.BookInfo {
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
			if index == NumChapter {
				goto ForEnd
			}
		}
		bar.Add(1)

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

// DownloaderChapter 下载小说章节
func (this BookTXT) DownloaderChapter(ResultChan chan chan edl.Chapter, pc edl.ProxyChapter, wg *sync.WaitGroup) {
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
			tmp = strings.Replace(tmp, "</p>", "", -1)
			tmp = strings.Replace(tmp, "(https://)", "", -1)

			tmp = strings.Replace(tmp, "请记住本书首发域名：booktxt.net。顶点小说手机版阅读网址：m.booktxt.net", "", -1)

			//tmp = tmp + "\r\n"
			//返回数据，填写Content内容
			result = edl.Chapter{
				Title:   pc.C.Title,
				Link:    pc.C.Link,
				Content: tmp,
			}
		} else {
			doc, _ := htmlquery.LoadURL(pollURL)
			contentNode, err := htmlquery.FindOne(doc, "//div[@id='content']")
			if err != nil {
				fmt.Println(err.Error())
			}
			contentText := htmlquery.OutputHTML(contentNode, true)

			//替换两个 html换行
			tmp := strings.Replace(contentText, "<br/><br/>", "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(tmp, "<br/>", "\r\n", -1)

			//把 readx(); 替换成 ""
			tmp = strings.Replace(tmp, "</p>", "", -1)
			tmp = strings.Replace(tmp, "(https://)", "", -1)

			tmp = strings.Replace(tmp, "请记住本书首发域名：booktxt.net。顶点小说手机版阅读网址：m.booktxt.net", "", -1)

			//tmp = tmp + "\r\n"
			//返回数据，填写Content内容
			result = edl.Chapter{
				Title:   pc.C.Title,
				Link:    pc.C.Link,
				Content: tmp,
			}
		}
		//fmt.Printf("result.Content= %s\n", result.Content)
		c <- result
		wg.Done()
	}(pc)
}
