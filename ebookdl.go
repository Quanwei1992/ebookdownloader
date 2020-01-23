package ebookdownloader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/Chain-Zhang/pinyin"
	"github.com/unknwon/com"
)

type BookInfo struct {
	Name        string
	Author      string
	Description string
	IsMobi      bool      //当为true的时候生成mobi
	IsAzw3      bool      //当为true的时候生成azw3,
	HasVolume   bool      //是否有小说分卷，默认为false；当设置为true的时候，Volumes里面需要包含分卷信息
	Volumes     []Volume  //小说分卷信息，一般不设置
	Chapters    []Chapter //小说章节信息
}

type Volume struct {
	PrevChapterId int
	PrevChapter   Chapter
	CurrentVolume string
	NextChapterId int
	NextChapter   Chapter
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
	DownloaderChapter(ResultChan chan chan Chapter, pc ProxyChapter, wg *sync.WaitGroup)
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
func (this *BookInfo) SetKindleEbookType(isMobi bool, isAzw3 bool) {
	this.IsMobi = isMobi
	this.IsAzw3 = isAzw3
}

//设置 是否包含分卷信息
// func ChangeVolumeState
func (this *BookInfo) ChangeVolumeState(hasVolume bool) {
	this.HasVolume = hasVolume
}

//返回 HasVolume的状态，true,false
func (this BookInfo) VolumeState() bool {
	return this.HasVolume
}

func (this BookInfo) PrintVolumeInfo() {
	volumes := this.Volumes
	if this.VolumeState() {
		for index := 0; index < len(volumes); index++ {
			fmt.Printf("index = %d\n", index)
			fmt.Printf("PrevChapterId= %d\n", volumes[index].PrevChapterId)
			fmt.Printf("PrevChapter.Title = %s\n", volumes[index].PrevChapter.Title)
			fmt.Printf("CurrentVolume = %s\n", volumes[index].CurrentVolume)
			fmt.Printf("NextChapterId= %d\n", volumes[index].NextChapterId)
			fmt.Printf("NextChapter.Title = %s\n", volumes[index].NextChapter.Title)
		}
	} else {
		fmt.Printf("没有找到本书的分卷信息")
	}
}

//生成txt电子书
func (this BookInfo) GenerateTxt() {
	chapters := this.Chapters //小说的章节信息
	volumes := this.Volumes   //小说的分卷信息
	content := ""             //用于存放（分卷、）章节内容
	outfpath := "./outputs/"
	outfname := outfpath + this.Name + "-" + this.Author + ".txt"

	for index := 0; index < len(chapters); index++ {
		if len(volumes) > 0 && this.VolumeState() {
			for vindex := 0; vindex < len(volumes); vindex++ {

				if volumes[vindex].PrevChapterId == index {
					//fmt.Printf("volumes[vindex].NextChapterId= %d\n", volumes[vindex].PrevChapterId) //用于测试
					//fmt.Printf("ChapterIndex =  %d\n", index)                                        //用于测试
					//fmt.Printf("CurrentVolume = %s\n", volumes[vindex].CurrentVolume)                //用于测试
					content += "\n" + "## " + volumes[vindex].CurrentVolume + " ##" + "\n"
				}
			}
		}
		//fmt.Printf("Title = %s\n", chapters[index].Title)                //用于测试
		//fmt.Printf("Content = %s\n", chapters[index].Content)            //用于测试
		content += "\n" + "### " + chapters[index].Title + " ###" + "\n" //每一章的标题，使用 ‘### 第n章 标题 ###’ 为格式
		content += chapters[index].Content + "\n\n"                      //每一章内容的结尾，使用两个换行符

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
	if this.VolumeState() && len(Volumes) > 0 {
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
		if this.VolumeState() && len(Volumes) > 0 {
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

	if !com.IsExist("./outputs/") {
		os.MkdirAll(path.Dir("./outputs/"), os.ModePerm)
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
	if this.IsAzw3 {
		outfname += ".azw3"
	}
	//-dont_append_source ,禁止mobi 文件中附加源文件
	//cmd := exec.Command("./tools/kindlegen.exe", "-dont_append_source", savepath+"/content.opf", "-c1", "-o", outfname)
	cmd := KindlegenCmd("-dont_append_source", savepath+"/content.opf", "-c1", "-o", outfname)
	cmd.Run()

	// 把生成的mobi文件复制到 outputs/目录下面
	com.Copy(savepath+"/"+outfname, "./outputs/"+outfname)
}

//AsycChapter
func AsycChapter(ResultChan chan chan Chapter, chapter chan Chapter) {
	for {
		c := <-ResultChan
		tmp := <-c
		chapter <- tmp
	}

}
