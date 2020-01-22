// +build gui

package main

import (
	"runtime"

	"github.com/AllenDang/gimu"
	"github.com/AllenDang/gimu/nk"
)

var (
	textedit       = gimu.NewTextEdit()
	selected       int
	comboLabel     string
	num1           uint = 11
	num2           uint = 33
	propertyInt    int32
	propertyFloat  float32
	checked        bool
	option         int = 1
	selected1      bool
	selected2      bool
	showPopup      bool
	showAboutPopup bool
	picture        *gimu.Texture
	slider         int32 = 33
	customFont     *nk.Font
)

func msgbox(w *gimu.Window) {
	gimu.SetFont(w.MasterWindow().GetContext(), customFont)

	opened := w.Popup("完成状态", nk.PopupStatic, nk.WindowTitle|nk.WindowNoScrollbar|nk.WindowClosable, nk.NkRect(30, 10, 300, 100), func(w *gimu.Window) {
		w.Row(25).Dynamic(1)
		// Custom font

		w.Label("弹出窗口", "LC")
		buttonState := w.Button("关闭")
		if buttonState {
			showPopup = false
			w.ClosePopup()
		}
		gimu.SetFont(w.MasterWindow().GetContext(), w.MasterWindow().GetDefaultFont())

	})
	if !opened {
		showPopup = false
	}
}

func AboutMsgBox(w *gimu.Window) {
	gimu.SetFont(w.MasterWindow().GetContext(), customFont)

	opened := w.Popup("关于本程序", nk.PopupStatic, nk.WindowTitle|nk.WindowNoScrollbar|nk.WindowClosable, nk.NkRect(30, 10, 300, 100), func(w *gimu.Window) {
		w.Row(25).Dynamic(1)
		// Custom font

		w.Label("本程序由 sndnvaps<sndnvaps@gmail.com>开发", "LC")
		buttonState := w.Button("关闭")
		if buttonState {
			showAboutPopup = false
			w.ClosePopup()
		}
		gimu.SetFont(w.MasterWindow().GetContext(), w.MasterWindow().GetDefaultFont())

	})
	if !opened {
		showAboutPopup = false
	}
}

func inputwidgets(w *gimu.Window) {
	gimu.SetFont(w.MasterWindow().GetContext(), customFont)
	// Menu
	w.Menubar(func(w *gimu.Window) {
		w.Row(25).Static(60, 60, 60)
		// Menu 1
		w.Menu("Menu1", "CC", 200, 100, func(w *gimu.Window) {
			w.Row(25).Dynamic(1)
			w.MenuItemLabel("Menu item 1", "LC")
			w.MenuItemLabel("Menu item 2", "LC")
			w.Button("Button inside menu")
		})
		// Menu 2
		w.Menu("Menu2", "CC", 100, 100, func(w *gimu.Window) {
			w.Row(25).Dynamic(1)
			w.MenuItemLabel("Menu item 1", "LC")
			w.SliderInt(0, &slider, 100, 1)
			w.MenuItemLabel("Menu item 2", "LC")
		})

		// Menu 3
		w.Menu("Help", "CC", 100, 100, func(w *gimu.Window) {
			w.Row(25).Dynamic(1)
			state := w.Button("关于")
			if state {
				showAboutPopup = true
			}
			/*
				if showAboutPopup {
					AboutMsgBox(w)
				}
			*/
		})

	})
	// 弹出窗口
	/*
		w.Row(25).Dynamic(1)
		w.LabelColored("Colored label", color.RGBA{200, 100, 100, 255}, "LC")
		state := w.Button("Click me to show a popup window")
		if state {
			showPopup = true
		}

		if showPopup {
			msgbox(w)
		}
	*/

	w.Row(25).Dynamic(1)
	w.Label("Radio", "LC")
	w.Row(25).Dynamic(3)
	if op1 := w.Radio("Option 1", option == 1); op1 {
		option = 1
	}
	if op2 := w.Radio("Option 2", option == 2); op2 {
		option = 2
	}
	if op3 := w.Radio("Option 3", option == 3); op3 {
		option = 3
	}

	w.Row(25).Static(100, 150)
	w.Label("小说的BookId", "LC")
	textedit.Edit(w, nk.EditField, gimu.EditFilterDefault) //inputbot ,输入框

	//fmt.Println(textedit.GetString())

	w.Row(25).Static(100, 100)
	w.Label("作者", "LC")
	textedit.Edit(w, nk.EditField, gimu.EditFilterDefault) //inputbot ,输入框

	w.Row(25).Static(100)
	w.Label("简介", "LC")
	w.Row(300).Static(600)
	textedit.Edit(w, nk.EditBox, gimu.EditFilterDefault) //inputbot ,输入框

	w.Row(25).Static(100)
	w.Tooltip("点击并进行小说内容的下载")
	w.Button("下载并生成文件")

	gimu.SetFont(w.MasterWindow().GetContext(), w.MasterWindow().GetDefaultFont())
}

func updatefn(w *gimu.Window) {
	width, height := w.MasterWindow().GetSize()
	bounds := nk.NkRect(0, 0, float32(width), float32(height))

	w.Window("小说下载器@sndnvaps", bounds, nk.WindowNoScrollbar, func(w *gimu.Window) {
		_, h := w.MasterWindow().GetSize()
		w.Row(int(float32(h-10)) - 9).Dynamic(1) //只有一个group-1
		w.Group("Group1", nk.WindowBorder|nk.WindowTitle, func(w *gimu.Window) {
			inputwidgets(w)

		})

	})
}

func main() {
	runtime.LockOSThread()

	// Create master window
	wnd := gimu.NewMasterWindow("小说下载器@sndnvaps", 1000, 800, gimu.MasterWindowFlagDefault)
	//Load font
	config := nk.NkFontConfig(14)
	config.SetOversample(1, 1)
	config.SetRange(nk.NkFontChineseGlyphRanges())
	//
	customFont = gimu.LoadFontFromFile("./fonts/FZYTK.TTF", 14, &config)

	wnd.Main(updatefn)
}
