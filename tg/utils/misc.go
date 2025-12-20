package utils

import (
	"errors"
	"fmt"
	"strconv"
)

type CountMap map[*int]bool

type ItemWithErr[T any] struct {
	Data T
	Err  error
}

func ConstFct[T any](value T) func() T {
	return func() T {
		return value
	}
}

func IsErr(err error, options ...error) bool {
	for _, o := range options {
		if errors.Is(err, o) {
			return true
		}
	}
	return false
}

func ExecErr(f interface{}) error {
	switch f := f.(type) {
	case func() error:
		return f()
	case func():
		f()
		return nil
	default:
		return fmt.Errorf("invalid function: %T", f)
	}
}

func Must[T any](src T, err error) T {
	if err != nil {
		panic(err)
	}
	return src
}

func MustSucceed(err error) {
	if err != nil {
		panic(err)
	}
}

func StrToFloat64(str string) float64 {
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return v
}

func Count(m CountMap) {
	for ptr, v := range m {
		if v {
			*ptr = *ptr + 1
		}
	}
}

/*
func Perm(count int, from float64, till float64, step float64, fct func(caseNr int, v []float64)) {
	var caseNr int
	permStep(count, from, till, step, fct, nil, &caseNr)

}
*/

/*
func permStep(count int, from float64, till float64, step float64, fct func(caseNr int, v []float64), prefix []float64, caseNr *int) {
	if count == 0 {
		*caseNr++
		fct(*caseNr, prefix)
		return
	}
	for v := from; v < till; v += step {
		permStep(count-1, from, till, step, fct, append(prefix, v), caseNr)
	}
}
*/

func DataOrErr[T any](data T, err error) ItemWithErr[T] {
	return ItemWithErr[T]{Data: data, Err: err}
}
