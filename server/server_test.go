package server

import (
	"github.com/seanflannery10/oak/assert"
	"testing"
)

func TestNew(t *testing.T) {
	srv := New("test", nil)

	assert.Equal(t, srv.Addr, "test")
	assert.SameType(t, srv, &Server{})
}
