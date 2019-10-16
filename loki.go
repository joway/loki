package loki

import (
	"bufio"
	"fmt"
	"github.com/gobwas/glob"
	"github.com/logrusorgru/aurora"
	"io"
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
	logger         = New(loggerRootName)
)

// Logger is a instance to handle logs
type Logger struct {
	name      string
	level     int
	env       string
	patterns  []glob.Glob
	formatter Formatter
	handler   Handler

	timeFormat string
}

// New create a Logger instance with with its name
func New(name string) Logger {
	//get patterns
	env, isExisted := os.LookupEnv("LOKI_ENV")
	if !isExisted {
		//enable all loggers by default
		env = "*"
	}

	l := Logger{
		name:       name,
		level:      INFO,
		env:        env,
		formatter:  NewStandardFormatter(),
		handler:    NewConsoleHandler(),
		timeFormat: time.RFC3339,
	}
	l.SetLogEnv(env)

	return l
}

// Check if logger's name match the LOKI_ENV setting
func (l Logger) Check() bool {
	if l.name == loggerRootName {
		return true
	}
	for _, g := range l.patterns {
		matched := g.Match(l.name)
		if matched {
			return true
		}
	}
	return false
}

// SetLogEnv set the LOKI_ENV of logger
func (l *Logger) SetLogEnv(env string) {
	pats := strings.Split(env, ",")
	var patterns []glob.Glob
	for _, p := range pats {
		patterns = append(patterns, glob.MustCompile(p))
	}
	l.patterns = patterns
}

// SetLevel set the level of logger
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// SetHandler set the handler of logger
func (l *Logger) SetHandler(handler Handler) {
	l.handler = handler
}

// SetFormatter set the formatter of logger
func (l *Logger) SetFormatter(formatter Formatter) {
	l.formatter = formatter
}

// SetTimeFormatter set the time format string of logger
func (l *Logger) SetTimeFormatter(format string) {
	l.timeFormat = format
}

// Compile return final compiled log string
func (l Logger) Compile(a ...interface{}) string {
	msg := l.formatter.format(a...)
	prefix := ""
	if l.timeFormat != "" {
		ts := time.Now().Format(l.timeFormat)
		prefix += fmt.Sprintf("%s ", ts)
	}
	if l.name != loggerRootName {
		prefix += fmt.Sprintf("%s", l.name)
	}
	if prefix == "" {
		return msg
	} else {
		return fmt.Sprintf("%s %s", prefix, msg)
	}
}

// Debug output level DEBUG log
func (l Logger) Debug(a ...interface{}) {
	if l.Check() && DEBUG >= l.level {
		_ = l.handler.debug(l.Compile(a...))
	}
}

// Info output level INFO log
func (l Logger) Info(a ...interface{}) {
	if l.Check() && INFO >= l.level {
		_ = l.handler.info(l.Compile(a...))
	}
}

// Warn output level WARN log
func (l Logger) Warn(a ...interface{}) {
	if l.Check() && WARN >= l.level {
		_ = l.handler.warn(l.Compile(a...))
	}
}

// Error output level ERROR log
func (l Logger) Error(a ...interface{}) {
	if l.Check() && ERROR >= l.level {
		_ = l.handler.error(l.Compile(a...))
	}
}

// Fatal output level ERROR log and exit with code 1
func (l Logger) Fatal(a ...interface{}) {
	if l.Check() && ERROR >= l.level {
		_ = l.handler.error(l.Compile(a...))
		os.Exit(1)
	}
}

// SetLevel set the level of logger
func SetLevel(level int) {
	logger.SetLevel(level)
}

// SetHandler set the handler of logger
func SetHandler(handler Handler) {
	logger.SetHandler(handler)
}

// SetTimeFormatter set the time format string of logger
func SetTimeFormatter(format string) {
	logger.timeFormat = format
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
	format, ok := a[0].(string)
	if !ok {
		return fmt.Sprintf("Logger format error with args %s", a)
	}
	return fmt.Sprintf(format, a[1:]...)
}

// Handler handle the output process
type Handler interface {
	debug(output string) error
	info(output string) error
	warn(output string) error
	error(output string) error
}

// ConsoleHandler output logs to console
type ConsoleHandler struct {
	Handler
}

// NewConsoleHandler return ConsoleHandler instance
func NewConsoleHandler() Handler {
	return ConsoleHandler{}
}

func (handler ConsoleHandler) debug(output string) error {
	_, err := fmt.Println(output)
	return err
}
func (handler ConsoleHandler) info(output string) error {
	_, err := fmt.Println(aurora.Blue(output))
	return err
}
func (handler ConsoleHandler) warn(output string) error {
	_, err := fmt.Println(aurora.Green(output))
	return err
}
func (handler ConsoleHandler) error(output string) error {
	_, err := fmt.Println(aurora.Red(output))
	return err
}

// FileHandler output logs to console
type FileHandler struct {
	Handler
	writer *bufio.Writer
}

// NewFileHandler return FileHandler instance
func NewFileHandler(fp *os.File, flushIntervalMs int) Handler {
	h := FileHandler{
		writer: bufio.NewWriter(fp),
	}
	timer := time.NewTimer(time.Duration(flushIntervalMs) * time.Millisecond)
	go func() {
		<-timer.C
		_ = h.writer.Flush()
	}()
	return h
}

func (handler FileHandler) append(msg string) error {
	data := []byte(msg + "\n")
	n, err := handler.writer.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	return err
}
func (handler FileHandler) debug(output string) error {
	return handler.append(output)
}
func (handler FileHandler) info(output string) error {
	return handler.append(output)
}
func (handler FileHandler) warn(output string) error {
	return handler.append(output)
}
func (handler FileHandler) error(output string) error {
	return handler.append(output)
}
