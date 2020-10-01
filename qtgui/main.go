package main

import (
	"os"

	"github.com/visualfc/goqt/ui"
)

func main() {
	ui.RunEx(os.Args, func() {
		w := NewMainWindow()
		w.mw.Show()
	})
}
