package optional

import (
	"fmt"
	"reflect"
)

type Optional[T any] struct {
	// The non-exported member prevents unsafe access to nil pointers
	rawPointer *T
}

//goland:noinspection GoUnusedExportedFunction
func FromPointer[T any](ptr *T) Optional[T] {
	return Optional[T]{
		rawPointer: ptr,
	}
}

func FromValue[T any](val T) Optional[T] {
	return FromPointer(&val)
}

func InSlice[T any](s []T, index int) Optional[T] {
	if index < 0 || index >= len(s) {
		return FromPointer[T](nil)
	}
	return FromValue(s[index])
}

func InMap[K comparable, V any](m map[K]V, k K) Optional[V] {
	if v, found := m[k]; found {
		return FromValue(v)
	} else {
		return FromPointer[V](nil)
	}
}

func (o *Optional[T]) Nil() bool {
	return o.rawPointer == nil
}

// Unwrap safely calls the provided function with a not-nil pointer
func (o Optional[T]) Unwrap(f func(safePtr *T)) {
	if ptr := o.rawPointer; ptr != nil {
		f(ptr)
	}
}

// Or returns the unwrapped pointer, but only if it's not nil
// Otherwise, it returns a pointer to a valid default
func (o *Optional[T]) Or(alternative T) *T {
	if ptr := o.rawPointer; ptr != nil {
		return ptr
	} else {
		return &alternative
	}
}

func (o *Optional[T]) Equal(cmp T) bool {
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

func (o *Optional[T]) SetPointer(ptr *T) {
	o.rawPointer = ptr
}

func (o *Optional[T]) SetValue(ptr T) {
	o.rawPointer = &ptr
}

// let the compiler verify interface compatibility
var _ fmt.Stringer = Optional[any]{}
var _ fmt.GoStringer = Optional[any]{}

//goland:noinspection GoMixedReceiverTypes
func (o Optional[T]) String() string {
	return fmt.Sprintf("<Optional:%v>", o.GoString())
}

//goland:noinspection GoMixedReceiverTypes
func (o Optional[T]) GoString() string {
	if o.rawPointer == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", *o.rawPointer)
}
