package redis

import (
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/utils"
)

type Config struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	DB         int    `json:"db"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	MaxRetries int    `json:"maxRetries"`
}

var config struct {
	Default *Config `json:"default"`
}

func loadRedisConfig() {
	err := utils.ParseJSONFile(env.RootPath()+"/conf/redis.json", &config)
	if err != nil {
		panic("load Redis config err:" + err.Error())
	}
}
