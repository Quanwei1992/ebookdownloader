// +build linux
// +build 386 amd64
// +build !arm !arm64

package ebookdownloader

import (
	"os/exec"
	"path/filepath"
)

func KindlegenCmd(args ...string) *exec.Cmd {
	path, _ := filepath.Abs("./tools/kindlegenLinux")
	cmd := exec.Command(path, args...)
	return cmd
}
