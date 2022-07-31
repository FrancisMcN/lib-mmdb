package field

import (
	"encoding/binary"
	"fmt"
)

type Int32 int32

func Int32FromBytes(b []byte, size uint32) Int32 {

	fp := FieldParserSingleton()
	// Skip past the control byte
	fp.offset++

	maxSize := uint32(4)
	var numBytes = make([]byte, maxSize)

	for i := uint32(0); i < size; i++ {
		numBytes[maxSize-size+i] = b[fp.offset+i]
	}

	ui := Int32(binary.BigEndian.Uint32(numBytes))
	fp.offset += size

	return ui

}

func (i Int32) String() string {
	return fmt.Sprintf("%d", i)
}

func (i Int32) Type() Type {
	return Int32Field
}

func (i Int32) Bytes() []byte {

	b := make([]byte, 2)
	b[0] = 0b0000_0000
	b[0] = 0b0000_0001
	if i == 0 {

	} else if i < 0 {
		b = append(b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(b[2:], uint32(i))
		b[2] |= 0b1000_0000
	} else if i >= 0 && i <= 255 {
		b = append(b, 0)
		b[2] = byte(i)
	} else if i >= 0 && i >= 256 && i <= 65_535 {
		b = append(b, 0, 0)
		binary.BigEndian.PutUint32(b[1:], uint32(i))
		b = append(b[:2], b[4:]...)
	} else if i >= 0 && i >= 65_536 && i <= 16777215 {
		b = append(b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(b[1:], uint32(i))
		b = append(b[:2], b[2:]...)
	} else if i >= 0 {
		b = append(b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(b[2:], uint32(i))
	}

	return b
}
