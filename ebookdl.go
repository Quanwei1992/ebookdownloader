package ebookdownloader

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/hidapple/isbn-gen/isbn"
	"github.com/unknwon/com"
)

//BookInfo 小说信息
type BookInfo struct {
	EBHost      string    `json:"ebook_host"`        //下载小说的网站
	EBookID     string    `json:"ebook_id"`          //对应小说网站的bookid
	BookISBN    string    `json:"isbn"`              //生成一个isbn码
	BookUUID    string    `json:"uuid"`              //生成一个uuid码，准备用于boltdb
	Name        string    `json:"bookname"`          //小说名字
	Author      string    `json:"author"`            //小说作者
	Description string    `json:"novel_description"` //小说简介
	CoverURL    string    `json:"cover_url"`         //小说封面图片地址
	IsMobi      bool      `json:"is_mobi"`           //当为true的时候生成mobi
	IsAzw3      bool      `json:"is_azw3"`           //当为true的时候生成azw3,
	HasVolume   bool      `json:"has_volume"`        //是否有小说分卷，默认为false；当设置为true的时候，Volumes里面需要包含分卷信息
	Volumes     []Volume  `json:"volumes"`           //小说分卷信息，一般不设置
	Chapters    []Chapter `json:"chapters"`          //小说章节信息
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
	GetBookInfo(ctx context.Context, bookid string, proxy string) BookInfo //获取小说的所有信息，包含小说名，作者，简介等信息
	GetBookBriefInfo(bookid string, proxy string) BookInfo                 //获取小说最基本的信息，不包含章节信息
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
