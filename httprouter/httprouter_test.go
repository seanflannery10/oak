package httprouter

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/seanflannery10/ossa/assert"
)

func TestNew(t *testing.T) {
	assert.SameType(t, New(), httprouter.New())
}
