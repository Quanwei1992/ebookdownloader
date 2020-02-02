// +build windows

package ebookdownloader

import (
	"os/exec"
	"path/filepath"
)

func KindlegenCmd(args ...string) *exec.Cmd {
	path, _ := filepath.Abs("./tools/kindlegen.exe")
	cmd := exec.Command(path, args...)
	return cmd
}
