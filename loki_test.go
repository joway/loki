package loki

import (
	"testing"
)

func TestLogger(t *testing.T) {
	SetLevel(DEBUG)
	SetFormatter(NewStandardFormatter())
	msg := "Hi, I'm loki! 你好，我是洛基！"

	Debug("msg: %s", msg)
	Info("msg: %s", msg)
	Warn("msg: %s", msg)
	Error("msg: %s", msg)
}
