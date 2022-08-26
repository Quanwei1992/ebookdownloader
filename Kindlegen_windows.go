//go:build windows

package ebookdownloader

import (
	_ "embed"
	"fmt"
	"syscall"

	"github.com/amenzhinsky/go-memexec"
)

/* Embedding one file into a slice of bytes */
//go:embed tools/win32/kindlegen.exe
var kindleBinary []byte

// KindlegenCmd 执行外部kindlegen命令
func KindlegenCmd(args ...string) {

	exe, err := memexec.New(kindleBinary)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer exe.Close()

	cmd := exe.Command(args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	cmd.Run()

}
