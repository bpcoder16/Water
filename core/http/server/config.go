package server

import (
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	loadHttpConfig()
	switch env.RunMode() {
	case env.RunModeRelease:
		gin.SetMode(gin.ReleaseMode)
	case env.RunModeTest:
		gin.SetMode(gin.TestMode)
	case env.RunModeDebug:
		gin.SetMode(gin.DebugMode)
	}
}

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
