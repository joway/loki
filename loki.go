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

// SetLevel set the level of logger
func SetLevel(level int) {
	logger.level = level
}

// SetFormatter set the formatter of logger
func SetFormatter(formatter Formatter) {
	logger.formatter = formatter
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

// Debug output level DEBUG log
func (l Logger) Debug(format string, a ...interface{}) {
	if DEBUG >= l.level {
		l.handler.output(l.formatter.format(format, a...))
	}
}

// Info output level INFO log
func (l Logger) Info(format string, a ...interface{}) {
	if INFO >= l.level {
		l.handler.output(aurora.Blue(l.formatter.format(format, a...)))
	}
}

// Warn output level WARN log
func (l Logger) Warn(format string, a ...interface{}) {
	if WARN >= l.level {
		l.handler.output(aurora.Green(l.formatter.format(format, a...)))
	}
}

// Error output level ERROR log
func (l Logger) Error(format string, a ...interface{}) {
	if ERROR >= l.level {
		l.handler.output(aurora.Red(l.formatter.format(format, a...)))
	}
}

// Fatal output level ERROR log and exit with code 1
func (l Logger) Fatal(format string, a ...interface{}) {
	Error(format, a...)
	os.Exit(1)
}

// Debug output level DEBUG log
func Debug(format string, a ...interface{}) {
	logger.Debug(format, a...)
}

// Info output level INFO log
func Info(format string, a ...interface{}) {
	logger.Info(format, a...)
}

// Warn output level WARN log
func Warn(format string, a ...interface{}) {
	logger.Warn(format, a...)
}

// Error output level ERROR log
func Error(format string, a ...interface{}) {
	logger.Error(format, a...)
}

// Fatal output level ERROR log and exit with code 1
func Fatal(format string, a ...interface{}) {
	logger.Fatal(format, a...)
}

// Formatter format the message
type Formatter interface {
	format(format string, a ...interface{}) string
}

// StandardFormatter default formatter
type StandardFormatter struct {
	Formatter
}

// NewStandardFormatter ...
func NewStandardFormatter() Formatter {
	return StandardFormatter{}
}

func (f StandardFormatter) format(format string, a ...interface{}) string {
	return fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), fmt.Sprintf(format, a...))
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
