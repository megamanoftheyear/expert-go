package cmp

import "reflect"

func Filled[T comparable](value *T) bool {
	if value == nil {
		return false
	}

	return *value == *(new(T))
}

func Nullable(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Ptr
}
