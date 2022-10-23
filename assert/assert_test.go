package assert

import (
	"errors"
	"fmt"
	"testing"
)

func TestEqual(t *testing.T) {
	testsNums := []struct {
		a float32
		e float32
		r bool
	}{
		{
			0,
			0,
			true,
		},
		{
			200,
			200,
			true,
		},
		{
			201,
			200,
			false,
		},
		{
			200.0,
			200.0,
			true,
		},
		{
			200.1,
			200.0,
			false,
		},
		{
			-200,
			-200,
			true,
		},
		{
			-201,
			-200,
			false,
		},
		{
			-200.0,
			-200.0,
			true,
		},
		{
			-200.1,
			-200.0,
			false,
		},
	}

	for _, tt := range testsNums {
		t.Run(fmt.Sprintf("%#v", tt.a), func(t *testing.T) {
			res := Equal(new(testing.T), tt.a, tt.e)
			if res != tt.r {
				t.Errorf("Equal(%#v, %#v) should return %#v", tt.a, tt.e, tt.r)
			}
		})
	}

	testsStrings := []struct {
		a string
		e string
		r bool
	}{
		{
			"123",
			"123",
			true,
		},
		{
			"124",
			"123",
			false,
		},
		{
			"",
			"",
			true,
		},
	}

	for _, tt := range testsStrings {
		t.Run(fmt.Sprintf("%#v", tt.a), func(t *testing.T) {
			res := Equal(new(testing.T), tt.a, tt.e)
			if res != tt.r {
				t.Errorf("Equal(%#v, %#v) should return %#v", tt.a, tt.e, tt.r)
			}
		})
	}
}

func TestSameType(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		res := SameType(new(testing.T), "foo", "bar")
		if res != true {
			t.Errorf("SameType(%#v, %#v) should return %#v", "foo", "bar", "true")
		}
	})

	t.Run("Int", func(t *testing.T) {
		res := SameType(new(testing.T), 1, 2)
		if res != true {
			t.Errorf("SameType(%#v, %#v) should return %#v", "foo", "bar", "true")
		}
	})

	t.Run("Float", func(t *testing.T) {
		res := SameType(new(testing.T), 1.2, 2.3)
		if res != true {
			t.Errorf("SameType(%#v, %#v) should return %#v", "foo", "bar", "true")
		}
	})

	t.Run("Bool", func(t *testing.T) {
		res := SameType(new(testing.T), true, false)
		if res != true {
			t.Errorf("SameType(%#v, %#v) should return %#v", "foo", "bar", "true")
		}
	})
}

func TestContains(t *testing.T) {
	tests := []struct {
		a string
		e string
		r bool
	}{
		{
			"this is a test",
			"test",
			true,
		},
		{
			"this is a",
			"test",
			false,
		},
		{
			"",
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.e), func(t *testing.T) {
			res := Contains(new(testing.T), tt.a, tt.e)
			if res != tt.r {
				t.Errorf("Contains(%#v, %#v) should return %#v", tt.a, tt.e, tt.r)
			}
		})
	}
}

func TestNilError(t *testing.T) {
	tests := []struct {
		a error
		r bool
	}{
		{
			nil,
			true,
		},
		{
			errors.New("test"),
			false,
		},
		{
			errors.New(""),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.a), func(t *testing.T) {
			res := NilError(new(testing.T), tt.a)
			if res != tt.r {
				t.Errorf("NilError(%#v) should return %#v", tt.a, tt.r)
			}
		})
	}
}
