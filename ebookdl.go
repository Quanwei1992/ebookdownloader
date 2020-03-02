package ebookdownloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Chain-Zhang/pinyin"
	"github.com/hidapple/isbn-gen/isbn"
	"github.com/unknwon/com"
)

//BookInfo 小说信息
type BookInfo struct {
	EBHost      string    `json:"ebook_host"` //下载小说的网站
	EBookID     string    `json:"ebook_id"`   //对应小说网站的bookid
	BookISBN    string    `json:"isbn"`       //生成一个isbn码
	Name        string    `json:"bookname"`
	Author      string    `json:"author"`
	Description string    `json:"novel_description"`
	IsMobi      bool      `json:"is_mobi"`    //当为true的时候生成mobi
	IsAzw3      bool      `json:"is_azw3"`    //当为true的时候生成azw3,
	HasVolume   bool      `json:"has_volume"` //是否有小说分卷，默认为false；当设置为true的时候，Volumes里面需要包含分卷信息
	Volumes     []Volume  `json:"volumes"`    //小说分卷信息，一般不设置
	Chapters    []Chapter `json:"chapters"`   //小说章节信息
}

//Volume 定义小说分卷信息
type Volume struct {
	PrevChapterID int     `json:"prev_chapter_id"`
	PrevChapter   Chapter `json:"prev_chapter"`
	CurrentVolume string  `json:"current_volume_name"`
	NextChapterID int     `json:"next_chapter_id"`
	NextChapter   Chapter `json:"next_chapter"`
}

// Chapter 定义小说章节信息
type Chapter struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"chapter_url_link"`
}

type ProxyChapter struct {
	Proxy string
	C     Chapter
}

//EBookDLInterface 小说下载器接口interface
type EBookDLInterface interface {
	GetBookInfo(bookid string, proxy string) BookInfo      //获取小说的所有信息，包含小说名，作者，简介等信息
	GetBookBriefInfo(bookid string, proxy string) BookInfo //获取小说最基本的信息，不包含章节信息
	DownloaderChapter(ResultChan chan chan Chapter, pc ProxyChapter, wg *sync.WaitGroup)
	DownloadChapters(Bi BookInfo, proxy string) BookInfo
}

//ReadAllString 读取文件内容，并存入string,最终返回
func ReadAllString(filename string) string {
	filepathAbs, _ := filepath.Abs(filename)
	tmp, _ := ioutil.ReadFile(filepathAbs)
	return string(tmp)
}

//WriteFile 写入文件操作
func WriteFile(filename string, data []byte) error {
	filepathAbs, _ := filepath.Abs(filename)
	os.MkdirAll(path.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filepathAbs, data, 0655)
}

//GenerateISBN GenerateISBN
func (this *BookInfo) GenerateISBN() {
	//bookISBN 设置小说的urn码
	bookISBN, _ := isbn.NewISBN("cn", "")
	bookISBNStr := bookISBN.Number()
	this.BookISBN = bookISBNStr
	//fmt.Printf("bookISBN = %s\n", bookISBNStr)
}

//SetISBN 对小说的ISBN码进行设置
func (this *BookInfo) SetISBN(value string) {
	this.BookISBN = value
}

//ISBN 返回小说的ISBN码
func (this BookInfo) ISBN() string {
	return this.BookISBN
}

//SetKindleEbookType 现在设置，mobi和awz3格式不能同时设置为true
func (this *BookInfo) SetKindleEbookType(isMobi bool, isAzw3 bool) {
	this.IsMobi = isMobi
	this.IsAzw3 = isAzw3
}

// ChangeVolumeState 设置 是否包含分卷信息
func (this *BookInfo) ChangeVolumeState(hasVolume bool) {
	this.HasVolume = hasVolume
}

//VolumeState 返回 HasVolume的状态，true,false
func (this BookInfo) VolumeState() bool {
	return this.HasVolume
}

//Split BookInfo里面的Chapter,以300章为一组进行分割
//当少于300章的里面，全部分为一卷；当有1000卷的时候，分为4卷；
//当分割有n个卷的时候，剩下的章节大于50章，重开一个分卷，当少于50的时候，分割到最后一个分卷里面
func (this BookInfo) Split() []BookInfo {
	chapters := this.Chapters       //小说的章节信息
	volumes := this.Volumes         //小说的分卷信息
	Bookname := this.Name           //小说的名称
	Author := this.Author           //小说的作者
	Description := this.Description //小说的描述信息
	IsAzw3 := this.IsAzw3           //是否生成azw3
	IsMobi := this.IsMobi           //是否生成mobi
	HasVolume := this.HasVolume     //是否包含分卷信息

	var bis []BookInfo //slice of BookInfo

	//当剩下章节大于200章时候
	var tmp1 BookInfo
	var BiggerThan50 bool = false //大于50章的时候设置为 true

	chapterCount := len(chapters)             //有多少章节
	count := (float64)(chapterCount / 300.00) //把章节分成几个部分
	if count < 1 {
		count = math.Ceil(count) //向上取整， 0.8 -> 1
		tmp := BookInfo{
			Name:        Bookname,
			Author:      Author,
			Description: Description,
			IsMobi:      IsMobi,
			IsAzw3:      IsAzw3,
			HasVolume:   HasVolume,
			Volumes:     volumes,
			Chapters:    chapters, //因为少于500章，所以全部分在一起
		}
		bis = append(bis, tmp)
	} else {
		count = math.Floor(count) //向下取下 3.1 -> 3; 2.5 -> 2
		for index := 0; index < (int)(count); index++ {
			tmp := BookInfo{
				Name:        Bookname,
				Author:      Author,
				Description: Description,
				IsMobi:      IsMobi,
				IsAzw3:      IsAzw3,
				HasVolume:   HasVolume,
				Volumes:     volumes,
				Chapters:    chapters[index*300 : (index+1)*300], //chapters[startindex:endindex]
				// [0*500:(0+1)*500] , [1*500:(1+1)*500], [2*500:(2+1)*500]
			}
			if index == (int)(count-1) && ((chapterCount - (index+1)*300) < 50) { //因为count 是向下取整的，所以需要进行一下处理
				tmp.Chapters = chapters[index*300 : chapterCount] // 1680 - 1000 == 680
			} else if index == (int)(count-1) && ((chapterCount - (index+1)*300) > 50) { //当剩下的章节多于200章，重新分割一个新的分卷
				tmp1 = tmp
				tmp1.Chapters = chapters[(index+1)*300 : chapterCount]
				BiggerThan50 = true
			}

			bis = append(bis, tmp)
		}
	}

	if BiggerThan50 { //当最后剩下的章节大于200时，再加多一个分割卷
		bis = append(bis, tmp1)
	}
	fmt.Printf("共分%d个下载单元", len(bis))
	return bis
}

//PrintVolumeInfo 用于打印 小说分卷信息
func (this BookInfo) PrintVolumeInfo() {
	volumes := this.Volumes
	if this.VolumeState() {
		for index := 0; index < len(volumes); index++ {
			fmt.Printf("index = %d\n", index)
			fmt.Printf("PrevChapterID= %d\n", volumes[index].PrevChapterID)
			fmt.Printf("PrevChapter.Title = %s\n", volumes[index].PrevChapter.Title)
			fmt.Printf("CurrentVolume = %s\n", volumes[index].CurrentVolume)
			fmt.Printf("NextChapterID= %d\n", volumes[index].NextChapterID)
			fmt.Printf("NextChapter.Title = %s\n", volumes[index].NextChapter.Title)
		}
	} else {
		fmt.Printf("没有找到本书的分卷信息")
	}
}

//GenerateTxt 生成txt电子书
func (this BookInfo) GenerateTxt() {
	chapters := this.Chapters //小说的章节信息
	volumes := this.Volumes   //小说的分卷信息
	content := ""             //用于存放（分卷、）章节内容
	outfpath := "./outputs/" + this.Name + "-" + this.Author + "/"
	outfname := outfpath + this.Name + "-" + this.Author + ".txt"
	txtAbsPath, _ := filepath.Abs(outfname)
	//当txt文件存在的时候删除它
	if com.IsExist(txtAbsPath) {
		os.RemoveAll(txtAbsPath)
	}
	//创建目录
	os.MkdirAll(filepath.Dir(txtAbsPath)+string(os.PathSeparator), os.ModePerm)
	//创建文件
	fptr, _ := os.Create(txtAbsPath)
	defer fptr.Close()

	for index := 0; index < len(chapters); index++ {
		content = "" //每次循环，都初始化一次
		if len(volumes) > 0 && this.VolumeState() {
			for vindex := 0; vindex < len(volumes); vindex++ {

				if volumes[vindex].PrevChapterID == index {
					//fmt.Printf("volumes[vindex].NextChapterID= %d\n", volumes[vindex].PrevChapterID) //用于测试
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
		fptr.Write(([]byte)(content))                                    //一章一章地往txt文件中写入
		fptr.Sync()                                                      //同步修改的文件
	}

	fptr.Close() //关闭文件
	//WriteFile(outfname, []byte(content))
}

//GenerateJSON 生成json格式的数据
func (this BookInfo) GenerateJSON() error {
	outfpath := "./outputs/" + this.Name + "-" + this.Author + "/"
	outfname := outfpath + this.Name + "-" + this.Author + ".json"
	jsonAbsPath, _ := filepath.Abs(outfname)
	//fmt.Println(jsonAbsPath)
	//当txt文件存在的时候删除它
	if com.IsExist(jsonAbsPath) {
		os.RemoveAll(jsonAbsPath)
	}
	//创建目录
	fmt.Println("jsonpath=", filepath.Dir(jsonAbsPath))
	err := os.MkdirAll(filepath.Dir(jsonAbsPath)+string(os.PathSeparator), os.ModePerm)
	if err != nil {
		return err
	}
	//创建文件
	fptr, err := os.Create(jsonAbsPath)
	if err != nil {
		return err
	}
	defer fptr.Close()

	// 带JSON缩进格式写文件
	data, err := json.MarshalIndent(this, "", "  ")
	if err != nil {
		return err
	}

	//写入文件中
	fptr.Write(data)
	fptr.Sync()
	fptr.Close()

	return nil
}

//GenerateMobi 生成mobi格式电子书
func (this BookInfo) GenerateMobi() {
	chapters := this.Chapters //章节信息
	Volumes := this.Volumes   //分卷信息
	//tpl_cover := ReadAllString("./tpls/tpl_cover.html")
	tplBookToc := ReadAllString("./tpls/tpl_book_toc.html")
	tplChapter := ReadAllString("./tpls/tpl_chapter.html")
	tplVolume := ReadAllString("./tpls/tpl_volume.html")
	tplContent := ReadAllString("./tpls/tpl_content.opf")
	tplStyle := ReadAllString("./tpls/tpl_style.css")
	tplToc := ReadAllString("./tpls/tpl_toc.ncx")
	//将文件名转换成拼音
	strPinyin, _ := pinyin.New(this.Name).Split("-").Mode(pinyin.WithoutTone).Convert()
	savepath := "./tmp/" + strPinyin
	savepath, _ = filepath.Abs(savepath) //使用绝对路径
	if com.IsExist(savepath) {
		os.RemoveAll(savepath)
	}
	os.MkdirAll(path.Dir(savepath)+string(os.PathSeparator), os.ModePerm)

	//设置生成mobi的输出目录
	outputpath := "./outputs/" + this.Name + "-" + this.Author + "/"
	outputpath, _ = filepath.Abs(outputpath)
	outputpath = outputpath + string(os.PathSeparator) //使用绝对路径
	//fmt.Println(outputpath)
	if !com.IsExist(outputpath) {
		os.MkdirAll(outputpath, os.ModePerm)
	}

	// 生成封面
	GenerateCover(this)

	//bookISBN 设置小说的urn码
	bookISBNStr := this.ISBN()
	//fmt.Printf("bookISBN = %s\n", bookISBNStr)

	//cover := strings.Replace(tpl_cover, "___BOOK_NAME___", this.Name, -1)
	//cover = strings.Replace(cover, "___BOOK_AUTHOR___", this.Author, -1)
	//WriteFile(savepath+"/cover.html", []byte(cover))

	//分卷
	if this.VolumeState() && len(Volumes) > 0 {
		for index := 0; index < len(Volumes); index++ {
			vinfo := Volumes[index] //vinfo表示第一分卷信息
			tplVolumeTmp := tplVolume
			volumeid := fmt.Sprintf("Volume%d", index)
			volume := strings.Replace(tplVolumeTmp, "___VOLUME_ID___", volumeid, -1)
			volume = strings.Replace(volume, "___VOLUME_NAME___", vinfo.CurrentVolume, -1)
			cpath := fmt.Sprintf("%s/volume%d.html", savepath, index)
			WriteFile(cpath, []byte(volume))
		}
	}

	// 章节
	tocContent := ""
	naxTocContent := ""
	opfToc := ""
	opfSpine := ""
	tocLine := ""
	naxTocLine := ""
	opfTocLine := ""
	for index := 0; index < len(chapters); index++ {
		// cinfo表示第一个章节的内容
		cinfo := chapters[index]
		tplChapterTmp := tplChapter
		chapterid := fmt.Sprintf("Chapter%d", index)
		//fmt.Printf("Chapterid =%s", chapterid)
		chapter := strings.Replace(tplChapterTmp, "___CHAPTER_ID___", chapterid, -1)
		chapter = strings.Replace(chapter, "___CHAPTER_NAME___", cinfo.Title, -1)
		contentTmp := cinfo.Content
		contentLines := strings.Split(contentTmp, "\r")
		content := ""
		for _, v := range contentLines {
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
				if Volumes[vindex].PrevChapterID == index {
					//分卷信息,在book-toc.html中插入分卷信息
					tocLine = fmt.Sprintf("<dt class=\"tocl1\"><a href=\"volume%d.html\">%s</a></dt>\n", vindex, Volumes[vindex].CurrentVolume)
					tocContent = tocContent + tocLine

					//分卷信息，在toc.ncx中插入分卷信息
					naxTocLine = fmt.Sprintf("<navPoint id=\"volume%d\" playOrder=\"%d\">\n", vindex, vindex+1)
					naxTocContent = naxTocContent + naxTocLine

					naxTocLine = fmt.Sprintf("<navLabel><text>%s</text></navLabel>\n", Volumes[vindex].CurrentVolume)
					naxTocContent = naxTocContent + naxTocLine

					naxTocLine = fmt.Sprintf("<content src=\"volume%d.html\"/>\n</navPoint>\n", vindex)
					naxTocContent = naxTocContent + naxTocLine

					//分卷信息,在content.opf中插入分卷信息
					opfTocLine = fmt.Sprintf("<item id=\"volume%d\" href=\"volume%d.html\" media-type=\"application/xhtml+xml\"/>\n", vindex, vindex)
					opfToc = opfToc + opfTocLine

					opfSpineLine := fmt.Sprintf("<itemref idref=\"volume%d\" linear=\"yes\"/>\n", vindex)
					opfSpine = opfSpine + opfSpineLine
				}
			}
		}
		tocLine = fmt.Sprintf("<dt class=\"tocl2\"><a href=\"chapter%d.html\">%s</a></dt>\n", index, cinfo.Title)
		tocContent = tocContent + tocLine

		// naxToc
		naxTocLine = fmt.Sprintf("<navPoint id=\"chapter%d\" playOrder=\"%d\">\n", index, index+1)
		naxTocContent = naxTocContent + naxTocLine

		naxTocLine = fmt.Sprintf("<navLabel><text>%s</text></navLabel>\n", cinfo.Title)
		naxTocContent = naxTocContent + naxTocLine

		naxTocLine = fmt.Sprintf("<content src=\"chapter%d.html\"/>\n</navPoint>\n", index)
		naxTocContent = naxTocContent + naxTocLine

		opfTocLine = fmt.Sprintf("<item id=\"chapter%d\" href=\"chapter%d.html\" media-type=\"application/xhtml+xml\"/>\n", index, index)
		opfToc = opfToc + opfTocLine

		opfSpineLine := fmt.Sprintf("<itemref idref=\"chapter%d\" linear=\"yes\"/>\n", index)
		opfSpine = opfSpine + opfSpineLine
	}

	// style
	WriteFile(savepath+"/style.css", []byte(tplStyle))

	// 目录
	bookToc := strings.Replace(tplBookToc, "___CONTENT___", tocContent, -1)
	WriteFile(savepath+"/book-toc.html", []byte(bookToc))

	naxToc := strings.Replace(tplToc, "___BOOK_ID___", "11111", -1)
	naxToc = strings.Replace(naxToc, "___BOOK_NAME___", this.Name, -1)
	naxToc = strings.Replace(naxToc, "___BOOK_AUTHOR___", this.Author, -1)
	naxToc = strings.Replace(naxToc, "___NAV___", naxTocContent, -1)
	WriteFile(savepath+"/toc.ncx", []byte(naxToc))

	// opf
	opfContent := strings.Replace(tplContent, "___MANIFEST___", opfToc, -1)
	opfContent = strings.Replace(opfContent, "___SPINE___", opfSpine, -1)
	opfContent = strings.Replace(opfContent, "___BOOK_ID___", "11111", -1)
	opfContent = strings.Replace(opfContent, "___BOOK_NAME___", this.Name, -1)
	opfContent = strings.Replace(opfContent, "___BOOK_AUTHOR___", this.Author, -1)
	opfContent = strings.Replace(opfContent, "__ISBN__", bookISBNStr, -1)
	//设置初始时间
	opfContent = strings.Replace(opfContent, "___CREATE_TIME___", time.Now().Format("2006-01-02 15:04:05"), -1)
	//写入简介信息
	opfContent = strings.Replace(opfContent, "___DESCRIPTION___", this.Description, -1)
	//写入发布者信息
	opfContent = strings.Replace(opfContent, "___PUBLISHER___", "sndnvaps", -1)
	//把修改内容写入到content.opf文件中
	WriteFile(savepath+"/content.opf", []byte(opfContent))

	//把封面复制到 tmp 目录当中
	coverPath, _ := filepath.Abs("./cover.jpg")
	err := com.Copy(coverPath, savepath+"/cover.jpg")
	if err != nil {
		fmt.Println(err.Error())
	}
	//把封面复制到 outputs/小说名-作者/cover.jpg
	err = com.Copy(coverPath, outputpath+"cover.jpg")
	if err != nil {
		fmt.Println(err.Error())
	}
	//删除当前目前的cover.jpg文件
	os.RemoveAll(coverPath)

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
	com.Copy(savepath+string(os.PathSeparator)+outfname, outputpath+string(os.PathSeparator)+outfname)
}

//AsycChapter 同步下载章节的content内容
func AsycChapter(ResultChan chan chan Chapter, chapter chan Chapter) {
	for {
		c := <-ResultChan
		tmp := <-c
		chapter <- tmp
	}

}
