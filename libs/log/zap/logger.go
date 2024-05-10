package zap

import (
	"fmt"
	"github.com/bpcoder16/Water/logit"
	"go.uber.org/zap"
	"io"
)

var _ logit.Logger = (*Logger)(nil)

type Logger struct {
	log    *zap.Logger
	msgKey string
}

func (l *Logger) Log(level logit.Level, keyValues ...interface{}) error {
	keyValuesLen := len(keyValues)
	if keyValuesLen == 0 || keyValuesLen%2 != 0 {
		l.log.Warn(fmt.Sprint("keyValues must appear in pairs: ", keyValues))
		return nil
	}

	data := make([]zap.Field, 0, (keyValuesLen/2)+1)
	var msg string
	for i := 0; i < keyValuesLen; i += 2 {
		if keyValues[i].(string) == l.msgKey {
			msg, _ = keyValues[i+1].(string)
			continue
		}
		data = append(data, zap.Any(fmt.Sprint(keyValues[i]), keyValues[i+1]))
	}

	switch level {
	case logit.LevelDebug:
		l.log.Debug(msg, data...)
	case logit.LevelInfo:
		l.log.Info(msg, data...)
	case logit.LevelWarn:
		l.log.Warn(msg, data...)
	case logit.LevelError:
		l.log.Error(msg, data...)
	case logit.LevelFatal:
		l.log.Fatal(msg, data...)
	}
	return nil
}

func (l *Logger) Close() {
	_ = l.log.Sync()
}

func NewLogger(w io.Writer) *Logger {
	return &Logger{
		log:    NewZapLogger(w),
		msgKey: logit.DefaultMessageKey,
	}
}

func NewWaterLogger(debugInfoWriter, warnErrorFatalWriter io.Writer) *Logger {
	return &Logger{
		log:    NewWaterZapLogger(debugInfoWriter, warnErrorFatalWriter),
		msgKey: logit.DefaultMessageKey,
	}
}
