package loki

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoggerExample(t *testing.T) {
	SetLevel(DEBUG)
	SetFormatter(NewStandardFormatter())
	msg := "Hi, I'm loki! 你好，我是洛基！"

	Debug("msg: %s", msg)
	Info("msg: %s", msg)
	Warn("msg: %s", msg)
	Error("msg: %s", msg)
	//Fatal("msg: %s", msg)
}

func TestLoggerCheck(t *testing.T) {
	//LOKI_ENV=app:xxx
	//root
	assert.True(t, logger.Check())

	l1 := New("app:xxx")
	assert.True(t, l1.Check())

	l2 := New("app:xxx1")
	assert.True(t, l2.Check())

	l3 := New("app:x")
	assert.False(t, l3.Check())
}
