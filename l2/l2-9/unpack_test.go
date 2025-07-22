package unpack

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		isErr    bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"", "", false},
		{"45", "", true},
		{"b0", "", false},
		{"b1", "b", false},
		{"\\32", "33", false},
		{"a10", "aaaaaaaaaa", false},
		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
		{"\\", "", true},
		{"a\\", "", true},
		{"x\\10", "x", false},
		{"x\\12", "x11", false},
		{"\\\\3", "\\\\\\", false}, // Экранирую символ экранирования
		{"\\", "", true},
		{"2a", "", true},
	}

	for _, test := range tests {
		result, err := Unpack(test.input)
		if test.isErr {
			if err == nil {
				t.Errorf("expected error for input %q", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("for input %q expected %q, got %q", test.input, test.expected, result)
			}
		}
	}
}
