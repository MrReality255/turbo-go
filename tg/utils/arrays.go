package utils

import "sort"

func ArrayCast[S Numeric, T Numeric](src []S) []T {
	result := make([]T, 0, len(src))
	for _, item := range src {
		result = append(result, T(item))
	}
	return result
}

func ArrayClone[T any](src []T) []T {
	return ArrayFilter(src, nil)
}

func ArrayFilter[T any](src []T, checkFct func(item T) bool) []T {
	result := make([]T, 0, len(src))
	for _, item := range src {
		if checkFct == nil || checkFct(item) {
			result = append(result, item)
		}
	}
	return result
}

func ArrayGroupBy[T any, C comparable](items []T, keyFct func(item T) C) map[C][]T {
	result := make(map[C][]T)
	for _, item := range items {
		key := keyFct(item)
		result[key] = append(result[key], item)
	}
	return result
}

func ArrayHasAny[T any](arr []T, fct func(item T) bool) bool {
	for _, item := range arr {
		if fct == nil || fct(item) {
			return true
		}
	}
	return false
}

func ArrayMap[T any, V any](src []T, fct func(item T) V) []V {
	result, err := ArrayMapErr(src, func(item T) (V, error) {
		return fct(item), nil
	})
	// that can never happen
	if err != nil {
		panic(err)
	}
	return result
}

func ArrayMapEx[T any, V any](src []T, fct func(item T) (V, bool)) []V {
	result, _ := ArrayMapExErr(src, func(item T) (V, bool, error) {
		v, ok := fct(item)
		return v, ok, nil
	})
	return result
}

func ArrayMapExErr[T any, V any](src []T, fct func(item T) (V, bool, error)) ([]V, error) {
	result := make([]V, 0, len(src))

	for _, item := range src {
		item, ok, err := fct(item)
		switch {
		case err != nil:
			return nil, err
		case ok:
			result = append(result, item)
		}
	}
	return result, nil
}

func ArrayMapErr[T any, V any](src []T, fct func(item T) (V, error)) ([]V, error) {
	return ArrayMapExErr(src, func(item T) (V, bool, error) {
		v, err := fct(item)
		return v, true, err
	})
}

func ArraySelect[T any](src []*T, checkFct func(prev *T, next *T) bool) *T {
	var selected *T
	for _, item := range src {
		if selected == nil || checkFct(selected, item) {
			selected = item
		}
	}
	return selected
}

func ArraySort[T any](src []T, lessFct func(item1 T, item2 T) bool) []T {
	result := ArrayClone(src)
	SortArray(result, lessFct)
	return result
}

func ArrayToMap[T comparable, S comparable, V any](ar []T, keyFct func(item T) S, valueFct func(item T) V) map[S]V {
	return ArrayToMapEx(ar, func(item T) (S, bool) {
		return keyFct(item), true
	}, valueFct)
}

func ArrayToMapEx[T comparable, S comparable, V any](
	ar []T, keyFct func(item T) (S, bool), valueFct func(item T) V,
) map[S]V {
	result := make(map[S]V)
	for _, item := range ar {
		key, ok := keyFct(item)
		if ok {
			result[key] = valueFct(item)
		}
	}
	return result
}

func ArrayChoose[T any](src []*T, checkFct func(prev *T, next *T) bool) *T {
	var selected *T
	for _, item := range src {
		if selected == nil || checkFct(selected, item) {
			selected = item
		}
	}
	return selected
}

func SortArray[T any](src []T, lessFct func(item1 T, item2 T) bool) {
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
