package zap

import (
	"fmt"
	"github.com/bpcoder16/Water/log"
	"go.uber.org/zap"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log    *zap.Logger
	msgKey string
}

func (l *Logger) Log(level log.Level, keyValues ...interface{}) error {
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
	case log.LevelDebug:
		l.log.Debug(msg, data...)
	case log.LevelInfo:
		l.log.Info(msg, data...)
	case log.LevelWarn:
		l.log.Warn(msg, data...)
	case log.LevelError:
		l.log.Error(msg, data...)
	case log.LevelFatal:
		l.log.Fatal(msg, data...)
	}
	return nil
}
