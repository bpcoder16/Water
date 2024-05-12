package logit

import (
	"context"
	"fmt"
	"os"
)

// DefaultMessageKey default message key.
var DefaultMessageKey = "type"

// HelperOption is Helper option.
type HelperOption func(*Helper)

// Helper is a logger helper.
type Helper struct {
	logger  Logger
	msgKey  string
	sprint  func(...interface{}) string
	sprintf func(format string, a ...interface{}) string
}

// WithMessageKey with message key.
func WithMessageKey(k string) HelperOption {
	return func(opts *Helper) {
		opts.msgKey = k
	}
}

// WithSprint with sprint
func WithSprint(sprint func(...interface{}) string) HelperOption {
	return func(opts *Helper) {
		opts.sprint = sprint
	}
}

// WithSprintf with sprintf
func WithSprintf(sprintf func(format string, a ...interface{}) string) HelperOption {
	return func(opts *Helper) {
		opts.sprintf = sprintf
	}
}

// NewHelper new a logger helper.
func NewHelper(logger Logger, opts ...HelperOption) *Helper {
	options := &Helper{
		msgKey:  DefaultMessageKey, // default message key
		logger:  logger,
		sprint:  fmt.Sprint,
		sprintf: fmt.Sprintf,
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

// WithContext returns a shallow copy of h with its context changed
// to ctx. The provided ctx must be non-nil.
func (h *Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{
		msgKey:  h.msgKey,
		logger:  WithContext(ctx, h.logger),
		sprint:  h.sprint,
		sprintf: h.sprintf,
	}
}

func (h *Helper) WithValues(kv ...interface{}) *Helper {
	return &Helper{
		msgKey:  h.msgKey,
		logger:  With(h.logger, kv...),
		sprint:  h.sprint,
		sprintf: h.sprintf,
	}
}

// Enabled returns true if the given level above this level.
// It delegates to the underlying *Filter.
func (h *Helper) Enabled(level Level) bool {
	if l, ok := h.logger.(*Filter); ok {
		return level >= l.level
	}
	return true
}

// Log Print logit by level and keyValues
func (h *Helper) Log(level Level, keyValues ...interface{}) error {
	return h.logger.Log(level, keyValues...)
}

// Debug logs a message at debug level.
func (h *Helper) Debug(a ...interface{}) {
	if !h.Enabled(LevelDebug) {
		return
	}
	_ = h.logger.Log(LevelDebug, h.msgKey, h.sprint(a...))
}

// DebugF logs a message at debug level.
func (h *Helper) DebugF(format string, a ...interface{}) {
	if !h.Enabled(LevelDebug) {
		return
	}
	_ = h.logger.Log(LevelDebug, h.msgKey, h.sprintf(format, a...))
}

// DebugW logs a message at debug level.
func (h *Helper) DebugW(keyValues ...interface{}) {
	_ = h.logger.Log(LevelDebug, keyValues...)
}

// Info logs a message at info level.
func (h *Helper) Info(a ...interface{}) {
	if !h.Enabled(LevelInfo) {
		return
	}
	_ = h.logger.Log(LevelInfo, h.msgKey, h.sprint(a...))
}

// InfoF logs a message at info level.
func (h *Helper) InfoF(format string, a ...interface{}) {
	if !h.Enabled(LevelInfo) {
		return
	}
	_ = h.logger.Log(LevelInfo, h.msgKey, h.sprintf(format, a...))
}

// InfoW logs a message at info level.
func (h *Helper) InfoW(keyValues ...interface{}) {
	_ = h.logger.Log(LevelInfo, keyValues...)
}

// Warn logs a message at warn level.
func (h *Helper) Warn(a ...interface{}) {
	if !h.Enabled(LevelWarn) {
		return
	}
	_ = h.logger.Log(LevelWarn, h.msgKey, h.sprint(a...))
}

// WarnF logs a message at warn level.
func (h *Helper) WarnF(format string, a ...interface{}) {
	if !h.Enabled(LevelWarn) {
		return
	}
	_ = h.logger.Log(LevelWarn, h.msgKey, h.sprintf(format, a...))
}

// WarnW logs a message at warn level.
func (h *Helper) WarnW(keyValues ...interface{}) {
	_ = h.logger.Log(LevelWarn, keyValues...)
}

// Error logs a message at error level.
func (h *Helper) Error(a ...interface{}) {
	if !h.Enabled(LevelError) {
		return
	}
	_ = h.logger.Log(LevelError, h.msgKey, h.sprint(a...))
}

// ErrorF logs a message at error level.
func (h *Helper) ErrorF(format string, a ...interface{}) {
	if !h.Enabled(LevelError) {
		return
	}
	_ = h.logger.Log(LevelError, h.msgKey, h.sprintf(format, a...))
}

// ErrorW logs a message at error level.
func (h *Helper) ErrorW(keyValues ...interface{}) {
	_ = h.logger.Log(LevelError, keyValues...)
}

// Fatal logs a message at fatal level.
func (h *Helper) Fatal(a ...interface{}) {
	_ = h.logger.Log(LevelFatal, h.msgKey, h.sprint(a...))
	os.Exit(1)
}

// FatalF logs a message at fatal level.
func (h *Helper) FatalF(format string, a ...interface{}) {
	_ = h.logger.Log(LevelFatal, h.msgKey, h.sprintf(format, a...))
	os.Exit(1)
}

// FatalW logs a message at fatal level.
func (h *Helper) FatalW(keyValues ...interface{}) {
	_ = h.logger.Log(LevelFatal, keyValues...)
	os.Exit(1)
}
