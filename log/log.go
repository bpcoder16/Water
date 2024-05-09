package log

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
	c, ok := l.(*wLogger)
	if !ok {
		return &wLogger{logger: l, prefix: kv, hasValuer: containsValuer(kv), ctx: context.Background()}
	}
	kvs := make([]interface{}, 0, len(c.prefix)+len(kv))
	kvs = append(kvs, c.prefix...)
	kvs = append(kvs, kv...)
	return &wLogger{
		logger:    c.logger,
		prefix:    kvs,
		hasValuer: containsValuer(kvs),
		ctx:       c.ctx,
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
