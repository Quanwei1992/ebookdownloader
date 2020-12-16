package ebookdownloader

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Chain-Zhang/pinyin"
	"github.com/unknwon/com"
)

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
	//GenerateCover(this)

	//下载封面
	this.GetCover()

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
