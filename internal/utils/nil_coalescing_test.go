package utils

import "testing"

func TestNilCoalescing(t *testing.T) {

	t.Run("Test string coalescing", func(t *testing.T) {
		defaultValue := "default"
		got := CoalesceStr(nil, defaultValue)

		if defaultValue != got {
			t.Errorf("Expected: %q Got: %q", defaultValue, got)
		}

		expected := "expected"
		got = CoalesceStr(&expected, defaultValue)

		if got != expected {
			t.Errorf("Expected: %q Got: %q", defaultValue, got)
		}
	})

	t.Run("Test int coalescing", func(t *testing.T) {
		defaultValue := 100
		got := CoalesceInt(nil, defaultValue)

		if defaultValue != got {
			t.Errorf("Expected: %q Got: %q", defaultValue, got)
		}

		expected := 100
		got = CoalesceInt(&expected, defaultValue)

		if got != expected {
			t.Errorf("Expected: %q Got: %q", defaultValue, got)
		}
	})
}
