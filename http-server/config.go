package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"path/filepath"
)

var (
	CFG_PATH string
)

var (
	//Cfg for global use
	Cfg *ini.File
)

type Config struct {
	Host     string `json:"host"`      //对外地址
	InerHost string `json:"iner_host"` //内部地址
	Port     string `json:"port"`
	URL_BASE string `json:"url_base"`
}

type KalaConfig struct {
	Host     string `json:"host"`     //kala job Schemule host
	Port     string `json:"port"`     //kala port
	URL_BASE string `json:"url_base"` //url_base
}

var (
	conf  Config
	kconf KalaConfig
)

func ConfInit() {

	var err error
	cfg_abs_path, _ := filepath.Abs(CFG_PATH)
	Cfg, err = ini.Load(cfg_abs_path)
	if err != nil {
		panic(fmt.Errorf("fail to load config file '%s': %v", CFG_PATH, err))
	}

	conf.Host = Cfg.Section("server").Key("host").MustString("localhost")
	conf.InerHost = Cfg.Section("server").Key("iner_host").MustString("localhost")
	conf.Port = Cfg.Section("server").Key("port").MustString("8080")

	conf.URL_BASE = "http://" + conf.Host + ":" + conf.Port

	kconf.Host = Cfg.Section("kala").Key("host").MustString("localhost")
	kconf.Port = Cfg.Section("kala").Key("port").MustString("8081")
	kconf.URL_BASE = "http://" + kconf.Host + ":" + kconf.Port

}
