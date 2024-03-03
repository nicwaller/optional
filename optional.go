package optional

import (
	"github.com/nicwaller/optional/wrapper"
)

//goland:noinspection GoUnusedExportedFunction
func NewOptional[T any](ptr *T) Optional[T] {
	return Optional[T](wrapper.New[T](ptr))
}

type Optional[T any] wrapper.Wrapper[T]
