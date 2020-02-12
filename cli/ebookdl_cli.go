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
	isJson := c.Bool("json")      //把下载到的小说信息保存到json数据当中
	isPV := c.Bool("printvolume") //打印分卷信息，只用做调试时使用
	isMeta := c.Bool("meta")      //保存meta信息到 小说目录当中

	var bookinfo edl.BookInfo              //初始化变量
	var EBDLInterface edl.EBookDLInterface //初始化接口

	var metainfo edl.Meta //用于保存小说的meta信息
	txtfilepath := ""     //定义 txt下载后，获取得到的 地址
	mobifilepath := ""    //定义 mobi下载后，获取得到的 地址
	cover_url_path := ""  //定义下载小说后，封面的url地址

	//isTxt 或者 isMobi必须一个为真，或者两个都为真
	if (isTxt || isMobi || isAzw3) || (isTxt && isMobi) || (isTxt && isAzw3) || isPV || isJson {

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
		} else {
			//下载章节内容
			fmt.Printf("正在下载电子书的相应章节，请耐心等待！\n")
			bookinfo = EBDLInterface.DownloadChapters(bookinfo, proxy)
		}

		if isJson {
			//生成 json格式后，直接退出程序
			err := bookinfo.GenerateJson()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		//生成txt文件
		if isTxt {
			fmt.Printf("\n正在生成txt版本的电子书，请耐心等待！\n")
			bookinfo.GenerateTxt()
			if isMeta { //配置meta信息
				txtfilepath = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + bookinfo.Name + "-" + bookinfo.Author + ".txt"
			}
		}
		//生成mobi格式电子书
		if isMobi {
			fmt.Printf("\n正在生成mobi版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(true /* isMobi */, false /* isAzw3 */)
			bookinfo.GenerateMobi()
			if isMeta { //配置meta信息
				mobifilepath = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + bookinfo.Name + "-" + bookinfo.Author + ".mobi"
				cover_url_path = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + "cover.jpg"
			}

		}
		//生成awz3格式电子书
		if isAzw3 {
			fmt.Printf("\n正在生成Azw3版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(false /* isMobi */, true /* isAzw3 */)
			bookinfo.GenerateMobi()
			if isMeta { //配置meta信息
				mobifilepath = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + bookinfo.Name + "-" + bookinfo.Author + ".azw3"
				cover_url_path = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + "cover.jpg"
			}
		}
		if isMeta {
			metainfo = edl.Meta{
				Ebhost:      ebhost,
				Bookid:      bookid,
				BookName:    bookinfo.Name,
				Author:      bookinfo.Author,
				CoverUrl:    cover_url_path,
				Description: bookinfo.Description,
				TxtUrlPath:  txtfilepath,
				MobiUrlPath: mobifilepath,
			}

			metainfo.WriteFile("./outputs/" + bookinfo.Name + "-" + bookinfo.Author + "/meta.json")
		}

	} else {
		cli.ShowAppHelpAndExit(c, 0)
		return nil
	}
	fmt.Printf("已经完成生成电子书！\n")

	return nil
}

//转换json文件到ebook格式
func ConvJson2Ebook(c *cli.Context) error {

	jsonPath := c.String("json")
	if jsonPath == "" {
		cli.ShowAppHelpAndExit(c, 0)
		return nil
	}

	isTxt := c.Bool("txt")
	isMobi := c.Bool("mobi")
	isAzw3 := c.Bool("azw3")
	isMeta := c.Bool("meta") //保存meta信息到 小说目录当中

	var metainfo edl.Meta //用于保存小说的meta信息
	txtfilepath := ""     //定义 txt下载后，获取得到的 地址
	mobifilepath := ""    //定义 mobi下载后，获取得到的 地址
	cover_url_path := ""  //定义下载小说后，封面的url地址

	//isTxt 或者 isMobi必须一个为真，或者两个都为真
	if (isTxt || isMobi || isAzw3) || (isTxt && isMobi) || (isTxt && isAzw3) {

		// isMobi && isAzw3 当同时为真的时候，退出进程
		if isMobi && isAzw3 {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		}
		bookinfo, err := edl.LoadBookJsonData(jsonPath)
		if err != nil {
			return err
		}

		//生成txt文件
		if isTxt {
			fmt.Printf("\n正在生成txt版本的电子书，请耐心等待！\n")
			bookinfo.GenerateTxt()
			if isMeta { //配置meta信息
				txtfilepath = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + bookinfo.Name + "-" + bookinfo.Author + ".txt"
			}
		}
		//生成mobi格式电子书
		if isMobi {
			fmt.Printf("\n正在生成mobi版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(true /* isMobi */, false /* isAzw3 */)
			bookinfo.GenerateMobi()
			if isMeta { //配置meta信息
				mobifilepath = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + bookinfo.Name + "-" + bookinfo.Author + ".mobi"
				cover_url_path = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + "cover.jpg"
			}

		}
		//生成awz3格式电子书
		if isAzw3 {
			fmt.Printf("\n正在生成Azw3版本的电子书，请耐心等待！\n")
			bookinfo.SetKindleEbookType(false /* isMobi */, true /* isAzw3 */)
			bookinfo.GenerateMobi()
			if isMeta { //配置meta信息
				mobifilepath = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + bookinfo.Name + "-" + bookinfo.Author + ".azw3"
				cover_url_path = "public/" + bookinfo.Name + "-" + bookinfo.Author + "/" + "cover.jpg"
			}
		}
		if isMeta {
			metainfo = edl.Meta{
				//Ebhost:      ebhost,
				//Bookid:      bookid,
				BookName:    bookinfo.Name,
				Author:      bookinfo.Author,
				CoverUrl:    cover_url_path,
				Description: bookinfo.Description,
				TxtUrlPath:  txtfilepath,
				MobiUrlPath: mobifilepath,
			}

			metainfo.WriteFile("./outputs/" + bookinfo.Name + "-" + bookinfo.Author + "/meta.json")
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
	app.Copyright = "© 2019 - 2020 Jimes Yang<sndnvaps@gmail.com>"
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
			Name:  "json",
			Usage: "当使用的时候，把下载得到的小说内容写入到json文件当中",
		},
		cli.BoolFlag{
			Name:  "meta",
			Usage: "把小说的meta信息写入到文件当中",
		},
		cli.BoolFlag{
			Name:  "printvolume,pv",
			Usage: "打印分卷信息，只于调试时使用！(使用此功能的时候，不会下载章节内容)",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "conv",
			Usage: " 转换json格式到其它格式，支持txt,mobi,azw3",
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
					Name:  "meta",
					Usage: "生成meta文件",
				},
			},
			Action: ConvJson2Ebook,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
