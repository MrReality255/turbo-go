package utils

import (
	"fmt"
	"testing"
)

func TestAsString(t *testing.T, testCase int, hint string, required string, actual any) {
	TestString(t, testCase, hint, required, fmt.Sprintf("%v", actual))
}

func TestString(t *testing.T, testCase int, hint string, required string, actual string) {
	if required != actual {
		t.Errorf("%s #%v failed: \nexpected %s\nreceived %s", hint, testCase, required, actual)
	}
}

func TestStringF(t *testing.T, testCase int, hint string, want string, got string, args ...any) {
	TestString(t, testCase, hint, want, fmt.Sprintf(got, args...))
}
