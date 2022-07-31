package field

import (
	"encoding/binary"
	"strings"
)

type Array []Field

func ArrayFromBytes(b []byte, length uint32) Array {

	fp := FieldParserSingleton()
	// Skip past the control byte
	fp.offset += 1
	//fmt.Println(b[fp.offset:fp.offset+10], fmt.Sprintf("%x", b[fp.offset:fp.offset+10]))
	//fmt.Println(fp.offset)
	//fmt.Println(fmt.Sprintf("%x", b[fp.offset:]))
	array := make([]Field, 0)
	for i := uint32(0); i < length; i++ {
		//fmt.Println(fmt.Sprintf("x = %x, b = %x", b[fp.offset:fp.offset+5], b[:]))
		//fmt.Println("test", fmt.Sprintf("%x", b[fp.offset:fp.offset+10]), b[fp.offset:fp.offset+10])
		f := fp.Parse(b)
		array = append(array, f)
	}
	return array

}

func (a Array) String() string {
	var sb strings.Builder
	sb.WriteString("[ ")
	for _, val := range a {
		sb.WriteString(val.String())
		sb.WriteString(" ")
	}
	sb.WriteString("]")
	return sb.String()
}

func (a Array) Type() Type {
	return ArrayField
}

func (a Array) Size() int {
	return len(a)
}

func (a Array) Get(i int) Field {
	return a[i]
}

func (a Array) Bytes() []byte {

	l := len(a)
	b := make([]byte, 2)
	b[0] = byte(l)
	b[1] = 0b0000_0100

	if l == 29 {
		b = append(b, 0)
		b[2] = byte(l - 29)
	} else if l == 30 {
		b = append(b, 0, 0)
		binary.BigEndian.PutUint16(b[2:], uint16(l)-285)
	} else if l == 31 {
		b = append(b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(b[2:], uint32(l)-65_821)
		b = append(b[:2], b[3:]...)
	}

	for _, elem := range a {
		b = append(b, elem.Bytes()...)
	}

	return b
}
