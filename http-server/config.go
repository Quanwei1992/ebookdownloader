package main

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	CFG_PATH string = "conf/ebdl_conf.ini"
)

var (
	//Cfg for global use
	Cfg *ini.File
)

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	URL_BASE string `json:"url_base"`
}

var (
	conf Config
)

func init() {

	var err error
	Cfg, err = ini.Load(CFG_PATH)
	if err != nil {
		panic(fmt.Errorf("fail to load config file '%s': %v", CFG_PATH, err))
	}

	conf.Host = Cfg.Section("server").Key("host").MustString("localhost")
	conf.Port = Cfg.Section("server").Key("port").MustString("8080")

	conf.URL_BASE = "http://" + conf.Host + ":" + conf.Port

}
