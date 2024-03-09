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
func OptionalPointer[T any](ptr *T) Optional[T] {
	return Optional[T]{
		rawPointer: ptr,
	}
}

func OptionalValue[T any](val T) Optional[T] {
	return OptionalPointer(&val)
}

func OptionalFromSlice[T any](s []T, index int) Optional[T] {
	if index < 0 || index >= len(s) {
		return OptionalPointer[T](nil)
	}
	return OptionalValue(s[index])
}

func OptionalFromMap[K comparable, V any](m map[K]V, k K) Optional[V] {
	if v, found := m[k]; found {
		return OptionalValue(v)
	} else {
		return OptionalPointer[V](nil)
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

func (o Optional[T]) IfElse(iff func(safePtr *T), elsef func()) {
	if ptr := o.rawPointer; ptr != nil {
		if iff != nil {
			iff(ptr)
		}
	} else {
		if elsef != nil {
			elsef()
		}
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

func Delve[T any](maybeSlice any, indices ...int) Optional[T] {
	if maybeSlice == nil || indices == nil {
		return OptionalPointer[T](nil)
	}
	if len(indices) == 0 {
		return OptionalPointer[T](nil)
	}
	index := indices[0]
	if index < 0 {
		return OptionalPointer[T](nil)
	}
	rSliceValue := reflect.ValueOf(maybeSlice)
	//if rSliceValue.Kind() != reflect.Slice {
	//	return OptionalPointer[T](nil)
	//}
	if 0 == rSliceValue.Len() || index >= rSliceValue.Len() {
		return OptionalPointer[T](nil)
	}

	innerValue := rSliceValue.Index(index)
	if len(indices) == 1 {
		if finalValue, ok := innerValue.Interface().(T); ok {
			return OptionalValue(finalValue)
		} else {
			return OptionalPointer[T](nil)
		}
	} else {
		subIndices := indices[1:]
		if subSlice, ok := innerValue.Interface().([]T); ok {
			return Delve[T](subSlice, subIndices...)
		} else if subSlice, ok := innerValue.Interface().([][]T); ok {
			return Delve[T](subSlice, subIndices...)
		} else if subSlice, ok := innerValue.Interface().([][][]T); ok {
			return Delve[T](subSlice, subIndices...)
		} else if subSlice, ok := innerValue.Interface().([][][][]T); ok {
			return Delve[T](subSlice, subIndices...)
		} else if subSlice, ok := innerValue.Interface().([][][][][]T); ok {
			return Delve[T](subSlice, subIndices...)
		} else {
			return OptionalPointer[T](nil)
		}
	}
}

//// if the optional is a collection, index into a slice or map collection
//func (o *Optional[T]) Get(index any) Optional[T] {
//	if o.rawPointer == nil {
//		return OptionalPointer[T](nil)
//	}
//
//	iVal := reflect.ValueOf(index)
//	rVal := reflect.ValueOf(*o.rawPointer)
//	indexKind := iVal.Kind()
//	switch rVal.Kind() {
//	case reflect.Slice:
//		if indexKind != reflect.Int {
//			return OptionalPointer[T](nil)
//		}
//		var indexIntVal int = iVal.Interface().(int)
//		var innerVal T = rVal.Index(indexIntVal).Interface().(T)
//		return OptionalValue[T](innerVal)
//	case reflect.Map:
//		innerVal := rVal.MapIndex(reflect.ValueOf(index)).Interface().(T)
//		return OptionalValue[T](innerVal)
//	default:
//		return OptionalPointer[T](nil)
//	}
//}

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
