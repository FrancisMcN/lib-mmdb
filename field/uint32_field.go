package field

import (
	"encoding/binary"
	"fmt"
)

type Uint32 uint32

func Uint32FromBytes(b []byte, size uint32) Uint32 {
	fp := FieldParserSingleton()
	// Skip past the control byte
	fp.offset++

	maxSize := uint32(4)
	var numBytes = make([]byte, maxSize)

	for i := uint32(0); i < size; i++ {
		numBytes[maxSize-size+i] = b[fp.offset+i]
	}

	ui := Uint32(binary.BigEndian.Uint32(numBytes))
	fp.offset += size
	return ui

}

func (ui Uint32) String() string {
	return fmt.Sprintf("%d", ui)
}

func (ui Uint32) Type() Type {
	return Uint32Field
}

func (ui Uint32) Bytes() []byte {
	b := make([]byte, 1)
	b[0] = 0b1100_0000
	if ui == 0 {

	} else if ui <= 255 {
		b[0] |= 0b0000_0001
		b = append(b, 0)
		b[1] = byte(ui)
	} else if ui >= 256 && ui <= 65_535 {
		b[0] |= 0b0000_0010
		b = append(b, 0, 0)
		binary.BigEndian.PutUint16(b[1:], uint16(ui))
	} else if ui >= 65_536 && ui <= 16777215 {
		b[0] |= 0b0000_0011
		b = append(b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(b[1:], uint32(ui))
		b = append(b[:1], b[2:]...)
	} else {
		b[0] |= 0b0000_0100
		b = append(b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(b[1:], uint32(ui))
	}

	return b
}
