package bootstrap

import (
	"github.com/bpcoder16/Water/conf"
	"github.com/bpcoder16/Water/env"
)

// MustLoadAppConfig 加载 app.toml ,若失败，会 panic
func MustLoadAppConfig(appConfigPath string) *conf.AppConfig {
	var config conf.AppConfig
	err := conf.ParseConfig(appConfigPath, &config)
	if err != nil {
		panic("parse app config failed: " + err.Error())
	}
	env.Default = env.New(config.Env)
	return &config
}
