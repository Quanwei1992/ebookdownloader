// +build linux
// +build !windows
// +build !386 !amd64 !arm64
// +build arm

package ebookdownloader

import (
	"os/exec"
	"path/filepath"
)

//KindlegenCmd 执行外部kindlegen命令
func KindlegenCmd(args ...string) *exec.Cmd {
	var cmds []string
	path, _ := filepath.Abs("./tools/kindlegenLinux")
	qemu_path, _ := filepath.Abs("./tools/qemu-i386-static-armhf")
	cmds = append(cmds, path)
	cmds = append(cmds, args...)
	cmd := exec.Command(qemu_path, cmds...)
	return cmd
}
