package field

import (
	"encoding/binary"
	"fmt"
)

type Pointer uint32

func PointerFromBytes(b []byte) Pointer {

	fp := FieldParserSingleton()
	offset := fp.offset
	ptrSize := (b[offset+0] & 0b0001_1000) >> 3
	bytes := make([]byte, 4)
	ptr := uint32(0)
	switch ptrSize {
	case 0:
		byte1 := b[offset] & 0b0000_0111
		byte2 := b[offset+1]
		bytes[2] = byte1
		bytes[3] = byte2
		ptr = binary.BigEndian.Uint32(bytes)
		//fmt.Println(fmt.Sprintf("%x", b[:20]))
		fp.offset += 2
		//*data = (*data)[2:]
	case 1:
		byte1 := b[offset] & 0b0000_0111
		byte2 := b[offset+1]
		byte3 := b[offset+2]
		bytes[1] = byte1
		bytes[2] = byte2
		bytes[3] = byte3
		ptr = binary.BigEndian.Uint32(bytes) + 2048
		//*data = (*data)[3:]
		fp.offset += 3
	case 2:
		byte1 := b[offset] & 0b0000_0111
		byte2 := b[offset+1]
		byte3 := b[offset+2]
		byte4 := b[offset+3]
		bytes[0] = byte1
		bytes[1] = byte2
		bytes[2] = byte3
		bytes[3] = byte4
		ptr = binary.BigEndian.Uint32(bytes) + 526336
		fp.offset += 4
		//*data = (*data)[4:]
	default:
		ptr = binary.BigEndian.Uint32(b[1:5])
		fp.offset += 5
		//*data = (*data)[5:]
	}
	return Pointer(ptr)

}

func (p Pointer) String() string {
	return fmt.Sprintf("&%d", p)
}

func (p Pointer) Type() Type {
	return PointerField
}

func (p Pointer) Resolve(b []byte) Field {
	fp := FieldParserSingleton()
	off := fp.offset
	fp.SetOffset(uint32(p))
	f := fp.Parse(b)
	fp.SetOffset(off)
	return f
}

func (p Pointer) Bytes() []byte {

	b := make([]byte, 1)
	b[0] = 0b0010_0000
	b[0] |= 0b0000_0111
	if p <= 255 {
		b = append(b, 0)
		if p < 7 {
			b[0] &= byte(p)
			b[0] |= 0b0010_0000
		} else {
			b[1] = byte(p) - 7
		}
	} else if p >= 256 && p <= 65_535 {
		b = append(b, 0, 0)
		binary.BigEndian.PutUint16(b[1:], uint16(p)-7)
		b[0] |= 0b0010_1000
	} else if p >= 65_536 && p <= 16777215 {
		b = append(b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(b[1:], uint32(p)-7)
		b = append(b[:1], b[2:]...)
		b[0] |= 0b0011_0000
	} else {
		b = append(b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(b[1:], uint32(p)-7)
		b[0] |= 0b0011_1000
	}

	return b
}
