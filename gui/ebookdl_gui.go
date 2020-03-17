package main

import (
	"fmt"
	"sync"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	edl "github.com/sndnvaps/ebookdownloader"
	ebook "github.com/sndnvaps/ebookdownloader/ebook-sources"
)

var (
	bookname     string   //对应小说名
	bookid       string   //对应小说网的bookid
	proxy        string   //代理，默认为空
	items        []string //定义下载小说的默认网站
	itemSelected int32    //0 -> xsbiquge.com ; 1 -> 999xs.com; 2 -> 23us.la
	checked      bool     //生成txt
	checked2     bool     //生成mobi
	checked3     bool     //生成epub
	multiline    string   //小说简介
	author       string   //小说作者

	compareResult string //版本对比信息

	bookinfo      edl.BookInfo         //初始化bookinfo变量
	ebdlInterface edl.EBookDLInterface //初始化接口
)

var (
	lock sync.Mutex
)

var (
	//Version 版本信息
	Version string = "dev"
	//Commit git commit信息
	Commit string = "b40f73c79"
	//BuildTime 编译时间
	BuildTime string = "2020-02-16 16:34"
)

func btnClickMeClicked() {
	fmt.Println("输入内容为=", bookid)
}

func comboChanged() {
	fmt.Println(items[itemSelected])
}

func btnPopupCLicked() {
	imgui.OpenPopup("Confirm")
}

func btnPopupCLicked1() {
	imgui.OpenPopup("Confirm1")
}

func btnPopupCLicked2() {
	obj, err := edl.UpdateCheck()
	if err == nil {
		compareResult = obj.Compare(Version)
	}
	imgui.OpenPopup("Comfirm2")
}

//EbookDownloaderRun 下载小说功能
func EbookDownloaderRun() {
	multiline = items[itemSelected]
	id := bookid
	ebhost := items[itemSelected]
	isTxt := checked
	isMobi := checked2
	isEpub := checked3
	p := proxy

	var cmdArgs []string //定义命令用到的参数

	switch ebhost {
	case "xsbiquge.com":
		cmdArgs = append(cmdArgs, "--ebhost=xsbiquge.com")
		xsbiquge := ebook.NewXSBiquge()
		ebdlInterface = xsbiquge //实例化接口
	case "biduo.cc":
		cmdArgs = append(cmdArgs, "--ebhost=biduo.cc")
		biduo := ebook.NewBiDuo()
		ebdlInterface = biduo //实例化接口
	case "biquwu.cc":
		cmdArgs = append(cmdArgs, "--ebhost=biquwu.cc")
		biquwu := ebook.NewBiquwu()
		ebdlInterface = biquwu //实例化接口
	case "booktxt.net":
		cmdArgs = append(cmdArgs, "--ebhost=booktxt.net")
		booktxt := ebook.NewBookTXT()
		ebdlInterface = booktxt //实例化接口
	case "999xs.com":
		cmdArgs = append(cmdArgs, "--ebhost=999xs.com")
		xs999 := ebook.New999XS()
		ebdlInterface = xs999 //实例化接口
	case "23us.la":
		cmdArgs = append(cmdArgs, "--ebhost=23us.la")
		xs23 := ebook.New23US()
		ebdlInterface = xs23 //实例化接口
	}

	//add --bookid={{.bookid}}
	cmdArgs = append(cmdArgs, fmt.Sprintf("--bookid=%s", id))

	bookinfo = ebdlInterface.GetBookBriefInfo(id, p)

	bookname = bookinfo.Name
	author = bookinfo.Author
	multiline = bookinfo.Description

	if isTxt {
		cmdArgs = append(cmdArgs, "--txt")
	}
	if isMobi {
		cmdArgs = append(cmdArgs, "--mobi")
	}
	if isEpub {
		cmdArgs = append(cmdArgs, "--epub")
	}
	//添加生成meta.json参数
	cmdArgs = append(cmdArgs, "--meta")

	cmd := EbookdownloaderCliCmd(cmdArgs...)
	lock.Lock()
	cmd.Run()
	lock.Unlock()

}
func loop() {
	// Create main menu bar for master window.
	g.MainMenuBar(
		g.Layout{
			g.Menu("帮助", g.Layout{

				g.Button("关于作者", btnPopupCLicked),
				g.Popup("Confirm", g.Layout{
					g.Label("作者: sndnvaps<sndnvaps@gmail.com>"),
					g.Label("github: https://github.com/sndnvaps"),
					g.Line(
						g.Button("Yes", func() { imgui.CloseCurrentPopup() }),
					),
				}),

				g.Button("关于本软件", btnPopupCLicked1),
				g.Popup("Confirm1", g.Layout{
					g.Label("本软件是基于go语言编写的!"),
					g.Label("项目地址: https://github.com/sndnvaps/ebookdownloader"),
					g.Label("版本号：" + Version),
					g.Label("CommitHash：" + Commit),
					g.Label("编译时间：" + BuildTime),
					g.Line(
						g.Button("Yes", func() { imgui.CloseCurrentPopup() }),
					),
				}),
				g.Button("检查更新", btnPopupCLicked2),
				g.Popup("Confirm2", g.Layout{
					g.Label("检查软件是否需要更新"),
					g.Label("项目地址: https://github.com/sndnvaps/ebookdownloader"),
					g.Label("当前版本号：" + Version),
					g.Label("检测结果：" + compareResult),
					g.Line(
						g.Button("Yes", func() { imgui.CloseCurrentPopup() }),
					),
				}),
			},
			),
		},
	).Build()

	// Build a new window
	size := g.Context.GetPlatform().DisplaySize()
	g.Window("小说下载器@sndnvaps", 0, 20, size[0], size[1], g.Layout{
		g.Label("EBookDownloader"),
		g.Line(
			g.Label("bookid"),
			g.InputText("##bookid", 0, &bookid),
			g.Tooltip("输入对应小说网的bookid"),
		),
		g.Line(
			g.Label("proxy"),
			g.InputText("##proxy", 0, &proxy),
			g.Tooltip("输入代理地址，默认为空"),
		),

		g.Line(
			g.Checkbox("生成txt", &checked, func() {
				fmt.Println(checked)
			}),
			g.Checkbox("生成mobi", &checked2, func() {
				fmt.Println(checked2)
			}),
			g.Checkbox("生成epub", &checked3, func() {
				fmt.Println(checked3)
			}),
		),

		g.Combo("选择要用到的默认下载网站", items[itemSelected], items, &itemSelected, 0, comboChanged),
		g.Line(
			g.Label("小说名"),
			g.InputText("##bookname", 0, &bookname),
		),
		g.Line(
			g.Label("作者"),
			g.InputText("##author", 0, &author),
		),

		g.Line(
			g.Label("简介"),
		),
		g.Line(
			g.InputTextMultiline("##multiline", &multiline, 0, 0, 0, nil, nil),
		),
		g.Line(
			g.Button("下载", EbookDownloaderRun),
			g.Tooltip("点击下载对应网站的小说对应的bookid"),
		),
	})
}

//加载中文字体
func loadFont() {
	fonts := g.Context.IO().Fonts()

	ranges := imgui.NewGlyphRanges()

	builder := imgui.NewFontGlyphRangesBuilder()
	//builder.AddText("铁憨憨你好！")
	builder.AddRanges(fonts.GlyphRangesChineseFull())
	builder.BuildRanges(ranges)

	fontBytes, _ := Asset("fonts/WenQuanYiMicroHei.ttf")
	fonts.AddFontFromMemoryTTFV(fontBytes, 14, imgui.DefaultFontConfig, ranges.Data())
}
func main() {

	//初始化配置文件
	ConfInit()

	items = make([]string, 6)
	//定义items里面的变量
	items[0] = "xsbiquge.com"
	items[1] = "biduo.cc"
	items[2] = "biquwu.cc"
	items[3] = "booktxt.net"
	items[4] = "999xs.com"
	items[5] = "23us.la"

	w := g.NewMasterWindow("EBookDownloader@"+Version, 800, 600, false, loadFont)
	w.Main(loop)
}
