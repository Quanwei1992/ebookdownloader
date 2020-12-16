package ebookdownloader

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/goki/freetype"
)

const (
	fontSize = 40 //字体尺寸
)

//GenerateCover 生成封面 cover.jpg
func GenerateCover(this BookInfo) {

	//需要添加内容的图片
	coverAbs, _ := filepath.Abs("./tpls/cover.jpg")
	//fmt.Println(coverAbs)
	imgfile, err := os.Open(coverAbs)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer imgfile.Close()

	jpgimg, _ := jpeg.Decode(imgfile)
	img := image.NewNRGBA(jpgimg.Bounds())

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, jpgimg.At(x, y))
		}
	}

	//需要一个ttf字体文件
	fontAbs, _ := filepath.Abs("./fonts/WenQuanYiMicroHei.ttf")
	fontBytes, err := ioutil.ReadFile(fontAbs)
	if err != nil {
		log.Println(err.Error())
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err.Error())
	}

	f := freetype.NewContext()
	f.SetDPI(72)
	f.SetFont(font)
	f.SetFontSize(fontSize)
	f.SetClip(jpgimg.Bounds())
	f.SetDst(img)
	f.SetSrc(image.Black) //设置字体颜色

	pt := freetype.Pt(img.Bounds().Dx()-370, img.Bounds().Dy()-590) //字体出现的位置
	//尝试把字符串，坚着写入图片中
	NameRune := []rune(this.Name)
	f.DrawString(string(NameRune[0]), pt) // 第一个中文字符
	for index := 1; index < len(NameRune); index++ {
		pt.Y += f.PointToFixed(50)
		f.DrawString(string(NameRune[index]), pt) //写入 小说名
	}

	f.SetFontSize(35)                                                     //重新设置 字体大小为35
	ptAuthor := freetype.Pt(img.Bounds().Dx()-320, img.Bounds().Dy()-500) //字体出现的位置
	f.DrawString(this.Author+" ©著", ptAuthor)                             //写入小说作者名

	newCoverpath, _ := filepath.Abs("./cover.jpg")
	newfile, err := os.Create(newCoverpath)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer newfile.Close()

	err = jpeg.Encode(newfile, img, &jpeg.Options{Quality: 100})
	if err != nil {
		fmt.Println(err.Error())
	}
}

//DownloadCoverImage 下载小说的封面图片
func (this BookInfo) DownloadCoverImage(coverURL string) error {
	res, err := http.Get(coverURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Printf("封面地址[%s]下载失败，改为直接生成封面!\n", coverURL)
		GenerateCover(this)
		//直接在此处结束进程，返回到上级进程中
		return errors.New("封面下载失败，改为直接生成封面")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("封面地址[%s]下载失败，改为直接生成封面!\n", coverURL)
		GenerateCover(this)
		//直接在此处结束进程，返回到上级进程中
		return err
	}
	//使用
	newCoverpath, _ := filepath.Abs("./cover.jpg")
	ioutil.WriteFile(newCoverpath, body, 0666)

	return nil

}

//GetCover 主要用于从 起点中文网上提取小说的封面
func (this BookInfo) GetCover() error {
	options := []chromedp.ExecAllocatorOption{
		//chromedp.Flag("headless", false), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 20*time.Second)
	defer cancel()

	var nodes []*cdp.Node
	searchLink := "https://www.qidian.com/search?kw=" + this.Name
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(searchLink),
		chromedp.WaitVisible(`div[id="result-list"]`),
		chromedp.Nodes("//div[@id='result-list']//div[@class='book-img-text']//ul//li[1]//div[@class='book-img-box']//a//img", &nodes),
	)
	//当执行出错的时候，优化执行生成封面，再返回错误信息
	if err != nil {
		GenerateCover(this)
		return err
	}
	//当执行出错的时候，优化执行生成封面，再返回错误信息
	if len(nodes) < 1 {
		GenerateCover(this)
		return errors.New("无法获取到封面地址，或者小说名字错误！")
	}
	CoverURL := "https:" + nodes[0].AttributeValue("src")
	//到最后返回nil
	return this.DownloadCoverImage(CoverURL)
}
