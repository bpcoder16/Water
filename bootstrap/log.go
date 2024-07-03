package bootstrap

import (
	"context"
	"github.com/bpcoder16/Water/conf"
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/libs/file/filerotatelogs"
	"github.com/bpcoder16/Water/libs/log/zap"
	"github.com/bpcoder16/Water/logit"
	"github.com/bpcoder16/Water/middlewares"
	"path"
	"time"
)

func initLoggers(_ context.Context, _ *conf.AppConfig) {
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
	logit.SetLogger(
		logit.NewFilter(
			logit.With(
				zap.NewWaterLogger(debugInfoWriter, warnErrorFatalWriter),
				middlewares.LogIdKey,
				func() logit.Valuer {
					return func(ctx context.Context) interface{} {
						logId := ctx.Value(middlewares.LogIdKey)
						if logId == nil {
							return "notSetLogId"
						}
						return logId
					}
				}(),
				logit.DefaultMessageKey,
				func() logit.Valuer {
					return func(ctx context.Context) interface{} {
						msg := ctx.Value(logit.DefaultMessageKey)
						if msg == nil {
							return "Default"
						}
						return msg
					}
				}(),
			),
			logit.FilterLevel(func() logit.Level {
				if env.RunMode() == env.RunModeRelease {
					return logit.LevelInfo
				}
				return logit.LevelDebug
			}()),
		))
}
