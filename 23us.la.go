package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Aiicy/htmlquery"
	pool "github.com/dgrr/goslaves"
	"gopkg.in/schollz/progressbar.v2"
)

//参考地址，创建规则
//https://www.23us.la/html/151/151850/ -> 罪域的骨终为王
//https://www.23us.la/html/209/209550/ -> 文娱万岁
//https://www.23us.la/html/113/113444/ -> 不朽凡人
//顶点小说网 23us.la
type Ebook23US struct {
	Url string
}

func New23US() Ebook23US {
	return Ebook23US{
		Url: "https://www.23us.la",
	}
}

func (this Ebook23US) GetBookInfo(bookid string, proxy string) BookInfo {

	var bi BookInfo
	var volumes []Volume
	var chapters []Chapter
	pollURL := this.Url + "/" + "html/" + handleBookid(bookid) + "/"

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

		//获取书分卷信息
		dtNode, _ := htmlquery.Find(doc, "//dl[@class='chapterlist']//dt") //获取书分卷信息
		testVolStr := htmlquery.InnerText(dtNode[1])
		if TestContainVolume(testVolStr) {
			bi.ChangeVolumeState(true)
			if len(dtNode) == 2 { //就是说刚好两个节点，我们要去除第一个，只保留第二个
				var tmp Volume
				tmp.CurrentVolume = htmlquery.InnerText(dtNode[1])
				volumes = append(volumes, tmp)
			} else { //当len(dtNode) >= 3
				for index := 1; index < len(dtNode); index++ { //因为第一个为 最新章节部分，需要去掉
					var tmp Volume
					//tmp.PrevChapterId =
					PrevChapter, _ := htmlquery.FindOne(dtNode[index], "//preceding-sibling::dd[1]") // 根据当前节点，查找上一个dd节点
					aNode, _ := htmlquery.Find(PrevChapter, "//a")
					tmp.PrevChapter.Link = this.Url + htmlquery.SelectAttr(aNode[0], "href")
					tmp.PrevChapter.Title = htmlquery.InnerText(aNode[0])
					tmp.CurrentVolume = htmlquery.InnerText(dtNode[index])
					volumes = append(volumes, tmp)
				}
			}
			volumes[0].PrevChapterId = 0      //第一分卷，前面的章节，设置为0
			volumes[0].PrevChapter.Link = ""  //第一分卷，前面的章节，连接设置为空
			volumes[0].PrevChapter.Title = "" //第一分卷，前面的章节，标题设置为空
		}
		//获取书章节列表
		ddNode, _ := htmlquery.Find(doc, "//dl[@class='chapterlist']//dd")
		for i := 0; i < len(ddNode); i++ {
			var tmp Chapter
			aNode, _ := htmlquery.Find(ddNode[i], "//a")
			tmp.Link = this.Url + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])
			if bi.HasVolume && len(volumes) >= 2 { //正式写入 PrevChapterId
				for index := 1; index < len(volumes); index++ { //第二个分卷开始，前面就有章节内容了
					if volumes[index].PrevChapter.Link == tmp.Link {
						volumes[index].PrevChapterId = i
					}
				}
			}
			chapters = append(chapters, tmp)
		}

		//导入信息
		bi = BookInfo{
			Name:        bookName,
			Author:      author,
			Description: description,
			Volumes:     volumes,
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

		//获取书分卷信息
		dtNode, _ := htmlquery.Find(doc, "//dl[@class='chapterlist']//dt") //获取书分卷信息
		testVolStr := htmlquery.InnerText(dtNode[1])
		if TestContainVolume(testVolStr) {
			bi.ChangeVolumeState(true)
			if len(dtNode) == 2 { //就是说刚好两个节点，我们要去除第一个，只保留第二个
				var tmp Volume
				tmp.CurrentVolume = htmlquery.InnerText(dtNode[1])
				volumes = append(volumes, tmp)
			} else { //当len(dtNode) >= 3
				for index := 1; index < len(dtNode); index++ { //因为第一个为 最新章节部分，需要去掉
					var tmp Volume
					//tmp.PrevChapterId =
					PrevChapter, _ := htmlquery.FindOne(dtNode[index], "//preceding-sibling::dd[1]") // 根据当前节点，查找上一个dd节点
					aNode, _ := htmlquery.Find(PrevChapter, "//a")
					tmp.PrevChapter.Link = this.Url + htmlquery.SelectAttr(aNode[0], "href")
					tmp.PrevChapter.Title = htmlquery.InnerText(aNode[0])
					tmp.CurrentVolume = htmlquery.InnerText(dtNode[index])
					volumes = append(volumes, tmp)
				}
			}
			volumes[0].PrevChapterId = 0      //第一分卷，前面的章节，设置为0
			volumes[0].PrevChapter.Link = ""  //第一分卷，前面的章节，连接设置为空
			volumes[0].PrevChapter.Title = "" //第一分卷，前面的章节，标题设置为空
		}
		//获取书章节列表
		ddNode, _ := htmlquery.Find(doc, "//dl[@class='chapterlist']//dd")
		for i := 12; i < len(ddNode); i++ { //因为前面的12个ddNode为显示最新的12章，与后面的会重复，所以直接Drop
			var tmp Chapter
			aNode, _ := htmlquery.Find(ddNode[i], "//a")
			tmp.Link = this.Url + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])

			if bi.HasVolume && len(volumes) >= 2 { //正式写入 PrevChapterId
				for index := 1; index < len(volumes); index++ { //第二个分卷开始，前面就有章节内容了
					if volumes[index].PrevChapter.Link == tmp.Link {
						volumes[index].PrevChapterId = i
					}
				}
			}
			chapters = append(chapters, tmp)
		}

		//导入信息
		bi = BookInfo{
			Name:        bookName,
			Author:      author,
			Description: description,
			Volumes:     volumes,
			Chapters:    chapters,
		}
	}
	return bi
}

func (this Ebook23US) GetChapterContent(pc ProxyChapter) Chapter {
	pollURL := pc.C.Link
	proxy := pc.Proxy
	var result Chapter

	if proxy != "" {
		doc, _ := htmlquery.LoadURLWithProxy(pollURL, proxy)
		contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
		contentText := htmlquery.InnerText(contentNode)

		//替换字符串中的特殊字符 \xE3\x80\x80\xE3\x80\x80 为换行符 \n
		tmp := strings.Replace(contentText, "\xE3\x80\x80\xE3\x80\x80", "\r\n", -1)

		//把 readx(); 替换成 ""
		//tmp = strings.Replace(tmp, "999小说更新最快 电脑端:https://www.999xs.com/", "", -1)

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

		//替换字符串中的特殊字符 \xE3\x80\x80\xE3\x80\x80 为换行符 \n
		tmp := strings.Replace(contentText, "\xE3\x80\x80\xE3\x80\x80", "\r\n", -1)

		//把 readx(); 替换成 ""
		//tmp = strings.Replace(tmp, "999小说更新最快 电脑端:https://www.999xs.com/", "", -1)

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

//根据每个章节的 url连接，下载每章对应的内容Content当中
func (this Ebook23US) DownloadChapters(Bi BookInfo, proxy string) BookInfo {
	chapters := Bi.Chapters
	NumChapter := len(chapters)
	ch := make(chan Chapter, 1)
	locker := sync.Mutex{}
	var bar *progressbar.ProgressBar

	sp := pool.NewPool(0, func(obj interface{}) {
		locker.Lock()
		tmp := obj.(ProxyChapter)
		content := this.GetChapterContent(tmp)
		locker.Unlock()
		ch <- content

	})

	go excuteServe(&sp, chapters, proxy)

	//下载章节的时候显示进度条
	bar = progressbar.New(NumChapter)
	bar.RenderBlank()

	for i := 0; i < len(chapters); {
		select {
		case c := <-ch:
			chapters[i].Content = c.Content
			i++
		}
		bar.Add(1)
	}
	sp.Close()

	result := BookInfo{
		Name:        Bi.Name,
		Author:      Bi.Author,
		Description: Bi.Description,
		Chapters:    chapters,
	}

	return result
}

func TestContainVolume(src string) bool {
	return !strings.Contains(src, "正文卷")
}
