package utils

import "testing"

func TestIsNil(t *testing.T) {
	for idx, tc := range []struct {
		value  any
		result string
	}{
		{
			value:  int64(0),
			result: "true",
		},
		{
			value:  "",
			result: "true",
		},
		{
			value:  nil,
			result: "true",
		},
		{
			value:  int64(1),
			result: "false",
		},
	} {
		v := IsNil(tc.value)
		TestAsString(t, idx+1, "isNil", tc.result, v)
	}
}
