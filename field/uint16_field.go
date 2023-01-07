package field

import (
	"encoding/binary"
	"fmt"
)

type Uint16 uint16

func Uint16FromBytes(b []byte, size uint32) Uint16 {

	fp := FieldParserSingleton()
	// Skip past the control byte
	fp.offset++

	maxSize := uint32(2)
	var numBytes = make([]byte, maxSize)

	for i := uint32(0); i < size; i++ {
		numBytes[maxSize-size+i] = b[fp.offset+i]
	}

	ui := Uint16(binary.BigEndian.Uint16(numBytes))
	fp.offset += size
	return ui

}

func (ui Uint16) String() string {
	return fmt.Sprintf("%d", ui)
}

func (ui Uint16) Type() Type {
	return Uint16Field
}

func (ui Uint16) Bytes() []byte {
	b := make([]byte, 1)
	b[0] = 0b1010_0000
	if ui == 0 {

	} else if ui <= 255 {
		b = append(b, 0)
		b[0] |= 0b0000_0001
		b[1] = byte(ui)
	} else {
		b[0] |= 0b0000_0010
		b = append(b, 0, 0)
		binary.BigEndian.PutUint16(b[1:], uint16(ui))
	}
	return b
}
