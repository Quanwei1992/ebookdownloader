package main

import (
	"fmt"
	"path/filepath"

	"gopkg.in/ini.v1"
)

var (
	//CFGPATH 定义配置文件的路径
	CFGPATH string
)

var (
	//Cfg for global use
	Cfg *ini.File
)

// EBDBinPathConfig ebookdownloader_cli 运行程序路径
type EBDBinPathConfig struct {
	Path string `json:"bin_path"` //ebookdownloader_cli exec path -> /usr/local/bin/ebookdownloader_cli
}

var (
	ebdBinPathConf EBDBinPathConfig
)

//ConfInit 初始化配置文件参数
func ConfInit() {

	var err error
	CFGPATH = "./conf/ebdl_conf.ini"
	cfgAbsPath, _ := filepath.Abs(CFGPATH)
	Cfg, err = ini.Load(cfgAbsPath)
	if err != nil {
		panic(fmt.Errorf("fail to load config file '%s': %v", CFGPATH, err))
	}

	ebdBinPathConf.Path = Cfg.Section("ebookdownloader_cli").Key("path").MustString("./ebookdownloader_cli")
	ebdBinPathConf.Path, _ = filepath.Abs(ebdBinPathConf.Path)
}
