package utils

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type Counter struct {
	mx    sync.Mutex
	value int64
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

type LinearFunction[S Numeric, T Numeric] struct {
	k float64
	q float64
}

func NewLinearFunction[S Numeric, T Numeric](x1 float64, y1 float64, x2 float64, y2 float64) *LinearFunction[S, T] {
	k := (y2 - y1) / (x2 - x1)
	return &LinearFunction[S, T]{
		k: k,
		q: y1 - k*x1,
	}
}

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

func (f *LinearFunction[S, T]) Eval(x S) T {
	return T(f.k*float64(x) + f.q)
}

func (f *LinearFunction[S, T]) Reverse(y T) S {
	return S((float64(y) - f.q) / f.k)
}

func TruncStep[T constraints.Integer](src T, step T) T {
	return (src / step) * step
}

func NewCounter() *Counter {
	return &Counter{value: 0}
}

func (c *Counter) Next64() int64 {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.value++
	return c.value
}

func (c *Counter) Next() int {
	return int(c.Next64())
}

func (c *Counter) Inc() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.value++
}

func (c *Counter) Get() int64 {
	c.mx.Lock()
	defer c.mx.Unlock()
	return c.value
}

func Align[T constraints.Integer](src T, step T) T {
	return (src / step) * step
}
