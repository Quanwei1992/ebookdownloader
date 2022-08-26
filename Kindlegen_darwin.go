//go:build darwin
// +build darwin

package ebookdownloader

import (
	"os/exec"
	"path/filepath"
)

// KindlegenCmd 执行外部kindlegen命令
func KindlegenCmd(args ...string) *exec.Cmd {
	path, _ := filepath.Abs("./tools/kindlegenMac")
	cmd := exec.Command(path, args...)
	return cmd
}
