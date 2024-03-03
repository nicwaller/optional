package wrapper

import (
	"reflect"
)

type Wrapper[T any] struct {
	// The non-exported member prevents unsafe access to nil pointers
	rawPointer *T
}

func New[T any](ptr *T) Wrapper[T] {
	return Wrapper[T]{
		rawPointer: ptr,
	}
}

func (o Wrapper[T]) Nil() bool {
	return o.rawPointer == nil
}

// Unwrap safely calls the provided function with a not-nil pointer
func (o Wrapper[T]) Unwrap(f func(*T)) {
	if ptr := o.rawPointer; ptr != nil {
		f(ptr)
	}
}

// Or returns the unwrapped pointer, but only if it's not nil
// Otherwise, it returns a pointer to a valid default
func (o Wrapper[T]) Or(alternative T) *T {
	if ptr := o.rawPointer; ptr != nil {
		return ptr
	} else {
		return &alternative
	}
}

func (o Wrapper[T]) Equal(cmp T) bool {
	if o.rawPointer == nil {
		// This makes the compiler happy when comparing any
		return func(v any) bool {
			return v == nil
		}(cmp)
	}
	if o.rawPointer != nil {
		if v := reflect.ValueOf(*o.rawPointer); v.Comparable() {
			u := reflect.ValueOf(cmp)
			return v.Equal(u)
		}
	}
	return false
}
