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

func New() (l *Logger) {
	l = &Logger{
		minLevel:   uint32(LevelInfo),
		timeFormat: DateTime,
	}

	l.SetOutput(os.Stdout)
	return
}

func (l *Logger) GetLevel() Level {
	return Level(atomic.LoadUint32(&l.minLevel))
}

func (l *Logger) SetLevel(level Level) {
	atomic.StoreUint32(&l.minLevel, uint32(level))
}

func (l *Logger) UseUnixTime() {
	l.timeFormat = UnixTime
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

	line := l.jsonLine(level, message, properties)

	l.mu.Lock()
	defer l.mu.Unlock()

	fmt.Fprintln(l.output, line)
}

func (l *Logger) jsonLine(level Level, message string, properties map[string]string) string {
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
		return fmt.Sprintf("%s: unable to marshal log message: %s", LevelError.String(), err.Error())
	}

	return string(line)
}
