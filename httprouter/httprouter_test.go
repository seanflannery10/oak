package httprouter

import (
	"github.com/julienschmidt/httprouter"
	"github.com/seanflannery10/ossa/assert"
	"testing"
)

func TestNew(t *testing.T) {
	assert.SameType(t, New(), httprouter.New())
}
