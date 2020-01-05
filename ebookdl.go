package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/chain-zhang/pinyin"
	"gopkg.in/urfave/cli.v1"
)

type BookInfo struct {
	Name        string
	Author      string
	Description string
	IsMobi      bool      //当为true的时候生成mobi
	IsAwz3      bool      //当为true的时候生成awz3,
	HasVolume   bool      //是否有小说分卷，默认为false；当设置为true的时候，Volumes里面需要包含分卷信息
	Volumes     []Volume  //小说分卷信息，一般不设置
	Chapters    []Chapter //小说章节信息
}

type Volume struct {
	PrevChapterId int
	CurrentVolume string
	NextChapterId int
}
type Chapter struct {
	Title   string
	Content string
	Link    string
}

type ProxyChapter struct {
	Proxy string
	C     Chapter
}

//interface
type EBookDLInterface interface {
	GetBookInfo(bookid string, proxy string) BookInfo //获取小说的所有信息，包含小说名，作者，简介等信息
	GetChapterContent(pc ProxyChapter) Chapter
	DownloadChapters(Bi BookInfo, proxy string) BookInfo
}

//读取文件内容，并存入string,最终返回
func ReadAllString(filename string) string {
	tmp, _ := ioutil.ReadFile(filename)
	return string(tmp)
}

func WriteFile(filename string, data []byte) error {
	os.MkdirAll(path.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filename, data, 0655)
}

//设置生成mobi格式，或者生成awz3格式
//现在设置，mobi和awz3格式不能同时设置为true
func (this *BookInfo) SetKindleEbookType(isMobi bool, isAwz3 bool) {
	this.IsMobi = isMobi
	this.IsAwz3 = isAwz3
}

//设置 是否包含分卷信息
// func ChangeVolumeState
func (this *BookInfo) ChangeVolumeState(hasVolume bool) {
	this.HasVolume = hasVolume
}

//生成txt电子书
func (this BookInfo) GenerateTxt() {
	chapters := this.Chapters //小说的章节信息
	volumes := this.Volumes   //小说的分卷信息
	content := ""             //用于存放（分卷、）章节内容
	outfpath := "./outputs/"
	outfname := outfpath + this.Name + "-" + this.Author + ".txt"
	if len(volumes) > 0 && this.HasVolume {
		for index := 0; index < len(chapters); index++ {
			for vindex := 0; vindex < len(volumes); vindex++ {

				if volumes[vindex].PrevChapterId == index {
					//fmt.Printf("volumes[vindex].PrevChapterId= %d\n", volumes[vindex].PrevChapterId) //用于测试
					//fmt.Printf("ChapterIndex =  %d\n", index)                                        //用于测试
					content += "\n" + "## " + volumes[vindex].CurrentVolume + " ##" + "\n"
				}
			}
			content += "\n" + "### " + chapters[index].Title + " ###" + "\n" //每一章的标题，使用 ‘### 第n章 标题 ###’ 为格式
			content += chapters[index].Content + "\n\n"                      //每一章内容的结尾，使用两个换行符

		}

	} else {
		for index := 0; index < len(chapters); index++ {
			content += "\n" + "### " + chapters[index].Title + " ###" + "\n" //每一章的标题，使用 ‘### 第n章 标题 ###’ 为格式
			content += chapters[index].Content + "\n\n"                      //每一章内容的结尾，使用两个换行符
		}
	}

	WriteFile(outfname, []byte(content))
}

//生成mobi格式电子书
func (this BookInfo) GenerateMobi() {
	chapters := this.Chapters //章节信息
	Volumes := this.Volumes   //分卷信息
	//tpl_cover := ReadAllString("./tpls/tpl_cover.html")
	tpl_book_toc := ReadAllString("./tpls/tpl_book_toc.html")
	tpl_chapter := ReadAllString("./tpls/tpl_chapter.html")
	tpl_volume := ReadAllString("./tpls/tpl_volume.html")
	tpl_content := ReadAllString("./tpls/tpl_content.opf")
	tpl_style := ReadAllString("./tpls/tpl_style.css")
	tpl_toc := ReadAllString("./tpls/tpl_toc.ncx")
	//将文件名转换成拼音
	strPinyin, _ := pinyin.New(this.Name).Split("-").Mode(pinyin.WithoutTone).Convert()
	savepath := "./tmp/" + strPinyin
	if com.IsExist(savepath) {
		os.RemoveAll(savepath)
	}
	os.MkdirAll(path.Dir(savepath), os.ModePerm)

	// 生成封面
	GenerateCover(this)

	//cover := strings.Replace(tpl_cover, "___BOOK_NAME___", this.Name, -1)
	//cover = strings.Replace(cover, "___BOOK_AUTHOR___", this.Author, -1)
	//WriteFile(savepath+"/cover.html", []byte(cover))

	//分卷
	if this.HasVolume && len(Volumes) > 0 {
		for index := 0; index < len(Volumes); index++ {
			vinfo := Volumes[index] //vinfo表示第一分卷信息
			tpl_volume_tmp := tpl_volume
			volumeid := fmt.Sprintf("Volume%d", index)
			volume := strings.Replace(tpl_volume_tmp, "___VOLUME_ID___", volumeid, -1)
			volume = strings.Replace(volume, "___VOLUME_NAME___", vinfo.CurrentVolume, -1)
			cpath := fmt.Sprintf("%s/volume%d.html", savepath, index)
			WriteFile(cpath, []byte(volume))
		}
	}

	// 章节
	toc_content := ""
	nax_toc_content := ""
	opf_toc := ""
	opf_spine := ""
	toc_line := ""
	nax_toc_line := ""
	opf_toc_line := ""
	for index := 0; index < len(chapters); index++ {
		// cinfo表示第一个章节的内容
		cinfo := chapters[index]
		tpl_chapter_tmp := tpl_chapter
		chapterid := fmt.Sprintf("Chapter%d", index)
		//fmt.Printf("Chapterid =%s", chapterid)
		chapter := strings.Replace(tpl_chapter_tmp, "___CHAPTER_ID___", chapterid, -1)
		chapter = strings.Replace(chapter, "___CHAPTER_NAME___", cinfo.Title, -1)
		content_tmp := cinfo.Content
		content_lines := strings.Split(content_tmp, "\r")
		content := ""
		for _, v := range content_lines {
			content = content + fmt.Sprintf("<p class=\"a\">    %s</p>\n", v)
		}
		chapter = strings.Replace(chapter, "___CONTENT___", content, -1)
		cpath := fmt.Sprintf("%s/chapter%d.html", savepath, index)
		//for debug
		//fmt.Printf("cpath=%s", cpath)
		//fmt.Printf("chapter=%s", chapter)

		WriteFile(cpath, []byte(chapter))

		//分卷信息
		if this.HasVolume && len(Volumes) > 0 {
			for vindex := 0; vindex < len(Volumes); vindex++ {
				if Volumes[vindex].PrevChapterId == index {
					//分卷信息,在book-toc.html中插入分卷信息
					toc_line = fmt.Sprintf("<dt class=\"tocl1\"><a href=\"volume%d.html\">%s</a></dt>\n", vindex, Volumes[vindex].CurrentVolume)
					toc_content = toc_content + toc_line

					//分卷信息，在toc.ncx中插入分卷信息
					nax_toc_line = fmt.Sprintf("<navPoint id=\"volume%d\" playOrder=\"%d\">\n", vindex, vindex+1)
					nax_toc_content = nax_toc_content + nax_toc_line

					nax_toc_line = fmt.Sprintf("<navLabel><text>%s</text></navLabel>\n", Volumes[vindex].CurrentVolume)
					nax_toc_content = nax_toc_content + nax_toc_line

					nax_toc_line = fmt.Sprintf("<content src=\"volume%d.html\"/>\n</navPoint>\n", vindex)
					nax_toc_content = nax_toc_content + nax_toc_line

					//分卷信息,在content.opf中插入分卷信息
					opf_toc_line = fmt.Sprintf("<item id=\"volume%d\" href=\"volume%d.html\" media-type=\"application/xhtml+xml\"/>\n", vindex, vindex)
					opf_toc = opf_toc + opf_toc_line

					opf_spine_line := fmt.Sprintf("<itemref idref=\"volume%d\" linear=\"yes\"/>\n", vindex)
					opf_spine = opf_spine + opf_spine_line
				}
			}
		}
		toc_line = fmt.Sprintf("<dt class=\"tocl2\"><a href=\"chapter%d.html\">%s</a></dt>\n", index, cinfo.Title)
		toc_content = toc_content + toc_line

		// nax_toc
		nax_toc_line = fmt.Sprintf("<navPoint id=\"chapter%d\" playOrder=\"%d\">\n", index, index+1)
		nax_toc_content = nax_toc_content + nax_toc_line

		nax_toc_line = fmt.Sprintf("<navLabel><text>%s</text></navLabel>\n", cinfo.Title)
		nax_toc_content = nax_toc_content + nax_toc_line

		nax_toc_line = fmt.Sprintf("<content src=\"chapter%d.html\"/>\n</navPoint>\n", index)
		nax_toc_content = nax_toc_content + nax_toc_line

		opf_toc_line = fmt.Sprintf("<item id=\"chapter%d\" href=\"chapter%d.html\" media-type=\"application/xhtml+xml\"/>\n", index, index)
		opf_toc = opf_toc + opf_toc_line

		opf_spine_line := fmt.Sprintf("<itemref idref=\"chapter%d\" linear=\"yes\"/>\n", index)
		opf_spine = opf_spine + opf_spine_line
	}

	// style
	WriteFile(savepath+"/style.css", []byte(tpl_style))

	// 目录
	book_toc := strings.Replace(tpl_book_toc, "___CONTENT___", toc_content, -1)
	WriteFile(savepath+"/book-toc.html", []byte(book_toc))

	nax_toc := strings.Replace(tpl_toc, "___BOOK_ID___", "11111", -1)
	nax_toc = strings.Replace(nax_toc, "___BOOK_NAME___", this.Name, -1)
	nax_toc = strings.Replace(nax_toc, "___BOOK_AUTHOR___", this.Author, -1)
	nax_toc = strings.Replace(nax_toc, "___NAV___", nax_toc_content, -1)
	WriteFile(savepath+"/toc.ncx", []byte(nax_toc))

	// opf
	opf_content := strings.Replace(tpl_content, "___MANIFEST___", opf_toc, -1)
	opf_content = strings.Replace(opf_content, "___SPINE___", opf_spine, -1)
	opf_content = strings.Replace(opf_content, "___BOOK_ID___", "11111", -1)
	opf_content = strings.Replace(opf_content, "___BOOK_NAME___", this.Name, -1)
	opf_content = strings.Replace(opf_content, "___BOOK_AUTHOR___", this.Author, -1)
	//设置初始时间
	opf_content = strings.Replace(opf_content, "___CREATE_TIME___", time.Now().Format("2006-01-02 15:04:05"), -1)
	//写入简介信息
	opf_content = strings.Replace(opf_content, "___DESCRIPTION___", this.Description, -1)
	//写入发布者信息
	opf_content = strings.Replace(opf_content, "___PUBLISHER___", "sndnvaps", -1)
	//把修改内容写入到content.opf文件中
	WriteFile(savepath+"/content.opf", []byte(opf_content))

	if !com.IsExist("./outputs") {
		os.MkdirAll(path.Dir("./outputs"), os.ModePerm)
	}

	//把封面复制到 tmp 目录当中
	err := com.Copy("cover.jpg", savepath+"/cover.jpg")
	if err != nil {
		fmt.Println(err.Error())
	}
	//删除当前目前的cover.jpg文件
	os.RemoveAll("cover.jpg")

	// 生成
	outfname := this.Name + "-" + this.Author
	if this.IsMobi {
		outfname += ".mobi"
	}
	if this.IsAwz3 {
		outfname += ".awz3"
	}
	//-dont_append_source ,禁止mobi 文件中附加源文件
	//cmd := exec.Command("./tools/kindlegen.exe", "-dont_append_source", savepath+"/content.opf", "-c2", "-o", outfname)
	cmd := KindlegenCmd("-dont_append_source", savepath+"/content.opf", "-c2", "-o", outfname)
	cmd.Run()

	// 把生成的mobi文件复制到 outputs/目录下面
	com.Copy(savepath+"/"+outfname, "./outputs/"+outfname)
}

func EbookDownloader(c *cli.Context) error {
	//bookid := "91_91345" //91_91345, 0_642
	bookid := c.String("bookid")
	if bookid == "" {
		cli.ShowAppHelpAndExit(c, 0)
		return nil
	}
	//对应下载小说的网站，默认值为xsbiquge.com
	ebhost := c.String("ebhost")

	proxy := c.String("proxy")

	isTxt := c.Bool("txt")
	isMobi := c.Bool("mobi")
	isAwz3 := c.Bool("awz3")

	var bookinfo BookInfo              //初始化变量
	var EBDLInterface EBookDLInterface //初始化接口
	//isTxt 或者 isMobi必须一个为真，或者两个都为真
	if (isTxt || isMobi || isAwz3) || (isTxt && isMobi) || (isTxt && isAwz3) {

		if ebhost == "xsbiquge.com" {
			xsbiquge := NewXSBiquge()
			EBDLInterface = xsbiquge //实例化接口
		} else if ebhost == "999xs.com" {
			xs999 := New999XS()
			EBDLInterface = xs999 //实例化接口
		} else {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		}
		// isMobi && isAwz3 当同时为真的时候，退出进程
		if isMobi && isAwz3 {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		}
		bookinfo = EBDLInterface.GetBookInfo(bookid, proxy)

		//下载章节内容
		fmt.Printf("正在下载电子书的相应章节，请耐心等待！\n")
		bookinfo = EBDLInterface.DownloadChapters(bookinfo, proxy)
		//生成txt文件
		if isTxt {
			fmt.Printf("\n正在生成txt版本的电子书，请耐心等待！\n")
			bookinfo.GenerateTxt()
		}
		//生成mobi格式电子书
		if isMobi {
			fmt.Printf("\n正在生成mobi版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(true /* isMobi */, false /* isAwz3 */)
			bookinfo.GenerateMobi()
		}
		//生成awz3格式电子书
		if isAwz3 {
			fmt.Printf("\n正在生成Awz3版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(false /* isMobi */, true /* isAwz3 */)
			bookinfo.GenerateMobi()
		}

	} else {
		cli.ShowAppHelpAndExit(c, 0)
		return nil
	}
	fmt.Printf("已经完成生成电子书！\n")

	return nil
}

func main() {

	app := cli.NewApp()
	app.Name = "golang EBookDownloader"
	app.Compiled = time.Now()
	app.Version = "1.5.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jimes Yang",
			Email: "sndnvaps@gmail.com",
		},
	}
	app.Copyright = "(c) 2019 - 2020 Jimes Yang<sndnvaps@gmail.com>"
	app.Usage = "用于下载 笔趣阁(https://www.xsbiquge.com),999小说网(https://www.999xs.com/) 上面的电子书，并保存为txt格式或者(mobi格式,awz3格式)的电子书"
	app.Action = EbookDownloader
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "ebhost",
			Value: "xsbiquge.com",
			Usage: "定义下载ebook的网站地址(可选择xsbiquge.com,999xs.com)",
		},
		cli.StringFlag{
			Name:  "bookid,id",
			Usage: "对应 笔趣阁id(https://www.xsbiquge.com/0_642/),其中0_642就是book_id;对应999小说网id(https://www.999xs.com/files/article/html/0/591/),其中591为book_id",
		},
		cli.StringFlag{
			Name:  "proxy,p",
			Usage: "ip代理(http://ip:ipport),减少本机ip被小说网站封的可能性",
		},
		cli.BoolFlag{
			Name:  "txt",
			Usage: "当使用的时候，生成txt文件",
		},
		cli.BoolFlag{
			Name:  "mobi",
			Usage: "当使用的时候，生成mobi文件(不可与--awz3同时使用)",
		},
		cli.BoolFlag{
			Name:  "awz3",
			Usage: "当使用的时候，生成awz3文件(不可与--mobi同时使用)",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
