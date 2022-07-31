package field

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Float float32

func FloatFromBytes(b []byte, size uint32) Float {

	fp := FieldParserSingleton()
	// Skip past the control byte
	fp.offset++

	maxSize := uint32(4)
	var numBytes = make([]byte, maxSize)

	for i := uint32(0); i < maxSize; i++ {
		numBytes[maxSize-size+i] = b[fp.offset+i]
	}

	ui := Float(float32(binary.BigEndian.Uint64(numBytes)))
	fp.offset += size
	return ui

}

func (f Float) String() string {
	return fmt.Sprintf("%f", f)
}

func (f Float) Type() Type {
	return FloatField
}

func (f Float) Bytes() []byte {

	b := make([]byte, 6)
	b[0] = 0b0000_0000
	b[1] = 0b0000_1000

	binary.BigEndian.PutUint32(b[2:], math.Float32bits(float32(f)))

	return b
}
