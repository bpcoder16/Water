package cron

import (
	"errors"
	"github.com/go-co-op/gocron/v2"
	"reflect"
	"sync"
	"time"
)

var cronMap map[string]Interface
var registerCronMu sync.RWMutex

func init() {
	cronMap = make(map[string]Interface)
}

func RegisterCron(cronName string, cron Interface) {
	registerCronMu.Lock()
	cronMap[cronName] = cron
	registerCronMu.Unlock()
}

func getCron(cronConfig ConfigItem) (cron Interface, err error) {
	if len(cronMap) == 0 {
		err = errors.New("cron config list is empty")
		return
	}
	var exist bool
	var cronTemplate Interface
	registerCronMu.RLock()
	cronTemplate, exist = cronMap[cronConfig.Name]
	registerCronMu.RUnlock()
	if !exist {
		err = errors.New("cron config[" + cronConfig.Name + "] is not exist")
		return
	}
	cron, _ = reflect.New(reflect.TypeOf(cronTemplate).Elem()).Interface().(Interface)
	cron.Init(cronTemplate)
	return
}

func Run() {
	if !config.IsRunCron {
		return
	}
	c, errC := gocron.NewScheduler()
	if errC != nil {
		panic(errC)
	}
	for _, cronConfig := range config.CronList {
		cronController, err := getCron(cronConfig)
		cronConfigNew := cronConfig
		if err == nil {
			_, _ = c.NewJob(
				gocron.DurationJob(cronConfigNew.EveryMillisecond*time.Millisecond),
				gocron.NewTask(func() {
					cronController.Before(cronConfigNew.Name, config.LockPreName+":"+cronConfigNew.Name, cronConfigNew.DeadLockExpireSecond, cronConfigNew.MaxConcurrencyCnt)
					defer cronController.Defer()
					if cronController.GetIsRun() {
						cronController.Process()
						cronController.Run()
					}
				}),
			)
		}
	}
	c.Start()
}
