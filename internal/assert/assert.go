package assert

import (
	"errors"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want %v", actual, expected)
	}
}

func ErrorIs(t *testing.T, actual, expected error) {
	t.Helper()

	if !errors.Is(actual, expected) {
		t.Errorf("got: %v; want %v", actual, expected)
	}
}
