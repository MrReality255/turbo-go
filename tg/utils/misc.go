package utils

import (
	"fmt"
	"strconv"
)

type CountMap map[*int]bool

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

func PermEx(
	fctResult func(caseNr int, vector []float64),
	fcts ...func(prev []float64) (from float64, till float64, step float64),
) {
	var caseNr int
	permStep(fctResult, fcts, nil, &caseNr)
}

func permStep(
	result func(caseNr int, vector []float64),
	fcts []func(prev []float64) (from float64, till float64, step float64),
	prefix []float64, caseNr *int,
) {
	if len(fcts) == 0 {
		*caseNr++
		result(*caseNr, prefix)
		return
	}

	from, till, step := fcts[0](prefix)
	for x := from; x <= till; x += step {
		permStep(result, fcts[1:], append(prefix, x), caseNr)
	}
}

func Clamp[T float64 | int](value T, min T, max T) T {
	switch {
	case value < min:
		return min
	case value > max:
		return max
	default:
		return value
	}
}

func LeastSquareSlope(values []float64) float64 {
	var sumX, sumY, sumXY, sumX2 float64

	n := len(values)
	for idx, value := range values {
		xFloat := float64(idx)
		sumX += xFloat
		sumY += value
		sumXY = xFloat * value
		sumX2 = xFloat * xFloat
	}

	numerator := float64(n)*sumXY - sumX*sumY
	denominator := float64(n)*sumX2 - sumX*sumX
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

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
