package utils

import (
	"reflect"
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Integer | constraints.Float
}

type RingBuffer[C any] struct {
	mx    sync.Mutex
	items []C
	ptr   int
}

type customSort struct {
	lessFct func(a int, b int) bool
	swapFct func(a int, b int)
	length  int
}

func Coalesce[T any](values ...T) T {
	var src T
	for _, item := range values {
		if !IsNil(item) {
			return item
		}
	}
	return src
}

func IfPass[C any](cond bool, fct func() (*C, error)) (*C, error) {
	if !cond {
		return nil, nil
	}
	return fct()
}

func IfThenAny(cond bool, valTrue any, valFalse any) interface{} {
	if cond {
		return valTrue
	}
	return valFalse
}

func IfThenFct[C any](cond bool, fctTrue func() C, fctFalse func() C) C {
	if cond {
		return fctTrue()
	}
	return fctFalse()
}

func IfThen[C any](cond bool, valTrue C, valFalse C) C {
	if cond {
		return valTrue
	}
	return valFalse
}

// IsNil handles typed nils (e.g., interface containing a nil pointer).
func IsNil[T any](val T) bool {
	// Shortcut for actual nil (untyped nil or nil interface)
	if any(val) == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func Sort[T any](src []T, lessFct func(item1 T, item2 T) bool) {
	cs := &customSort{
		lessFct: func(i int, j int) bool {
			return lessFct(src[i], src[j])
		},
		swapFct: func(i int, j int) {
			src[i], src[j] = src[j], src[i]
		},
		length: len(src),
	}
	sort.Sort(cs)
}

func NewRingBuffer[C any](length int) *RingBuffer[C] {
	return &RingBuffer[C]{
		items: make([]C, 0, length),
		ptr:   0,
	}
}

func (r *RingBuffer[C]) Add(items ...C) {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, item := range items {
		if len(r.items) < cap(r.items) {
			r.items = append(r.items, item)
			r.ptr++
			continue
		}

		items[r.ptr] = item
		r.ptr++
		r.ptr = r.ptr % cap(r.items)
	}
}

func (r *RingBuffer[C]) Content() []C {
	r.mx.Lock()
	defer r.mx.Unlock()

	// G H I J E F
	// 0 1 2 3 4 5
	var (
		prefix []C
		suffix []C
	)
	prefix = r.items[r.ptr:len(r.items)]
	suffix = r.items[:r.ptr]
	return append(prefix, suffix...)
}

func (c customSort) Len() int {
	return c.length
}

func (c customSort) Less(i, j int) bool {
	return c.lessFct(i, j)
}

func (c customSort) Swap(i, j int) {
	c.swapFct(i, j)
}
