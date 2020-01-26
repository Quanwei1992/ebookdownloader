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
	Cfg  *ini.File
	Host string
	Port string
)

func init() {

	var err error
	Cfg, err = ini.Load(CFG_PATH)
	if err != nil {
		panic(fmt.Errorf("fail to load config file '%s': %v", CFG_PATH, err))
	}

	Host = Cfg.Section("server").Key("host").MustString("localhost")
	Port = Cfg.Section("server").Key("port").MustString("8080")

}
