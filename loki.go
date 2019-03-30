package loki

import (
	"fmt"
	"github.com/gobwas/glob"
	"github.com/logrusorgru/aurora"
	"os"
	"strings"
	"time"
)

var (
	DEBUG = 1
	INFO  = 2
	WARN  = 3
	ERROR = 4

	LoggerRootName = ""
	LoggerEnv      = os.Getenv("LOKI_ENV")

	logger = New(LoggerRootName)
	globs  []glob.Glob
)

func init() {
	if LoggerEnv != "" {
		patterns := strings.Split(LoggerEnv, ",")
		for _, pattern := range patterns {
			globs = append(globs, glob.MustCompile(pattern))
		}
	}
}

func SetLevel(level int) {
	logger.level = level
}

func SetFormatter(formatter Formatter) {
	logger.formatter = formatter
}

type Logger struct {
	name      string
	level     int
	formatter Formatter
	handler   Handler
}

func New(name string) Logger {
	return Logger{
		name:      name,
		level:     INFO,
		formatter: NewStandardFormatter(),
		handler:   NewConsoleHandler(),
	}
}

func (l Logger) Check() bool {
	if l.name == "" {
		return true
	}
	for _, g := range globs {
		matched := g.Match(l.name)
		if matched {
			return true
		}
	}
	return false
}

func (l Logger) Debug(format string, a ...interface{}) {
	if DEBUG >= l.level {
		l.handler.output(l.formatter.format(format, a...))
	}
}

func (l Logger) Info(format string, a ...interface{}) {
	if INFO >= l.level {
		l.handler.output(aurora.Blue(l.formatter.format(format, a...)))
	}
}

func (l Logger) Warn(format string, a ...interface{}) {
	if WARN >= l.level {
		l.handler.output(aurora.Green(l.formatter.format(format, a...)))
	}
}

func (l Logger) Error(format string, a ...interface{}) {
	if ERROR >= l.level {
		l.handler.output(aurora.Red(l.formatter.format(format, a...)))
	}
}

func (l Logger) Fatal(format string, a ...interface{}) {
	Error(format, a...)
	os.Exit(1)
}

func Debug(format string, a ...interface{}) {
	logger.Debug(format, a...)
}

func Info(format string, a ...interface{}) {
	logger.Info(format, a...)
}

func Warn(format string, a ...interface{}) {
	logger.Warn(format, a...)
}

func Error(format string, a ...interface{}) {
	logger.Error(format, a...)
}

func Fatal(format string, a ...interface{}) {
	logger.Fatal(format, a...)
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
	return fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), fmt.Sprintf(format, a...))
}

type Handler interface {
	output(output interface{}) error
}

type ConsoleHandler struct {
	Handler
}

func NewConsoleHandler() Handler {
	return ConsoleHandler{}
}

func (c ConsoleHandler) output(output interface{}) error {
	_, err := fmt.Println(output)
	return err
}
