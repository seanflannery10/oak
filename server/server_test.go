package server

import (
	"github.com/seanflannery10/oak/assert"
	"testing"
)

func TestNew(t *testing.T) {
	assert.SameType(t, New(), &Server{})
}
