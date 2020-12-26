package main

import (
	"os"

	"github.com/visualfc/goqt/ui"
)

func main() {

	ui.RunEx(os.Args, func() {
		//napp := ui.NewApplication(nil)
		w := NewMainWindow()
		w.mw.Show()

	})
}
