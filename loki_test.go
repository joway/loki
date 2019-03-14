package loki

import (
	"testing"
)

func TestLogger(t *testing.T) {
	SetLevel(INFO)
	SetFormatter(NewStandardFormatter())

	Info("x: %s", "hi")
	Info("xxx")
	Debug("xxx")
	Error("xxx")
}
