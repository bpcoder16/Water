package logit

import (
	"context"
	"sync"
)

var global = &loggerAppliance{}

type loggerAppliance struct {
	lock   sync.Mutex
	helper *Helper
}

func init() {
	global.SetLogger(DefaultLogger)
}

func GetGlobalHelper() *Helper {
	return global.helper
}

func (a *loggerAppliance) SetLogger(in Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	switch v := in.(type) {
	case *Helper:
		a.helper = v
	default:
		a.helper = NewHelper(v)
	}
}

func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

// Log Print logit by level and keyValues.
func Log(level Level, keyValues ...interface{}) error {
	return global.helper.Log(level, keyValues...)
}

// Context with context logger.
func Context(ctx context.Context) *Helper {
	return global.helper.WithContext(ctx)
}

// Debug logs a message at debug level.
func Debug(a ...interface{}) {
	global.helper.Debug(a...)
}

// DebugF logs a message at debug level.
func DebugF(format string, a ...interface{}) {
	global.helper.DebugF(format, a...)
}

// DebugW logs a message at debug level.
func DebugW(keyValues ...interface{}) {
	global.helper.DebugW(keyValues...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	global.helper.Info(a...)
}

// InfoF logs a message at info level.
func InfoF(format string, a ...interface{}) {
	global.helper.InfoF(format, a...)
}

// InfoW logs a message at info level.
func InfoW(keyValues ...interface{}) {
	global.helper.InfoW(keyValues...)
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	global.helper.Warn(a...)
}

// WarnF logs a message at warn level.
func WarnF(format string, a ...interface{}) {
	global.helper.WarnF(format, a...)
}

// WarnW logs a message at warn level.
func WarnW(keyValues ...interface{}) {
	global.helper.WarnW(keyValues...)
}

// Error logs a message at error level.
func Error(a ...interface{}) {
	global.helper.Error(a...)
}

// ErrorF logs a message at error level.
func ErrorF(format string, a ...interface{}) {
	global.helper.ErrorF(format, a...)
}

// ErrorW logs a message at error level.
func ErrorW(keyValues ...interface{}) {
	global.helper.ErrorW(keyValues...)
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	global.helper.Fatal(a...)
}

// FatalF logs a message at fatal level.
func FatalF(format string, a ...interface{}) {
	global.helper.FatalF(format, a...)
}

// FatalW logs a message at fatal level.
func FatalW(keyValues ...interface{}) {
	global.helper.FatalW(keyValues...)
}
