package main

import (
	"fmt"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
)

var (
	bookid       string
	items        []string //定义下载小说的默认网站
	itemSelected int32
	checked      bool
	checked2     bool
	dragInt      int32
	multiline    string
	author       string
	radioOp      int
)

var (
	Version   string = ""
	Commit    string = ""
	BuildTime string = ""
)

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
func MultilineChanged() {
	multiline = items[itemSelected]
}
func loop() {
	// Create main menu bar for master window.
	g.MainMenuBar(
		g.Layout{
			g.Menu("File", g.Layout{
				g.MenuItem("Open"),
				g.MenuItem("Save"),
				// You could add any kind of widget here, not just menu item.
				g.Menu("Save as ...", g.Layout{
					g.MenuItem("Excel file"),
					g.MenuItem("CSV file"),
					g.Button("Button inside menu", nil),
				},
				),
			},
			),

			g.Menu("Help", g.Layout{
				g.Button("关于作者", btnPopupCLicked),
				g.Popup("Confirm", g.Layout{
					g.Label("作者: sndnvaps<sndnvaps@gmail.com>"),
					g.Line(
						g.Button("Yes", func() { imgui.CloseCurrentPopup() }),
						g.Button("No", nil),
					),
				}),

				g.Button("关于本软件", btnPopupCLicked1),
				g.Popup("Confirm1", g.Layout{
					g.Label("本软件是基于go语言编写的!"),
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
			g.Checkbox("生成txt", &checked, func() {
				fmt.Println(checked)
			}),
			g.Checkbox("生成mobi", &checked2, func() {
				fmt.Println(checked2)
			}),
			g.Dummy(30, 0),
			g.RadioButton("xsbiquge.com", radioOp == 0, func() { radioOp = 0 }),
			g.RadioButton("999xs.com", radioOp == 1, func() { radioOp = 1 }),
			g.RadioButton("23us.la", radioOp == 2, func() { radioOp = 2 }),
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
			g.Button("下载", MultilineChanged),
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

	w := g.NewMasterWindow("Overview测试", 800, 600, false, loadFont)
	w.Main(loop)
}
