package field

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Double float64

func DoubleFromBytes(b []byte, size uint32) Double {

	fp := FieldParserSingleton()
	// Skip past the control byte
	fp.offset++
	maxSize := uint32(8)
	var numBytes = make([]byte, maxSize)

	for i := uint32(0); i < maxSize; i++ {
		numBytes[maxSize-size+i] = b[fp.offset+i]
	}

	ui := Double(math.Float64frombits(binary.BigEndian.Uint64(numBytes)))
	fp.offset += size
	return ui

}

func (d Double) String() string {
	return fmt.Sprintf("%f", d)
}

func (d Double) Type() Type {
	return DoubleField
}

func (d Double) Bytes() []byte {

	b := make([]byte, 9)
	b[0] = 0b0110_0000

	binary.BigEndian.PutUint64(b[1:], math.Float64bits(float64(d)))

	return b
}
