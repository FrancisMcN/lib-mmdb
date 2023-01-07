package field

import (
	"encoding/binary"
	"fmt"
	// "sort"
	"strings"
)

type Map struct {
	InternalMap map[Field]Field
	OrderedKeys []Field
}

func NewMap() *Map {
	return &Map{
		InternalMap: make(map[Field]Field),
		OrderedKeys: make([]Field, 0),
	}
}

func MapFromBytes(b []byte, items uint32) *Map {

	fp := FieldParserSingleton()
	// fmt.Println("map offset: ", fp.offset)
	// fmt.Println(fmt.Sprintf("x: %x", b[int(math.Min(float64(fp.offset-5), 0)):fp.offset + 5]))
	// Skip past the control byte
	fp.offset += 1

	m := NewMap()
	// m := make(map[Field]Field)
	//fmt.Println("--- map ---")
	for i := uint32(0); i < items; i++ {
		//fmt.Println(fp.offset, fmt.Sprintf("%x", b[fp.offset:fp.offset+10]))
		//fmt.Println("fp.offset", fp.offset)
		//fmt.Println(m)
		key := fp.Parse(b)
		// fmt.Println("key", key, "b", fmt.Sprintf("%x", b[fp.offset:fp.offset+5]))
		//fmt.Println("key", key)
		//if key.Type() == PointerField {
		//	//key = key.(Pointer).Resolve(b)
		//}
		//fmt.Println(fp.offset)
		//fmt.Println("-----")
		//fmt.Println(fp.offset, fmt.Sprintf("%x", b[fp.offset:int(math.Min(float64(len(b)), float64(fp.offset+10)))]))
		val := fp.Parse(b)
		// fmt.Println("val", val, "b", fmt.Sprintf("%x", b[fp.offset:fp.offset+5]))
		//fmt.Println("val", val)
		//if val.Type() == PointerField {
		//	//val = val.(Pointer).Resolve(b)
		//}
		// InternalMap[key.(String)] = val
		m.Put(key, val)
		// OrderedKeys = append(OrderedKeys, key)
		//fmt.Println(m)
	}
	//fmt.Println("--- --- ---")
	return m

}

func (m *Map) String() string {
	// var sb strings.Builder
	kv := make([]string, 0)
	//i := 0
	for _, k := range m.OrderedKeys {

		v := m.InternalMap[k]
		kv = append(kv, fmt.Sprintf("%s:%s", k, v))
		//if k.Type() == PointerField {
		//	//k = k.(Pointer).Resolve()
		//}
		//sb.WriteString(fmt.Sprintf("%s:%s", k, v))
		//i++
		//if i < len(m) {
		//	sb.WriteString(" ")
		//}
	}
	// sort.Strings(kv)
	// fmt.Println("kv", kv)
	// fmt.Println("join:", ))
	// sb.WriteString("[ ")
	// for i, v := range kv {
	// 	sb.WriteString(v)
	// 	i++
	// 	if i < len(m.InternalMap) {
	// 		sb.WriteString(" ")
	// 	}
	// }
	// sb.WriteString(" ]")
	return fmt.Sprintf("[%s]", strings.Join(kv, " "))
}

func (m *Map) Type() Type {
	return MapField
}

func (m *Map) Size() int {
	return len(m.InternalMap)
}

func (m *Map) Get(key Field) Field {
	return m.InternalMap[key.(String)]
}

func (m *Map) Put(key Field, val Field) {
	if _, f := m.InternalMap[key]; f {
		m.InternalMap[key] = val
	} else {
		m.InternalMap[key] = val
		m.OrderedKeys = append(m.OrderedKeys, key)
	}
}

func (m *Map) Bytes() []byte {

	l := len(m.InternalMap)
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

	// keys := make([]string, 0)
	// keyvals := make([]Field, 0)
	// for k, _ := range m {
	// 	keys = append(keys, k.String())
	// 	keyvals = append(keyvals, k)
	// }
	// sort.Strings(keys)
	// // sort.Ints(keys)
	// // fmt.Println("sorted", keys)
	// for i, _ := range keys {
	// k := keyvals[i]
	// b = append(b, k.Bytes()...)
	// v := m[k]
	// b = append(b, v.Bytes()...)
	// }

	for _, k := range m.OrderedKeys {
		v := m.InternalMap[k]
		b = append(b, k.Bytes()...)
		b = append(b, v.Bytes()...)
	}
	// for k, v := range m {
	// 	b = append(b, k.Bytes()...)
	// 	b = append(b, v.Bytes()...)
	// }

	return b
}
