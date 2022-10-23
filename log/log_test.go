package log

import (
	"bytes"
	"errors"
	"github.com/seanflannery10/oak/assert"
	"os"
	"os/exec"
	"testing"
)

func TestLogger(t *testing.T) {
	b := new(bytes.Buffer)
	l := New()

	l.SetOutput(b)
	l.SetLevel(LevelWarning)

	l.Debug("debug")
	assert.NotContains(t, b.String(), "debug")

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
