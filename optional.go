package optional

import (
	"github.com/nicwaller/optional/wrapper"
)

//goland:noinspection GoUnusedExportedFunction
func NewOptional[T any](ptr *T) wrapper.Wrapper[T] {
	return wrapper.New[T](ptr)
}

type Optional[T any] wrapper.Wrapper[T]
