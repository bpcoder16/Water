package logit

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
			DebugW("logit", msg)
			expected = append(expected, fmt.Sprintf("%s logit=%s", "DEBUG", msg))
		case LevelInfo:
			Info(msg)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "INFO", msg))
			InfoF(testCase.content[0].(string), testCase.content[1:]...)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "INFO", msg))
			InfoW("logit", msg)
			expected = append(expected, fmt.Sprintf("%s logit=%s", "INFO", msg))
		case LevelWarn:
			Warn(msg)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "WARN", msg))
			WarnF(testCase.content[0].(string), testCase.content[1:]...)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "WARN", msg))
			WarnW("logit", msg)
			expected = append(expected, fmt.Sprintf("%s logit=%s", "WARN", msg))
		case LevelError:
			Error(msg)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "ERROR", msg))
			ErrorF(testCase.content[0].(string), testCase.content[1:]...)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "ERROR", msg))
			ErrorW("logit", msg)
			expected = append(expected, fmt.Sprintf("%s logit=%s", "ERROR", msg))
		default:
		}
	}
	_ = Log(LevelInfo, DefaultMessageKey, "test logit")
	expected = append(expected, fmt.Sprintf("%s msg=%s", "INFO", "test logit"))

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
