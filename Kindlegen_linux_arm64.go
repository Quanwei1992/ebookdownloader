// +build linux
// +build !386 !amd64 !arm
// +build arm64

package ebookdownloader

import (
	"os/exec"
)

func KindlegenCmd(args ...string) *exec.Cmd {
	var cmds []string
	cmds = append(cmds, "./tools/kindlegenLinux")
	cmds = append(cmds, args...)
	cmd := exec.Command("./tools/qemu-i386-static-arm64", cmds...)
	return cmd
}
