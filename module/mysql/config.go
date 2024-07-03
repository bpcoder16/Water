package mysql

import (
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/utils"
)

type Config struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Database     string `json:"database"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
}

var config struct {
	Master *Config `json:"master"`
}

func loadMySQLConfig() {
	err := utils.ParseJSONFile(env.RootPath()+"/conf/mysql.json", &config)
	if err != nil {
		panic("load MySQL config err:" + err.Error())
	}
}
