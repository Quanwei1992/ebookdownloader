// +build linux
// +build 386 amd64
// +build !arm !arm64

package ebookdownloader

import (
	//"fmt"
	//"os"
	"os/exec"
	//"path/filepath"
)

//KindlegenCmd 执行外部kindlegen命令
func KindlegenCmd(args ...string) *exec.Cmd {
	//path, _ := filepath.Abs("./tools/kindlegenLinux")
	//mt.Println("KindlegenLinux path=",path)
	cmd := exec.Command("kindlegenLinux" ,args...)
	return cmd
}
