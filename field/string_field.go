package field

import (
	"encoding/binary"
	"math"
)

type String string

func StringFromBytes(b []byte, length uint32) String {

	fp := FieldParserSingleton()
	// Skip past the control byte
	fp.offset++
	//fmt.Println(fmt.Sprintf("str: %x, %d", b[fp.offset:fp.offset+10], length))

	s := String(b[fp.offset : fp.offset+length])
	fp.offset += length

	return s

}

func (s String) String() string {
	return string(s)
}

func (s String) Type() Type {
	return StringField
}

func (s String) Bytes() []byte {

	l := len(s)
	b := make([]byte, 1)
	b[0] = byte(math.Min(float64(l), float64(31)))
	b[0] &= 0b0101_1111
	b[0] |= 0b0100_0000

	if l >= 32 {
		if l >= 29 && l <= 284 {
			b = append(b, byte(l-29))
		} else if l >= 285 && l <= 65_821 {
			b = append(b, 0, 0)
			binary.BigEndian.PutUint16(b[1:2], uint16(l-285))
		} else {
			b = append(b, 0, 0, 0, 0)
			binary.BigEndian.PutUint32(b[1:3], uint32(l-65_821))
			b = append(b[:1], b[2:]...)
		}
	}

	b = append(b, []byte(s)...)
	return b
}
