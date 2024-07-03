package cron

import (
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/utils"
)

func init() {
	loadConfig()
}

var config Config

func loadConfig() {
	err := utils.ParseJSONFile(env.RootPath()+"/conf/cron.json", &config)
	if err != nil {
		panic("load cron config err:" + err.Error())
	}
}
