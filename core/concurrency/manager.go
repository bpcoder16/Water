package concurrency

import (
	"context"
	"github.com/bpcoder16/Water/module/gtask"
	"github.com/bpcoder16/Water/utils"
)

var (
	panicFunc func(interface{})
)

type ChanResult struct {
	Err    error
	Result interface{}

	uniqueName string // 不需要程序自行设置
}

func Init(panicHandler func(interface{})) {
	panicFunc = panicHandler
}

func Manager(ctx context.Context, taskMap map[string]func(ctx context.Context) ChanResult, logField string) (resultMap map[string]ChanResult, err error) {
	defer utils.TimeCostLog(ctx, "concurrency.Manager."+logField)()
	var g *gtask.Group
	g, ctx = gtask.WithContext(ctx)
	chanList := make(chan ChanResult, len(taskMap))
	for uniqueName, f := range taskMap {
		task := f
		uniqueNameNew := uniqueName
		g.Go(func() error {
			defer utils.TimeCostLog(ctx, "concurrency.Manager."+logField+"."+uniqueNameNew)()
			defer func() {
				if r := recover(); r != nil {
					panicFunc(r)
				}
			}()
			taskResult := task(ctx)
			taskResult.uniqueName = uniqueNameNew
			chanList <- taskResult
			return nil
		})
	}
	err = g.Wait()
	close(chanList)
	if err == nil {
		resultMap = make(map[string]ChanResult)
		for dataItem := range chanList {
			resultMap[dataItem.uniqueName] = dataItem
		}
	}
	return
}
