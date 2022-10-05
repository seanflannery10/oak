package log

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
