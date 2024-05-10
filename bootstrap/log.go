package bootstrap

import (
	"context"
	"github.com/bpcoder16/Water/conf"
	"github.com/bpcoder16/Water/contrib/file/filerotatelogs"
	"github.com/bpcoder16/Water/contrib/log/zap"
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/logit"
	"path"
	"time"
)

func initLoggers(ctx context.Context, _ *conf.AppConfig) {
	debugInfoWriter := filerotatelogs.NewWriter(
		path.Join(env.RootPath(), "log", env.AppName(), env.AppName()+".log"),
		time.Duration(86400*5)*time.Second,
		time.Duration(3600)*time.Second,
	)
	warnErrorFatalWriter := filerotatelogs.NewWriter(
		path.Join(env.RootPath(), "log", env.AppName(), env.AppName()+".wf.log"),
		time.Duration(86400*5)*time.Second,
		time.Duration(3600)*time.Second,
	)
	logit.SetLogger(logit.With(
		zap.NewWaterLogger(debugInfoWriter, warnErrorFatalWriter),
		logIdKey,
		func() logit.Valuer {
			return func(ctx context.Context) interface{} {
				logId := ctx.Value(logIdKey)
				if logId == nil {
					return "notSetLogId"
				}
				return logId
			}
		}(),
	))
}
