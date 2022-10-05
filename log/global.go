package log

import (
	"io"
	"strconv"
	"time"
)

var (
	DateTime = time.Now().UTC().Format(time.RFC3339)
	UnixTime = strconv.FormatInt(time.Now().Unix(), 10)

	GlobalLogger = New()
)

func GetLevel() Level {
	return GlobalLogger.GetLevel()
}

func SetLevel(level Level) {
	GlobalLogger.SetLevel(level)
}

func UseUnixTime() {
	GlobalLogger.UseUnixTime()
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
