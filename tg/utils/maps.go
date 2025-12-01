package utils

func MapClone[T comparable, S any](m map[T]S) map[T]S {
	r := make(map[T]S)
	for k, v := range m {
		r[k] = v
	}
	return r
}

func MapKeys[C comparable, T any](
	src map[C]T, lessFct func(item1 C, item2 C) bool,
) []C {
	return MapKeysIf(src, nil, lessFct)
}

func MapKeysIf[C comparable, T any](
	src map[C]T, condFct func(item C, value T) bool, lessFct func(item1 C, item2 C) bool,
) []C {
	result := make([]C, 0, len(src))
	for key, value := range src {
		if condFct == nil || condFct(key, value) {
			result = append(result, key)
		}
	}
	if lessFct != nil {
		Sort(result, lessFct)
	}
	return result

}

func MapValues[C comparable, T any](src map[C]T, lessFct func(item1 T, item2 T) bool) []T {
	return MapValuesIf(src, nil, lessFct)
}

func MapValuesIf[C comparable, T any](
	src map[C]T, selectFct func(key C, value T) bool, lessFct func(item1 T, item2 T) bool,
) []T {
	result := make([]T, 0, len(src))
	for key, value := range src {
		if selectFct == nil || selectFct(key, value) {
			result = append(result, value)
		}
	}
	if lessFct != nil {
		Sort(result, lessFct)
	}
	return result
}

func MapMap[K comparable, T any, K2 comparable, T2 any](
	src map[K]T,
	keyFct func(item K) K2,
	valueFct func(key K, item T) T2,
) map[K2]T2 {
	result := make(map[K2]T2)
	for k, v := range src {
		result[keyFct(k)] = valueFct(k, v)
	}
	return result
}
func MapToArray[T comparable, S any, C any](
	src map[T]S, fct func(key T, value S) C, lessFct func(item1 C, item2 C) bool) []C {
	result := make([]C, 0, len(src))
	for k, v := range src {
		result = append(result, fct(k, v))
	}
	if lessFct != nil {
		SortArray(result, lessFct)
	}
	return result
}

func MapMapValues[K comparable, T any, T2 any](src map[K]T, valueFct func(key K, item T) T2) map[K]T2 {
	return MapMap(
		src,
		func(item K) K {
			return item
		},
		valueFct,
	)
}
