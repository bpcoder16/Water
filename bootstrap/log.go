package bootstrap

import (
	"context"
	"github.com/bpcoder16/Water/conf"
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/libs/file/filerotatelogs"
	"github.com/bpcoder16/Water/libs/log/zap"
	"github.com/bpcoder16/Water/logit"
	"github.com/bpcoder16/Water/middlewares"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

func initLoggers(_ context.Context, conf *conf.AppConfig) {
	var debugInfoFilePath, warnErrorFatalFilePath string
	if len(conf.LogPath) > 0 {
		debugInfoFilePath = path.Join(conf.LogPath, env.AppName(), env.AppName()+".log")
		warnErrorFatalFilePath = path.Join(conf.LogPath, env.AppName(), env.AppName()+".wf.log")
	} else {
		debugInfoFilePath = path.Join(env.RootPath(), "log", env.AppName(), env.AppName()+".log")
		warnErrorFatalFilePath = path.Join(env.RootPath(), "log", env.AppName(), env.AppName()+".wf.log")
	}

	var debugInfoWriter, warnErrorFatalWriter io.Writer
	if conf.NotUseRotateLog {
		var err error
		debugInfoWriter, err = openFileForWriting(debugInfoFilePath)
		if err != nil {
			panic("openFileForWriting.debugInfoWriter.Err" + err.Error())
		}
		warnErrorFatalWriter, err = openFileForWriting(warnErrorFatalFilePath)
		if err != nil {
			panic("openFileForWriting.warnErrorFatalWriter.Err" + err.Error())
		}
	} else {
		debugInfoWriter = filerotatelogs.NewWriter(
			debugInfoFilePath,
			time.Duration(86400*30)*time.Second,
			time.Duration(3600)*time.Second,
		)
		warnErrorFatalWriter = filerotatelogs.NewWriter(
			warnErrorFatalFilePath,
			time.Duration(86400*30)*time.Second,
			time.Duration(3600)*time.Second,
		)
	}
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

func openFileForWriting(path string) (io.Writer, error) {
	dir := filepath.Dir(path)
	if errF := os.MkdirAll(dir, 0755); errF != nil {
		panic("Error creating directory:" + errF.Error())
	}
	// 打开文件（写入模式）
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
