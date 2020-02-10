// +build darwin

package ebookdownloader

import (
	"os/exec"
	"path/filepath"
)

func KindlegenCmd(args ...string) *exec.Cmd {
	path, _ := filepath.Abs("./tools/kindlegenMac")
	cmd := exec.Command(path, args...)
	return cmd
}
