package zap

import (
	"github.com/bpcoder16/Water/logit"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

type EncoderType int8

const (
	JSONEncoder EncoderType = iota
	ConsoleEncoder
)

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "time",
		MessageKey:    logit.DefaultMessageKey,
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format(time.DateTime + ".000"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

func getEncoder(encoderType EncoderType) zapcore.Encoder {
	switch encoderType {
	case JSONEncoder:
		return zapcore.NewJSONEncoder(getEncoderConfig())
	case ConsoleEncoder:
		return zapcore.NewConsoleEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

func NewZapLogger(w io.Writer) *zap.Logger {
	return zap.New(
		zapcore.NewTee(
			zapcore.NewCore(
				getEncoder(JSONEncoder),
				zapcore.AddSync(w),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return true
				}),
			),
		),
	)
}

func NewWaterZapLogger(debugInfoWriter, warnErrorFatalWriter io.Writer) *zap.Logger {
	return zap.New(
		zapcore.NewTee(
			zapcore.NewCore(
				getEncoder(JSONEncoder),
				zapcore.AddSync(debugInfoWriter),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl < zapcore.WarnLevel
				}),
			),
			zapcore.NewCore(
				getEncoder(JSONEncoder),
				zapcore.AddSync(warnErrorFatalWriter),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= zapcore.WarnLevel
				}),
			),
		),
	)
}
