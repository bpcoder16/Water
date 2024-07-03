package cron

import (
	"context"
	"github.com/bpcoder16/Water/core/concurrency"
	"github.com/bpcoder16/Water/core/lock/nonblock"
	"github.com/bpcoder16/Water/logit"
	"github.com/bpcoder16/Water/module/redis"
	"strconv"
	"time"
)

type Base struct {
	Ctx               context.Context
	MaxConcurrencyCnt int
	IsRun             bool
	LockName          string

	name                 string
	deadLockExpireSecond time.Duration
	baseTaskList         []func()
	processAddTaskList   []func()
}

func (b *Base) Before(name, lockName string, deadLockExpireSecond time.Duration, maxConcurrencyCnt int) {
	b.Ctx = context.WithValue(context.Background(), logit.DefaultMessageKey, "Cron")
	b.name = name
	b.LockName = lockName
	b.deadLockExpireSecond = deadLockExpireSecond * time.Second
	b.MaxConcurrencyCnt = maxConcurrencyCnt
	b.processAddTaskList = make([]func(), 0, 100)
	b.baseTaskList = make([]func(), 0, 100)
}

func (b *Base) AddBaseTaskList(task func()) {
	b.baseTaskList = append(b.baseTaskList, task)
}

func (b *Base) AddProcessAddTaskList(task func()) {
	b.processAddTaskList = append(b.processAddTaskList, task)
}

func (b *Base) taskPoolRun(taskList []func()) {
	if len(taskList) == 0 {
		return
	}
	taskMap := make(map[string]func(ctx context.Context) concurrency.ChanResult)
	if len(taskList) > b.MaxConcurrencyCnt {
		cnt := 0
		for index, item := range taskList {
			if cnt >= b.MaxConcurrencyCnt {
				_, _ = concurrency.Manager(b.Ctx, taskMap, b.name)
				cnt = 0
				taskMap = make(map[string]func(ctx context.Context) concurrency.ChanResult)
			}
			f := item
			taskMap[strconv.Itoa(index)] = func(ctx context.Context) concurrency.ChanResult {
				f()
				return concurrency.ChanResult{}
			}
			cnt++
		}
		if len(taskMap) > 0 {
			_, _ = concurrency.Manager(b.Ctx, taskMap, b.name)
		}
	} else {
		for index, item := range taskList {
			f := item
			taskMap[strconv.Itoa(index)] = func(ctx context.Context) concurrency.ChanResult {
				f()
				return concurrency.ChanResult{}
			}
		}
		_, _ = concurrency.Manager(b.Ctx, taskMap, b.name)
	}
}

func (b *Base) GetIsRun() bool {
	b.IsRun = nonblock.RedisLock(b.Ctx, b.LockName, b.deadLockExpireSecond)
	return b.IsRun
}

func (b *Base) Init(_ Interface) {
	b.baseTaskList = make([]func(), 0, 100)
}

func (b *Base) Process() {}

func (b *Base) Run() {
	b.taskPoolRun(append(b.baseTaskList, b.processAddTaskList...))
}

func (b *Base) Defer() {
	defer redis.GetDefaultRedis().Del(b.Ctx, b.LockName)
	if r := recover(); r != nil {
		redis.GetDefaultRedis().Del(b.Ctx, b.LockName)
		logit.Context(b.Ctx).ErrorW(b.name+".Err", r)
	} else {
		if b.IsRun {
			logit.Context(b.Ctx).DebugW(b.name+".Status", "Run")
		} else {
			logit.Context(b.Ctx).DebugW(b.name+".Status", "NotRun")
		}
	}
}
