package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelDisabled
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return ""
	}
}

type Logger struct {
	output     io.Writer
	minLevel   uint32
	timeFormat string
	mu         sync.Mutex
}

var (
	DateTime = time.Now().UTC().Format(time.RFC3339)
	UnixTime = strconv.FormatInt(time.Now().Unix(), 10)
)

func New() *Logger {
	return &Logger{
		output:     os.Stdout,
		minLevel:   uint32(LevelInfo),
		timeFormat: DateTime,
	}
}

func (l *Logger) GetLevel() Level {
	return Level(atomic.LoadUint32(&l.minLevel))
}

func (l *Logger) SetLevel(level Level) {
	atomic.StoreUint32(&l.minLevel, uint32(level))
}

func (l *Logger) SetTimeFormat(timeFormat string) {
	l.timeFormat = timeFormat
}

func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

func (l *Logger) Debug(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.print(LevelDebug, message, nil)
}

func (l *Logger) Info(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.print(LevelInfo, message, nil)
}

func (l *Logger) Warning(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.print(LevelWarning, message, nil)
}

func (l *Logger) Error(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) Fatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]string) {
	if level < l.GetLevel() {
		return
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       l.timeFormat,
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(fmt.Sprintf("%s: unable to marshal log message: %s", LevelError.String(), err.Error()))
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	fmt.Fprintln(l.output, string(line))
}

var GlobalLogger = New()

func GetLevel() Level {
	return GlobalLogger.GetLevel()
}

func SetLevel(level Level) {
	GlobalLogger.SetLevel(level)
}

func SetTimeFormat(timeFormat string) {
	GlobalLogger.SetTimeFormat(timeFormat)
}

func SetOutput(w io.Writer) {
	GlobalLogger.SetOutput(w)
}

func Debug(format string, v ...any) {
	GlobalLogger.Debug(format, v...)
}

func Info(format string, v ...any) {
	GlobalLogger.Info(format, v...)
}

func Warning(format string, v ...any) {
	GlobalLogger.Warning(format, v...)
}

func Error(err error, properties map[string]string) {
	GlobalLogger.Error(err, properties)
}

func Fatal(err error, properties map[string]string) {
	GlobalLogger.Fatal(err, properties)
}
