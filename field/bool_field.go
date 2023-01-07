package field

import (
	"fmt"
)

type Bool bool

func BoolFromBytes(b []byte, size uint32) Bool {

	var val bool
	fp := FieldParserSingleton()
	if size > 0 {
		val = true
	}
	fp.offset += 1
	return Bool(val)
}

func (b Bool) Type() Type {
	return BooleanField
}

func (b Bool) String() string {
	return fmt.Sprintf("%t", b)
}

func (b Bool) Bytes() []byte {

	bytes := make([]byte, 2)
	bytes[0] = 0b0000_0000
	bytes[1] = 0b0000_0111

	if b {
		bytes[0] = 1
	}

	return bytes
}
