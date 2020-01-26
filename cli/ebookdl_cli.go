package main

import (
	"fmt"
	"log"
	"os"

	edl "github.com/sndnvaps/ebookdownloader"
	"gopkg.in/urfave/cli.v1"
)

var (
	Version   string = "1.6.9"
	Commit    string = ""
	BuildTime string = ""
)

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
	isAzw3 := c.Bool("azw3")
	isPV := c.Bool("printvolume") //打印分卷信息，只用做调试时使用

	var bookinfo edl.BookInfo              //初始化变量
	var EBDLInterface edl.EBookDLInterface //初始化接口
	//isTxt 或者 isMobi必须一个为真，或者两个都为真
	if (isTxt || isMobi || isAzw3) || (isTxt && isMobi) || (isTxt && isAzw3) || isPV {

		if ebhost == "xsbiquge.com" {
			xsbiquge := edl.NewXSBiquge()
			EBDLInterface = xsbiquge //实例化接口
		} else if ebhost == "999xs.com" {
			xs999 := edl.New999XS()
			EBDLInterface = xs999 //实例化接口
		} else if ebhost == "23us.la" {
			xs23 := edl.New23US()
			EBDLInterface = xs23 //实例化接口
		} else {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		}
		// isMobi && isAzw3 当同时为真的时候，退出进程
		if isMobi && isAzw3 {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		}
		bookinfo = EBDLInterface.GetBookInfo(bookid, proxy)

		//打印分卷信息，只用于调试
		if isPV {
			bookinfo.PrintVolumeInfo()
			return nil
		} else {
			//下载章节内容
			fmt.Printf("正在下载电子书的相应章节，请耐心等待！\n")
			bookinfo = EBDLInterface.DownloadChapters(bookinfo, proxy)
		}
		//生成txt文件
		if isTxt {
			fmt.Printf("\n正在生成txt版本的电子书，请耐心等待！\n")
			bookinfo.GenerateTxt()
		}
		//生成mobi格式电子书
		if isMobi {
			fmt.Printf("\n正在生成mobi版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(true /* isMobi */, false /* isAzw3 */)
			bookinfo.GenerateMobi()
		}
		//生成awz3格式电子书
		if isAzw3 {
			fmt.Printf("\n正在生成Azw3版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(false /* isMobi */, true /* isAzw3 */)
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
	app.Version = Version + "-" + Commit + "-" + BuildTime
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jimes Yang",
			Email: "sndnvaps@gmail.com",
		},
	}
	app.Copyright = "(c) 2019 - 2020 Jimes Yang<sndnvaps@gmail.com>"
	app.Usage = "用于下载 笔趣阁(https://www.xsbiquge.com),999小说网(https://www.999xs.com/) ,顶点小说网(https://www.23us.la) 上面的电子书，并保存为txt格式或者(mobi格式,awz3格式)的电子书"
	app.Action = EbookDownloader
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "ebhost",
			Value: "xsbiquge.com",
			Usage: "定义下载ebook的网站地址(可选择xsbiquge.com,999xs.com,23us.la)",
		},
		cli.StringFlag{
			Name:  "bookid,id",
			Usage: "对应笔趣阁id(https://www.xsbiquge.com/0_642/),其中0_642就是book_id;\n对应999小说网id(https://www.999xs.com/files/article/html/0/591/),其中591为book_id;\n对应顶点小说网id(https://www.23us.la/html/113/113444/),其中113444为bookid",
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
			Usage: "当使用的时候，生成mobi文件(不可与--azw3同时使用)",
		},
		cli.BoolFlag{
			Name:  "azw3",
			Usage: "当使用的时候，生成azw3文件(不可与--mobi同时使用)",
		},
		cli.BoolFlag{
			Name:  "printvolume,pv",
			Usage: "打印分卷信息，只于调试时使用！(使用此功能的时候，不会下载章节内容)",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
