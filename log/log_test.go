package log

import (
	"bytes"
	"errors"
	"github.com/seanflannery10/oak/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	b := new(bytes.Buffer)
	l := New()

	l.SetOutput(b)

	l.Debug("debug")
	l.Info("info")
	l.Warning("warning")
	l.Error(errors.New("error"), map[string]string{"error": "error"})

	assert.StringContains(t, b.String(), "debug")
	assert.StringContains(t, b.String(), "info")
	assert.StringContains(t, b.String(), "warning")
	assert.StringContains(t, b.String(), "error")
}
