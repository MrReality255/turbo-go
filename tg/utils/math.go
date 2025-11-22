package utils

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type Counter struct {
	mx    sync.Mutex
	value int64
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
