package main

import (
	"fmt"
	"strings"

	"github.com/andlabs/ui"

	edl "github.com/sndnvaps/ebookdownloader"
	ebook "github.com/sndnvaps/ebookdownloader/ebook-sources"
)

var (
	//Version 版本信息
	Version string = "dev"
	//Commit git commit信息
	Commit string = "7caf59d"
	//BuildTime 编译时间
	BuildTime string = "2020-05-01 20:50"
)

func makeAboutWindow() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	//hbox.Append(vbox, false)

	data1 := fmt.Sprintf("%s\n", "Ebookdownloader")
	data2 := fmt.Sprintf("版本:%s\n", Version)
	data3 := fmt.Sprintf("CommitHash:%s\n", Commit)
	data4 := fmt.Sprintf("编译时间:%s\n", BuildTime)

	announceInfo := "       声明       \n"
	announceInfo += "  本软件用于下载网络上的小说，并生成对应的格式。\n目前支持生成txt,mobi,epub等三种格式。\n"
	announceInfo += "对于本软件在使用过程中造成的问题，一概不负责"
	announceInfoMultiEntry := ui.NewNonWrappingMultilineEntry()

	announceInfoMultiEntry.SetText(announceInfo)
	announceInfoMultiEntry.SetReadOnly(true)

	hbox.Append(ui.NewLabel(data2), false)
	hbox.Append(ui.NewLabel(data3), false)

	vbox.Append(ui.NewLabel(data1), false)
	vbox.Append(hbox, false)
	vbox.Append(ui.NewLabel(data4), false)
	vbox.Append(ui.NewHorizontalSeparator(), false)
	//富文件支持
	vbox.Append(announceInfoMultiEntry, true)

	return vbox
}

//主页
func makeHomeWindow() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hboxID := ui.NewHorizontalBox() //用于设置BookID
	hboxID.SetPadded(true)

	hboxProxy := ui.NewHorizontalBox() //用于设置Proxy
	hboxProxy.SetPadded(true)

	hboxChooseWebsite := ui.NewHorizontalBox() //用于选择下地小说数据的网站
	hboxChooseWebsite.SetPadded(true)

	hboxSaveType := ui.NewHorizontalBox() //用于设置生成小说的格式
	hboxSaveType.SetPadded(true)

	bookIDLabel := ui.NewLabel("BookID")
	bookIDInputEntry := ui.NewEntry()

	hboxID.Append(bookIDLabel, false)
	hboxID.Append(bookIDInputEntry, false)

	proxyLabel := ui.NewLabel("Proxy")
	proxyInputEntry := ui.NewEntry()

	hboxProxy.Append(proxyLabel, false)
	hboxProxy.Append(proxyInputEntry, false)

	fictionWebsiteCombox := ui.NewCombobox()
	fictionWebsiteCombox.Append("biqufan.com")
	fictionWebsiteCombox.Append("booktxt.com")
	fictionWebsiteCombox.Append("899zw.com")
	fictionWebsiteCombox.SetSelected(0) //设置默认选择为 biqufan.com
	fictionWebsiteLabel := ui.NewLabel("请选择要用到的下载源")

	hboxChooseWebsite.Append(fictionWebsiteCombox, false)
	hboxChooseWebsite.Append(fictionWebsiteLabel, false)

	checkboxTxt := ui.NewCheckbox("txt")
	checkboxMobi := ui.NewCheckbox("mobi")
	checkboxEpub := ui.NewCheckbox("epub")

	hboxSaveType.Append(checkboxTxt, false)
	hboxSaveType.Append(checkboxMobi, false)
	hboxSaveType.Append(checkboxEpub, false)

	group := ui.NewGroup("选择要保存的格式")
	group.SetMargined(true)
	//group.SetChild(ui.NewNonWrappingMultilineEntry())
	group.SetChild(hboxSaveType)

	var bookinfo edl.BookInfo              //初始化变量
	var EBDLInterface edl.EBookDLInterface //初始化接口

	runBtn := ui.NewButton("下载")
	runBtn.OnClicked(func(*ui.Button) {
		bookid := ""
		proxy := ""
		if strings.Compare(bookIDInputEntry.Text(), "") != 0 {
			bookid = bookIDInputEntry.Text()
		}

		if strings.Compare(proxyInputEntry.Text(), "") != 0 {
			proxy = proxyInputEntry.Text()
		}
		switch fictionWebsiteCombox.Selected() {
		case 0:
			xsbiquge := ebook.NewXSBiquge()
			EBDLInterface = xsbiquge //实例化接口
		case 1:
			booktxt := ebook.NewBookTXT()
			EBDLInterface = booktxt //实例化接口
		case 2:
			xs999 := ebook.New999XS()
			EBDLInterface = xs999 //实例化接口
		}

		bookinfo = EBDLInterface.GetBookInfo(bookid, proxy)
		bookinfo = EBDLInterface.DownloadChapters(bookinfo, proxy) //下载小说章节内容
		if checkboxTxt.Checked() {                                 //当被选择时，生成txt格式
			bookinfo.GenerateTxt()
		}
		if checkboxMobi.Checked() { //当被选择时，生成mobi格式
			bookinfo.SetKindleEbookType(true /* isMobi */, false /* isAzw3 */)
			bookinfo.GenerateMobi()
		}
		if checkboxEpub.Checked() { //当被选择时，生成epub格式
			bookinfo.GenerateEPUB()
		}

		MsgBoxInfo := fmt.Sprintf("小说名：%s\n作者：%s\n简介：\n\t%s", bookinfo.Name, bookinfo.Author, bookinfo.Description)
		ui.MsgBox(ui.NewWindow("MsgBox", 40, 50, false), "小说已经下载完成", MsgBoxInfo)

	})

	vbox.Append(hboxID, false)
	//分割线
	vbox.Append(ui.NewHorizontalSeparator(), false)
	vbox.Append(hboxProxy, false)
	vbox.Append(ui.NewHorizontalSeparator(), false)
	vbox.Append(hboxChooseWebsite, false)
	vbox.Append(ui.NewHorizontalSeparator(), false)
	//vbox.Append(hboxSaveType, false)
	vbox.Append(group, true)
	vbox.Append(ui.NewHorizontalSeparator(), false)
	vbox.Append(runBtn, false)

	return vbox
}

func InitUI() {
	mainwin := ui.NewWindow("Ebookdownloader", 400, 300, true)
	//关闭
	mainwin.OnClosing(func(*ui.Window) bool {

		ui.Quit()
		return true
	})
	//退出
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return false
	})

	form := ui.NewForm()
	form.SetPadded(true)

	mainwin.SetChild(form)
	mainwin.SetMargined(true)

	//主菜单
	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("主页", makeHomeWindow())
	tab.Append("关于", makeAboutWindow())

	mainwin.Show()
}