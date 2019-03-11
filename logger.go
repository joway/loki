package logger

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

var (
	logger = &Logger{
		level:     INFO,
		formatter: NewStandardFormatter(),
		handler:   NewConsoleHandler(),
	}

	DEBUG = 1
	INFO  = 2
	ERROR = 3
)

func SetLevel(level int) {
	logger.level = level
}

func SetFormatter(formatter Formatter) {
	logger.formatter = formatter
}

type Logger struct {
	level     int
	formatter Formatter
	handler   Handler
}

func Debug(format string, a ...interface{}) {
	if DEBUG >= logger.level {
		logger.handler.Output(logger.formatter.format(format, a...))
	}
}

func Info(format string, a ...interface{}) {
	if INFO >= logger.level {
		logger.handler.Output(aurora.Blue(logger.formatter.format(format, a...)))
	}
}

func Error(format string, a ...interface{}) {
	if ERROR >= logger.level {
		logger.handler.Output(aurora.Red(logger.formatter.format(format, a...)))
	}
}

type Formatter interface {
	format(format string, a ...interface{}) string
}

type StandardFormatter struct {
	Formatter
}

func NewStandardFormatter() Formatter {
	return StandardFormatter{}
}

func (f StandardFormatter) format(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

type Handler interface {
	Output(output interface{}) error
}

type ConsoleHandler struct {
	Handler
}

func NewConsoleHandler() Handler {
	return ConsoleHandler{}
}

func (c ConsoleHandler) Output(output interface{}) error {
	_, err := fmt.Println(output)
	return err
}
