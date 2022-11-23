package validator

import (
	"github.com/seanflannery10/ossa/assert"
	"testing"
)

func TestValidator(t *testing.T) {
	var v Validator

	assert.Equal(t, v.HasErrors(), false)

	v.AddError("test", "test error")
	v.AddFieldError("test", "test field error")

	assert.Equal(t, v.Errors["test"], "test error")
	assert.Equal(t, v.FieldErrors["test"], "test field error")

	v.Check(true, "test2", "test2")
	v.CheckField(true, "test2", "test field error2")

	assert.Equal(t, len(v.Errors), 1)
	assert.Equal(t, len(v.FieldErrors), 1)

	v.Check(false, "test3", "test3")
	v.CheckField(false, "test3", "test field error3")

	assert.Equal(t, v.Errors["test3"], "test3")
	assert.Equal(t, v.FieldErrors["test3"], "test field error3")

	assert.Equal(t, v.HasErrors(), true)
}
