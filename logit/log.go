package logit

import (
	"context"
	"log"
)

var DefaultLogger = NewStdLogger(log.Writer())

type Logger interface {
	Log(level Level, keyValues ...interface{}) error
}

type wLogger struct {
	logger    Logger
	prefix    []interface{}
	hasValuer bool
	ctx       context.Context
}

func (c *wLogger) Log(level Level, keyValues ...interface{}) error {
	kvs := make([]interface{}, 0, len(c.prefix)+len(keyValues))
	kvs = append(kvs, c.prefix...)
	if c.hasValuer {
		bindValues(c.ctx, kvs)
	}
	kvs = append(kvs, keyValues...)
	return c.logger.Log(level, kvs...)
}

// With logger fields.
func With(l Logger, kv ...interface{}) Logger {
	switch v := l.(type) {
	case *Filter:
		fv := *v
		fv.logger = With(fv.logger, kv...)
		return &fv
	case *wLogger:
		kvs := make([]interface{}, 0, len(v.prefix)+len(kv))
		kvs = append(kvs, v.prefix...)
		kvs = append(kvs, kv...)
		nw := *v
		nw.prefix = kvs
		nw.hasValuer = containsValuer(kvs)
		return &nw
	default:
		return &wLogger{logger: l, prefix: kv, hasValuer: containsValuer(kv), ctx: context.Background()}
	}
}

// WithContext returns a shallow copy of l with its context changed
// to ctx. The provided ctx must be non-nil.
func WithContext(ctx context.Context, l Logger) Logger {
	switch v := l.(type) {
	case *wLogger:
		lv := *v
		lv.ctx = ctx
		return &lv
	case *Filter:
		fv := *v
		fv.logger = WithContext(ctx, fv.logger)
		return &fv
	default:
		return &wLogger{logger: l, ctx: ctx}
	}
}
