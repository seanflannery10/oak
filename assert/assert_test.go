package assert

import (
	"fmt"
	"testing"
)

func TestEqual(t *testing.T) {
	testsNums := []struct {
		a float32
		e float32
		r bool
	}{
		{200, 200, true},
		{201, 200, false},
		{20.0, 20.0, true},
		{20.1, 20.0, false},
	}

	for _, tt := range testsNums {
		t.Run(fmt.Sprintf("%#v", tt.a), func(t *testing.T) {
			res := Equal(new(testing.T), tt.a, tt.e)
			if res != tt.r {
				t.Errorf("Equal(%#v, %#v) should return %#v", tt.e, tt.a, tt.r)
			}
		})
	}

	testsStrings := []struct {
		a string
		e string
		r bool
	}{
		{"123", "123", true},
		{"124", "123", false},
	}

	for _, tt := range testsStrings {
		t.Run(fmt.Sprintf("%#v", tt.a), func(t *testing.T) {
			res := Equal(new(testing.T), tt.a, tt.e)
			if res != tt.r {
				t.Errorf("Equal(%#v, %#v) should return %#v", tt.e, tt.a, tt.r)
			}
		})
	}
}
