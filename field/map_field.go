package field

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Map map[Field]Field

func MapFromBytes(b []byte, items uint32) Map {

	fp := FieldParserSingleton()
	//fmt.Println(fmt.Sprintf("x: %x", b[fp.offset]))
	// Skip past the control byte
	fp.offset += 1

	m := make(map[Field]Field)
	//fmt.Println("--- map ---")
	for i := uint32(0); i < items; i++ {
		//fmt.Println(fp.offset, fmt.Sprintf("%x", b[fp.offset:fp.offset+10]))
		key := fp.Parse(b)
		//fmt.Println("key", key)
		//if key.Type() == PointerField {
		//	//key = key.(Pointer).Resolve(b)
		//}
		//fmt.Println(fp.offset)
		//fmt.Println("-----")
		//fmt.Println(fp.offset, fmt.Sprintf("%x", b[fp.offset:int(math.Min(float64(len(b)), float64(fp.offset+10)))]))
		val := fp.Parse(b)
		//fmt.Println("val", val)
		//if val.Type() == PointerField {
		//	//val = val.(Pointer).Resolve(b)
		//}
		//fmt.Println(key, val)
		m[key.(String)] = val
	}
	//fmt.Println("--- --- ---")
	return m

}

func (m Map) String() string {
	var sb strings.Builder
	sb.WriteString("[ ")
	i := 0
	for k, v := range m {
		//if k.Type() == PointerField {
		//	//k = k.(Pointer).Resolve()
		//}
		sb.WriteString(fmt.Sprintf("%s:%s", k, v))
		i++
		if i < len(m) {
			sb.WriteString(" ")
		}
	}
	sb.WriteString(" ]")
	return sb.String()
}

func (m Map) Type() Type {
	return MapField
}

func (m Map) Size() int {
	return len(m)
}

func (m Map) Get(key Field) Field {
	return m[key.(String)]
}

func (m Map) Bytes() []byte {

	l := len(m)
	b := make([]byte, 1)
	b[0] = 0b1110_0000
	if l < 29 {
		b[0] |= byte(l)
	}
	if l == 29 {
		b = append(b, 0)
		b[1] = byte(l - 29)
	} else if l == 30 {
		b = append(b, 0, 0)
		binary.BigEndian.PutUint16(b[1:], uint16(l)-285)
	} else if l == 31 {
		b = append(b, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(b[1:], uint32(l)-65_821)
		b = append(b[:1], b[2:]...)
	}

	for k, v := range m {
		b = append(b, k.Bytes()...)
		b = append(b, v.Bytes()...)
		//fmt.Println(k, v)
	}

	return b
}
