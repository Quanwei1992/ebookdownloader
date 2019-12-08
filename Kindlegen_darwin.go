// +build darwin

package main

import (
	"os/exec"
)

func KindlegenCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("./tools/kindlegenMac", args...)
	return cmd
}
