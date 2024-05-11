package http

import (
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/utils"
)

var config struct {
	Server struct {
		Port string `json:"port"`
	}
}

func loadHttpConfig() {
	err := utils.ParseJSONFile(env.RootPath()+"/conf/http.json", &config)
	if err != nil {
		panic("load Http config err:" + err.Error())
	}
}
