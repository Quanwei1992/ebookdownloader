//go:build linux && (386 || amd64) && (!arm || !arm64)
// +build linux
// +build 386 amd64
// +build !arm !arm64

package ebookdownloader

import (
	_ "embed"
	"fmt"

	"github.com/amenzhinsky/go-memexec"
)

/* Embedding one file into a slice of bytes */
//go:embed tools/linux-x86/kindlegenLinux
var kindleBinary []byte

// KindlegenCmd 执行外部kindlegen命令
func KindlegenCmd(args ...string) {

	exe, err := memexec.New(kindleBinary)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer exe.Close()

	cmd := exe.Command(args...)

	cmd.Run()
}
