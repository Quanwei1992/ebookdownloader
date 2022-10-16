package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	edl "github.com/sndnvaps/ebookdownloader"
	ebook "github.com/sndnvaps/ebookdownloader/ebook-sources"
	cli "gopkg.in/urfave/cli.v1"
)

var (
	//Version 版本信息
	Version string = "dirty"
	//Commit git commit信息
	Commit string = "d26837e"
	//BuildTime 编译时间
	BuildTime string = "2022-08-21 11:38:21"
)

// EbookDownloader 下载电子书的接口
func EbookDownloader(c *cli.Context) error {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

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
	isEpub := c.Bool("epub")
	isJSON := c.Bool("json")      //把下载到的小说信息保存到json数据当中
	isPV := c.Bool("printvolume") //打印分卷信息，只用做调试时使用

	var bookinfo edl.BookInfo              //初始化变量
	var EBDLInterface edl.EBookDLInterface //初始化接口

	//isTxt 或者 isMobi必须一个为真，或者两个都为真
	if (isTxt || isMobi || isAzw3 || isEpub) ||
		(isTxt && isMobi) ||
		(isTxt && isAzw3) ||
		(isTxt && isEpub) ||
		isPV || isJSON {

		if ebhost == "biqugei.net" {
			xsbiquge := ebook.NewBiqugei()
			EBDLInterface = xsbiquge //实例化接口
		} else if ebhost == "biqugse.com" {
			biqugse := ebook.NewBiqugse()
			EBDLInterface = biqugse
		} else if ebhost == "xixiwx.net" {
			xixiwx := ebook.NewXixiwx()
			EBDLInterface = xixiwx
		} else if ebhost == "zhhbq.com" {
			zzhbq := ebook.NewZhhbq()
			EBDLInterface = zzhbq
		} else {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		}
		// isMobi && isAzw3 当同时为真的时候，退出进程
		if isMobi && isAzw3 {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		}
		bookinfo = EBDLInterface.GetBookInfo(ctx, bookid, proxy)
		//fmt.Println(bookinfo.Chapters) //只用于测试

		//打印分卷信息，只用于调试
		if isPV {
			bookinfo.PrintVolumeInfo()
		} else {
			//下载章节内容
			fmt.Printf("正在下载电子书的相应章节，请耐心等待！\n")
			bookinfo = EBDLInterface.DownloadChapters(bookinfo, proxy)
		}

		if isJSON {
			//生成 json格式后，直接退出程序
			fmt.Printf("\n正在生成json版本的电子书数据，请耐心等待！\n")
			err := bookinfo.GenerateJSON()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		//生成txt文件
		if isTxt {
			fmt.Printf("\n正在生成txt版本的电子书，请耐心等待！\n")
			bookinfo.SetDownloadCoverMethod(false)
			bookinfo.GenerateTxt()
		}
		//生成mobi格式电子书
		if isMobi {
			fmt.Printf("\n正在生成mobi版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(true /* isMobi */, false /* isAzw3 */)
			bookinfo.SetDownloadCoverMethod(true)
			bookinfo.GenerateMobi()

		}
		//生成awz3格式电子书
		if isAzw3 {
			fmt.Printf("\n正在生成Azw3版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(false /* isMobi */, true /* isAzw3 */)
			bookinfo.SetDownloadCoverMethod(true)
			bookinfo.GenerateMobi()
		}

		//生成epub格式电子书
		if isEpub {
			fmt.Printf("\n正在生成EPUB版本的电子书，请耐心等待！\n")
			bookinfo.SetDownloadCoverMethod(true)
			bookinfo.GenerateEPUB()
		}

	} else {
		cli.ShowAppHelpAndExit(c, 0)
		return nil
	}
	fmt.Printf("已经完成生成电子书！\n")

	return nil
}

// ConvJSON2Ebook 转换json文件到ebook格式
func ConvJSON2Ebook(c *cli.Context) error {

	jsonPath := c.String("json")
	if jsonPath == "" {
		cli.ShowAppHelpAndExit(c, 0)
		return nil
	}

	isTxt := c.Bool("txt")
	isMobi := c.Bool("mobi")
	isAzw3 := c.Bool("azw3")
	isEpub := c.Bool("epub")

	//isTxt 或者 isMobi必须一个为真，或者两个都为真
	if (isTxt || isMobi || isAzw3 || isEpub) ||
		(isTxt && isMobi) ||
		(isTxt && isAzw3) ||
		(isTxt && isEpub) {

		// isMobi && isAzw3 当同时为真的时候，退出进程
		// isMobi && isEpub 当同时为真的时候，退出进程
		// isAzw3 && isEpub 当同时为真的时候，退出进程
		// isMobi && isAzw3 && isEpub 当同时为真的时候，退出进程
		if (isMobi && isAzw3) ||
			(isMobi && isEpub) ||
			(isAzw3 && isMobi && isEpub) ||
			(isAzw3 && isEpub) {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		}
		bookinfo, err := edl.LoadBookJSONData(jsonPath)
		if err != nil {
			return err
		}

		//生成txt文件
		if isTxt {
			fmt.Printf("\n正在生成txt版本的电子书，请耐心等待！\n")
			bookinfo.SetDownloadCoverMethod(false)
			bookinfo.GenerateTxt()
		}
		//生成mobi格式电子书
		if isMobi {
			fmt.Printf("\n正在生成mobi版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(true /* isMobi */, false /* isAzw3 */)
			bookinfo.SetDownloadCoverMethod(true)
			bookinfo.GenerateMobi()

		}
		//生成awz3格式电子书
		if isAzw3 {
			fmt.Printf("\n正在生成Azw3版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(false /* isMobi */, true /* isAzw3 */)
			bookinfo.SetDownloadCoverMethod(true)
			bookinfo.GenerateMobi()
		}

		//生成epub格式电子书
		if isEpub {
			fmt.Printf("\n正在生成EPUB版本的电子书，请耐心等待！\n")
			bookinfo.SetDownloadCoverMethod(true)
			bookinfo.GenerateEPUB()
		}

	} else {
		cli.ShowAppHelpAndExit(c, 0)
		return nil
	}
	fmt.Printf("已经完成生成电子书！\n")

	return nil
}

// UpdateCheck 检查更新
func UpdateCheck(*cli.Context) error {
	result, err := edl.UpdateCheck()
	if err == nil {
		CompareResult := result.Compare(Version)
		fmt.Printf("版本检测结果[%s]\n", CompareResult)
		return nil
	}
	return err
}

func main() {

	app := cli.NewApp()
	app.Name = "golang EBookDownloader"
	app.Version = Version + "-" + Commit + "-" + BuildTime
	app.Authors = []cli.Author{
		{
			Name:  "Jimes Yang",
			Email: "sndnvaps@gmail.com",
		},
	}
	app.Copyright = "© 2019 - 2022 Jimes Yang<sndnvaps@gmail.com>"
	app.Usage = "用于下载 笔趣阁(http://www.biqugse.com/,http://www.biqugei.net,https://www.zhhbq.com/) 上面的电子书，并保存为txt格式或者(mobi格式,awz3格式)的电子书"
	app.Action = EbookDownloader
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "ebhost",
			Value: "biqugse.com",
			Usage: "定义下载ebook的网站地址(可选择biqugse.com,biqugei.net,zhhbq.com),西西文学(http://www.xixiwx.net/)",
		},
		cli.StringFlag{
			Name:  "bookid,id",
			Usage: "对应小说网链接最后一串数字,例如：1_1902",
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
			Name:  "epub",
			Usage: "当使用的时候，生成epub文件(不可与--mobi同时使用)",
		},
		cli.BoolFlag{
			Name:  "json",
			Usage: "当使用的时候，把下载得到的小说内容写入到json文件当中",
		},
		cli.BoolFlag{
			Name:  "printvolume,pv",
			Usage: "打印分卷信息，只于调试时使用！(使用此功能的时候，不会下载章节内容)",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "conv",
			Usage: " 转换json格式到其它格式，支持txt,mobi,azw3,epub",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json",
					Usage: "需要转换的json文件",
				},
				cli.BoolFlag{
					Name:  "txt",
					Usage: " 生成txt文件",
				},
				cli.BoolFlag{
					Name:  "mobi",
					Usage: "生成mobi文件",
				},
				cli.BoolFlag{
					Name:  "azw3",
					Usage: "生成azw3文件",
				},
				cli.BoolFlag{
					Name:  "epub",
					Usage: "当使用的时候，生成epub文件",
				},
				cli.BoolFlag{
					Name:  "meta",
					Usage: "生成meta文件",
				},
			},
			Action: ConvJSON2Ebook,
		},
		{
			Name:   "update_check",
			Usage:  "检查更新",
			Action: UpdateCheck,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
