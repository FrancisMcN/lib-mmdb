package field

import (
	"encoding/binary"
	"fmt"
)

type Uint64 uint64

func Uint64FromBytes(b []byte, size uint32) Uint64 {
	fp := FieldParserSingleton()
	// Skip past the control byte
	fp.offset++
	maxSize := uint32(8)
	var numBytes = make([]byte, maxSize)

	for i := uint32(0); i < size; i++ {
		numBytes[maxSize-size+i] = b[fp.offset+i]
	}

	ui := Uint64(binary.BigEndian.Uint64(numBytes))
	fp.offset += size
	return ui

}

func (ui Uint64) String() string {
	return fmt.Sprintf("%d", ui)
}

func (ui Uint64) Type() Type {
	return Uint64Field
}

func (ui Uint64) Bytes() []byte {
	b := make([]byte, 2)
	b[0] = 0b0000_0000
	b[1] = 0b0000_0010

	if ui == 0 {

	} else if ui <= 255 {
		b[0] |= 0b0000_0001
		b = append(b, 0)
		b[2] = byte(ui)
	} else if ui >= 256 && ui <= 65_535 {
		b[0] |= 0b0000_0010
		b = append(b, 0, 0)
		binary.BigEndian.PutUint16(b[1:], uint16(ui))
	} else if ui >= 65_536 && ui <= 16_777_215 {
		b[0] |= 0b0000_0011
		b = append(b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(b[2:], uint32(ui))
		b = append(b[:1], b[2:]...)
	} else if ui >= 16_777_216 && ui <= 4_294_967_295 {
		b[0] |= 0b0000_0100
		b = append(b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(b[2:], uint32(ui))
	} else {
		b[0] |= 0b0000_1000
		b = append(b, 0, 0, 0, 0, 0, 0, 0, 0)
		binary.BigEndian.PutUint64(b[2:], uint64(ui))
	}

	return b
}
