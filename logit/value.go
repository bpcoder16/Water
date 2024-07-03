package logit

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	// DefaultCaller is a Valuer that returns the file and line.
	DefaultCaller = Caller(4)

	// DefaultTimestamp is a Valuer that returns the current wallclock time.
	DefaultTimestamp = Timestamp(time.RFC3339)
)

// Valuer is returns a logit value.
type Valuer func(ctx context.Context) interface{}

// Value return the function value.
func Value(ctx context.Context, v interface{}) interface{} {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}

// Caller returns a Valuer that returns a pkg/file:line description of the caller.
func Caller(depth int) Valuer {
	return func(_ context.Context) interface{} {
		_, file, line, _ := runtime.Caller(depth)
		idx := strings.LastIndexByte(file, '/')
		if idx == -1 {
			return file[idx+1:] + ":" + strconv.Itoa(line)
		}
		idx = strings.LastIndexByte(file[:idx], '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

// Timestamp returns a timestamp Valuer with a custom time format.
func Timestamp(layout string) Valuer {
	return func(_ context.Context) interface{} {
		return time.Now().Format(layout)
	}
}

func bindValues(ctx context.Context, keyValues []interface{}) {
	for i := 1; i < len(keyValues); i += 2 {
		if v, ok := keyValues[i].(Valuer); ok {
			keyValues[i] = v(ctx)
		}
	}
}

func containsValuer(keyValues []interface{}) bool {
	for i := 1; i < len(keyValues); i += 2 {
		if _, ok := keyValues[i].(Valuer); ok {
			return true
		}
	}
	return false
}
