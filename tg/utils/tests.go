package utils

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T, hint string, testCase int, want string, got string) {
	if got != want {
		t.Errorf("Test %s %d: got `%s`, want `%s`", hint, testCase, got, want)
	}
}

func TestStringF(t *testing.T, hint string, testCase int, want string, got string, args ...any) {
	TestString(t, hint, testCase, want, fmt.Sprintf(got, args...))
}
