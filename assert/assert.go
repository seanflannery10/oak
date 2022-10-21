package assert

import (
	"reflect"
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) bool {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
		return false
	}

	return true
}

func SameType[T comparable](t *testing.T, actual, expected T) bool {
	t.Helper()

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("got: %v; want: %v", actual, expected)
		return false
	}

	return true
}

func StringContains(t *testing.T, actual, expectedSubstring string) bool {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
		return false
	}

	return true
}

func NilError(t *testing.T, actual error) bool {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
		return false
	}
	return true
}