package constraints

import (
	"golang.org/x/exp/constraints"
	"unsafe"
)

type String interface{ string | []byte | []rune }

type Hashable interface {
	constraints.Integer |
		constraints.Float |
		constraints.Complex |
		~string |
		uintptr |
		~unsafe.Pointer
}
