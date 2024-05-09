package log

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestGlobalLog(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := NewStdLogger(buf)
	SetLogger(logger)

	testCases := []struct {
		level   Level
		content []interface{}
	}{
		{
			LevelDebug,
			[]interface{}{"test debug"},
		},
		{
			LevelInfo,
			[]interface{}{"test info"},
		},
		{
			LevelInfo,
			[]interface{}{"test %s", "info"},
		},
		{
			LevelWarn,
			[]interface{}{"test warn"},
		},
		{
			LevelError,
			[]interface{}{"test error"},
		},
		{
			LevelError,
			[]interface{}{"test %s", "error"},
		},
	}

	var expected []string
	for _, testCase := range testCases {
		msg := fmt.Sprintf(testCase.content[0].(string), testCase.content[1:]...)
		switch testCase.level {
		case LevelDebug:
			Debug(msg)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "DEBUG", msg))
			DebugF(testCase.content[0].(string), testCase.content[1:]...)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "DEBUG", msg))
			DebugW("log", msg)
			expected = append(expected, fmt.Sprintf("%s log=%s", "DEBUG", msg))
		case LevelInfo:
			Info(msg)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "INFO", msg))
			InfoF(testCase.content[0].(string), testCase.content[1:]...)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "INFO", msg))
			InfoW("log", msg)
			expected = append(expected, fmt.Sprintf("%s log=%s", "INFO", msg))
		case LevelWarn:
			Warn(msg)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "WARN", msg))
			WarnF(testCase.content[0].(string), testCase.content[1:]...)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "WARN", msg))
			WarnW("log", msg)
			expected = append(expected, fmt.Sprintf("%s log=%s", "WARN", msg))
		case LevelError:
			Error(msg)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "ERROR", msg))
			ErrorF(testCase.content[0].(string), testCase.content[1:]...)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "ERROR", msg))
			ErrorW("log", msg)
			expected = append(expected, fmt.Sprintf("%s log=%s", "ERROR", msg))
		default:
		}
	}
	_ = Log(LevelInfo, DefaultMessageKey, "test log")
	expected = append(expected, fmt.Sprintf("%s msg=%s", "INFO", "test log"))

	expected = append(expected, "")

	t.Logf("Content: %s", buf.String())

	if buf.String() != strings.Join(expected, "\n") {
		t.Errorf("Expected: %s, got: %s", strings.Join(expected, "\n"), buf.String())
	}
}

func TestGlobalContext(t *testing.T) {
	buf := new(bytes.Buffer)
	SetLogger(NewStdLogger(buf))
	Context(context.Background()).InfoF("111")
	if buf.String() != "INFO msg=111\n" {
		t.Errorf("Expected:%s, got:%s", "INFO msg=111", buf.String())
	}
}
