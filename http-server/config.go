package main

import (
	"fmt"
	"path/filepath"

	"gopkg.in/ini.v1"
)

var (
	CFG_PATH string
)

var (
	//Cfg for global use
	Cfg *ini.File
)

//Config 定义http-server服务器的相关参数
type Config struct {
	Host     string `json:"host"`      //对外地址
	InerHost string `json:"iner_host"` //内部地址
	Port     string `json:"port"`
	URLBase  string `json:"url_base"`
}

// KalaConfig 定义 kala程序运行的相关参数
type KalaConfig struct {
	Host    string `json:"host"`     //kala job Schemule host
	Port    string `json:"port"`     //kala port
	URLBase string `json:"url_base"` //url_base
}

// EBDBinPathConfig ebookdownloader_cli 运行程序路径
type EBDBinPathConfig struct {
	Path string `json:"bin_path"` //ebookdownloader_cli exec path -> /usr/local/bin/ebookdownloader_cli
}

var (
	conf           Config
	kconf          KalaConfig
	ebdBinPathConf EBDBinPathConfig
)

//ConfInit 初始化配置文件参数
func ConfInit() {

	var err error
	cfgAbsPath, _ := filepath.Abs(CFG_PATH)
	Cfg, err = ini.Load(cfgAbsPath)
	if err != nil {
		panic(fmt.Errorf("fail to load config file '%s': %v", CFG_PATH, err))
	}

	conf.Host = Cfg.Section("server").Key("host").MustString("localhost")
	conf.InerHost = Cfg.Section("server").Key("iner_host").MustString("localhost")
	conf.Port = Cfg.Section("server").Key("port").MustString("8080")

	conf.URLBase = "http://" + conf.Host + ":" + conf.Port

	kconf.Host = Cfg.Section("kala").Key("host").MustString("localhost")
	kconf.Port = Cfg.Section("kala").Key("port").MustString("8081")
	kconf.URLBase = "http://" + kconf.Host + ":" + kconf.Port

	ebdBinPathConf.Path = Cfg.Section("ebookdownloader_cli").Key("path").MustString("./ebookdownloader_cli")
	ebdBinPathConf.Path, _ = filepath.Abs(ebdBinPathConf.Path)
}
