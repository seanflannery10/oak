package assert

import (
	"fmt"
	"testing"
)

func TestEqualInt(t *testing.T) {
	tests := []struct {
		a int
		e int
		r bool
	}{
		{200, 200, true},
		{201, 200, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.a), func(t *testing.T) {
			mockT := new(testing.T)
			res := Equal(mockT, tt.a, tt.e)
			if res != tt.r {
				t.Errorf("Equal(%#v, %#v) should return %#v", tt.e, tt.a, tt.r)
			}
		})
	}
}

func TestEqualString(t *testing.T) {
	tests := []struct {
		a string
		e string
		r bool
	}{
		{"123", "123", true},
		{"124", "123", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.a), func(t *testing.T) {
			mockT := new(testing.T)
			res := Equal(mockT, tt.a, tt.e)
			if res != tt.r {
				t.Errorf("Equal(%#v, %#v) should return %#v", tt.e, tt.a, tt.r)
			}
		})
	}
}
