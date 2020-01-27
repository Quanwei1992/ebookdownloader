// +build linux
// +build 386 amd64
// +build !arm !arm64

package ebookdownloader

import (
	"os/exec"
)

func KindlegenCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("./tools/kindlegenLinux", args...)
	return cmd
}
