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
	// DEBUG level of DEBUG
	DEBUG = 1
	// INFO level of DEBUG
	INFO = 2
	// WARN level of DEBUG
	WARN = 3
	// ERROR level of DEBUG
	ERROR = 4

	loggerRootName = ""
	loggerEnv      = os.Getenv("LOKI_ENV")
	logger         = New(loggerRootName)
	globs          []glob.Glob
)

func init() {
	if loggerEnv != "" {
		patterns := strings.Split(loggerEnv, ",")
		for _, pattern := range patterns {
			globs = append(globs, glob.MustCompile(pattern))
		}
	}
}

// Logger is a instance to handle logs
type Logger struct {
	name      string
	level     int
	formatter Formatter
	handler   Handler
}

// New create a Logger instance with with its name
func New(name string) Logger {
	return Logger{
		name:      name,
		level:     INFO,
		formatter: NewStandardFormatter(),
		handler:   NewConsoleHandler(),
	}
}

// Check if logger's name match the LOKI_ENV setting
func (l Logger) Check() bool {
	if l.name == loggerRootName {
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

// SetLevel set the level of logger
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// SetFormatter set the formatter of logger
func (l *Logger) SetFormatter(formatter Formatter) {
	l.formatter = formatter
}

// Debug output level DEBUG log
func (l Logger) Debug(a ...interface{}) {
	if DEBUG >= l.level {
		l.handler.output(l.formatter.format(a...))
	}
}

// Info output level INFO log
func (l Logger) Info(a ...interface{}) {
	if INFO >= l.level {
		l.handler.output(aurora.Blue(l.formatter.format(a...)))
	}
}

// Warn output level WARN log
func (l Logger) Warn(a ...interface{}) {
	if WARN >= l.level {
		l.handler.output(aurora.Green(l.formatter.format(a...)))
	}
}

// Error output level ERROR log
func (l Logger) Error(a ...interface{}) {
	if ERROR >= l.level {
		l.handler.output(aurora.Red(l.formatter.format(a...)))
	}
}

// Fatal output level ERROR log and exit with code 1
func (l Logger) Fatal(a ...interface{}) {
	Error(a...)
	os.Exit(1)
}

// SetLevel set the level of logger
func SetLevel(level int) {
	logger.SetLevel(level)
}

// SetFormatter set the formatter of logger
func SetFormatter(formatter Formatter) {
	logger.SetFormatter(formatter)
}

// Debug output level DEBUG log
func Debug(a ...interface{}) {
	logger.Debug(a...)
}

// Info output level INFO log
func Info(a ...interface{}) {
	logger.Info(a...)
}

// Warn output level WARN log
func Warn(a ...interface{}) {
	logger.Warn(a...)
}

// Error output level ERROR log
func Error(a ...interface{}) {
	logger.Error(a...)
}

// Fatal output level ERROR log and exit with code 1
func Fatal(a ...interface{}) {
	logger.Fatal(a...)
}

// Formatter format the message
type Formatter interface {
	format(a ...interface{}) string
}

// StandardFormatter default formatter
type StandardFormatter struct {
	Formatter
}

// NewStandardFormatter ...
func NewStandardFormatter() Formatter {
	return StandardFormatter{}
}

func (f StandardFormatter) format(a ...interface{}) string {
	if len(a) == 0 {
		return ""
	}

	ts := time.Now().Format(time.RFC3339)
	format, ok := a[0].(string)
	if !ok {
		return fmt.Sprintf("%s Logger format error with args %s", ts, a)
	}
	return fmt.Sprintf("%s %s", ts, fmt.Sprintf(format, a[1:]...))
}

// Handler handle the output process
type Handler interface {
	output(output interface{}) error
}

// ConsoleHandler output logs to console
type ConsoleHandler struct {
	Handler
}

// NewConsoleHandler ...
func NewConsoleHandler() Handler {
	return ConsoleHandler{}
}

func (c ConsoleHandler) output(output interface{}) error {
	_, err := fmt.Println(output)
	return err
}
