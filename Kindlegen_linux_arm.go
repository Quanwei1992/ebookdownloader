// +build linux
// +build !386 !amd64
// +build arm

package ebookdownloader

import (
	"os/exec"
)

func KindlegenCmd(args ...string) *exec.Cmd {
	var cmds []string
        cmds = append(cmds,"./tools/kindlegenLinux")
        cmds = append(cmds,args...)
        cmd := exec.Command("./tools/qemu-i386-static-armhf",cmds...)
	return cmd
}
