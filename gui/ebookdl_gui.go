package main

import (
	"fmt"
	"time"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	edl "github.com/sndnvaps/ebookdownloader"
)

var (
	bookid       string   //对应小说网的bookid
	proxy        string   //代理，默认为空
	items        []string //定义下载小说的默认网站
	itemSelected int32    //0 -> xsbiquge.com ; 1 -> 999xs.com; 2 -> 23us.la
	checked      bool     //生成txt
	checked2     bool     //生成mobi
	multiline    string   //小说简介
	author       string   //小说作者

	bookinfo      edl.BookInfo         //初始化bookinfo变量
	EBDLInterface edl.EBookDLInterface //初始化接口
)

var (
	Version   string = "v1.6.4"
	Commit    string = ""
	BuildTime string = ""
)

func init() {
	BuildTime = fmt.Sprintf("%s", time.Now().Format("2006/01/02 15:04:05"))
}

func btnClickMeClicked() {
	fmt.Println("输入内容为=", bookid)
}

func comboChanged() {
	fmt.Println(items[itemSelected])
}

func contextMenu1Clicked() {
	fmt.Println("Context menu 1 is clicked")
}

func contextMenu2Clicked() {
	fmt.Println("Context menu 2 is clicked")
}

func btnPopupCLicked() {
	imgui.OpenPopup("Confirm")
}

func btnPopupCLicked1() {
	imgui.OpenPopup("Confirm1")
}
func EbookDownloaderRun() {
	multiline = items[itemSelected]
	id := bookid
	ebhost := items[itemSelected]
	isTxt := checked
	isMobi := checked2
	p := proxy

	switch ebhost {
	case "xsbiquge.com":
		xsbiquge := edl.NewXSBiquge()
		EBDLInterface = xsbiquge //实例化接口
	case "999xs.com":
		xs999 := edl.New999XS()
		EBDLInterface = xs999 //实例化接口
	case "23us.la":
		xs23 := edl.New23US()
		EBDLInterface = xs23 //实例化接口
	}

	bookinfo = EBDLInterface.GetBookInfo(id, p)

	author = bookinfo.Author
	multiline = bookinfo.Description

	bookinfo = EBDLInterface.DownloadChapters(bookinfo, proxy)

	if isTxt {
		bookinfo.GenerateTxt()
	}
	if isMobi {
		bookinfo.SetKindleEbookType(true /* isMobi */, false /* isAzw3 */)
		bookinfo.GenerateMobi()
	}

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
						g.Button("No", nil),
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
						g.Button("No", nil),
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
		),

		g.Combo("选择要用到的默认下载网站", items[itemSelected], items, &itemSelected, 0, comboChanged),
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

	fontPath := "./fonts/fzytk.ttf"
	fonts.AddFontFromFileTTFV(fontPath, 14, imgui.DefaultFontConfig, ranges.Data())
}
func main() {
	items = make([]string, 3)
	//定义items里面的变量
	items[0] = "xsbiquge.com"
	items[1] = "999xs.com"
	items[2] = "23us.la"

	w := g.NewMasterWindow("EBookDownloader@"+Version, 800, 600, false, loadFont)
	w.Main(loop)
}
