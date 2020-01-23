// +build windows

package ebookdownloader

import (
	"os/exec"
)

func KindlegenCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("./tools/kindlegen.exe", args...)
	return cmd
}
