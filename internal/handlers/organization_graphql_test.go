package handlers

import (
	"testing"
)

func TestNilCoalescing(t *testing.T) {
	defaultValue := "default"
	got := nilCoalescing(nil, defaultValue)

	if defaultValue != got {
		t.Errorf("Expected: %q Got: %q", defaultValue, got)
	}

	expected := "expected"
	got = nilCoalescing(&expected, defaultValue)

	if got != expected {
		t.Errorf("Expected: %q Got: %q", defaultValue, got)
	}
}
