package main

import (
	"os/exec"
)

//EbookdownloaderCliCmd 传入参数，运行ebookdownloader_cli
func EbookdownloaderCliCmd(args ...string) *exec.Cmd {
	ebookdownloaderCliPath := ebdBinPathConf.Path
	cmd := exec.Command(ebookdownloaderCliPath, args...)
	return cmd
}
