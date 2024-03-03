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
func NewOptional[T any](ptr *T) Optional[T] {
	return Optional[T]{
		rawPointer: ptr,
	}
}

func NewOptionalVal[T any](val T) Optional[T] {
	return NewOptional(&val)
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

func (o *Optional[T]) Set(ptr *T) {
	o.rawPointer = ptr
}

// SetValue is a convenience function
// because Go doesn't allow taking address of a literal scalar type
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
