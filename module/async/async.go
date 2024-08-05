package async

import (
	"context"
	"github.com/bpcoder16/Water/logit"
)

type taskData struct {
	f      func() error
	errMsg string
	cnt    int
}

var fChan = make(chan taskData, 10000)

func AddQueue(f func() error, errMsg string) {
	fChan <- taskData{
		f:      f,
		errMsg: errMsg,
		cnt:    0,
	}
}

func Consumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case f := <-fChan:
			task(f)
		}
	}
}

func task(d taskData) {
	defer func() {
		if r := recover(); r != nil {
			logit.ErrorW("async.task", d.errMsg, "async.task.panic", r)
		}
	}()
	if err := d.f(); err != nil {
		d.cnt++
		if d.cnt >= 3 {
			logit.ErrorW("async.task", d.errMsg, "async.task.err", err)
			return
		}
		logit.WarnW("async.task", d.errMsg, "async.task.err", err)
		fChan <- d
	}
}
