package cron

import "time"

type Config struct {
	LockPreName string       `json:"lockPreName"`
	IsRunCron   bool         `json:"isRunCron"`
	CronList    []ConfigItem `json:"cronList"`
}

type ConfigItem struct {
	Name                 string        `json:"name"`
	EveryMillisecond     time.Duration `json:"everyMillisecond"`
	DeadLockExpireSecond time.Duration `json:"deadLockExpireSecond"`
	MaxConcurrencyCnt    int           `json:"maxConcurrencyCnt"`
}
