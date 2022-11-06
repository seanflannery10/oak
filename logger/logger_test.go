package logger

import (
	"bytes"
	"errors"
	"github.com/seanflannery10/ossa/assert"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

func TestLogger(t *testing.T) {
	b := new(bytes.Buffer)
	l := New()

	l.SetOutput(b)

	assert.Equal(t, l.GetLevel(), LevelInfo)

	l.Debug("debug")
	assert.NotContains(t, b.String(), "debug")

	l.SetLevel(LevelWarning)

	l.Info("info")
	assert.NotContains(t, b.String(), "info")

	l.SetLevel(LevelDebug)

	l.Debug("debug")
	assert.Contains(t, b.String(), "debug")

	l.Info("info")
	assert.Contains(t, b.String(), "info")

	l.Warning("warning")
	assert.Contains(t, b.String(), "warning")

	l.Error(errors.New("error"), map[string]string{"error": "error"})
	assert.Contains(t, b.String(), "error")
}

func TestLogger_Fatal(t *testing.T) {
	b := new(bytes.Buffer)
	l := New()

	l.SetOutput(b)

	if os.Getenv("TEST") == "1" {
		l.Fatal(errors.New("fatal"), map[string]string{"fatal": "fatal"})
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLogger")
	cmd.Env = append(os.Environ(), "TEST=1")

	err := cmd.Run()
	assert.Equal(t, err.Error(), "exit status 1")
}

func TestLogger_SetTimeFormat(t *testing.T) {
	b := new(bytes.Buffer)
	l := New()

	l.SetOutput(b)
	l.SetTimeFormat(DateTime)

	l.Debug("debug")

	m := regexp.MustCompile(`"time":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z"`)
	assert.Contains(t, b.String(), m.FindString(b.String()))

	l.SetTimeFormat(UnixTime)

	l.Info("info")

	m = regexp.MustCompile(`"time\\":\\"\d{6,}\\"`)
	assert.Contains(t, b.String(), m.FindString(b.String()))
}

func TestGlobalLogger(t *testing.T) {
	b := new(bytes.Buffer)

	SetOutput(b)

	assert.Equal(t, GetLevel(), LevelInfo)

	Debug("debug")
	assert.NotContains(t, b.String(), "debug")

	SetLevel(LevelWarning)

	Info("info")
	assert.NotContains(t, b.String(), "info")

	SetLevel(LevelDebug)

	Debug("debug")
	assert.Contains(t, b.String(), "debug")

	Info("info")
	assert.Contains(t, b.String(), "info")

	Warning("warning")
	assert.Contains(t, b.String(), "warning")

	Error(errors.New("error"), map[string]string{"error": "error"})
	assert.Contains(t, b.String(), "error")

	assert.Equal(t, GetLevel(), LevelDebug)
}

func TestGlobalLogger_Fatal(t *testing.T) {
	b := new(bytes.Buffer)

	SetOutput(b)

	if os.Getenv("TEST") == "1" {
		Fatal(errors.New("fatal"), map[string]string{"fatal": "fatal"})
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLogger")
	cmd.Env = append(os.Environ(), "TEST=1")

	err := cmd.Run()
	assert.Equal(t, err.Error(), "exit status 1")
}

func TestGlobalLogger_SetTimeFormat(t *testing.T) {
	b := new(bytes.Buffer)

	SetOutput(b)
	SetTimeFormat(DateTime)

	Debug("debug")

	m := regexp.MustCompile(`"time":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z"`)
	assert.Contains(t, b.String(), m.FindString(b.String()))

	SetTimeFormat(UnixTime)

	Info("info")

	m = regexp.MustCompile(`"time\\":\\"\d{6,}\\"`)
	assert.Contains(t, b.String(), m.FindString(b.String()))
}
