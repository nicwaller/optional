package optional

import (
	"github.com/nicwaller/optional/wrapper"
)

//goland:noinspection GoUnusedExportedFunction
func Optional[T any](ptr *T) wrapper.Wrapper[T] {
	return wrapper.New[T](ptr)
}
