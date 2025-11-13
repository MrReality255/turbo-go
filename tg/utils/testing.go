package utils

import (
	"fmt"
	"testing"
)

func TestAsString(t *testing.T, hint string, testCase int, required string, actual any) {
	TestString(t, hint, testCase, required, fmt.Sprintf("%v", actual))
}

func TestString(t *testing.T, hint string, testCase int, required string, actual string) {
	if required != actual {
		t.Errorf("%s #%v failed: \nexpected %s\nreceived %s", hint, testCase, required, actual)
	}
}

func TestStringF(t *testing.T, hint string, testCase int, want string, got string, args ...any) {
	TestString(t, hint, testCase, want, fmt.Sprintf(got, args...))
}
