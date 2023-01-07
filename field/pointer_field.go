package field

import (
	"encoding/binary"
	"fmt"
)

type Pointer uint32

func PointerFromBytes(b []byte) Pointer {

	fp := FieldParserSingleton()
	offset := fp.offset

	pointerSize := (b[offset] & 0b0001_1000) >> 3
	// fmt.Println("pointer size is: ",pointerSize)
	switch pointerSize {
	case 1:
		bytes := make([]byte, 4)
		// Get last 3 bits of pointer
		// First byte stays zero, so start at 1
		bytes[1] = b[offset] & 0b0000_0111
		bytes[2] = b[offset+1]
		bytes[3] = b[offset+2]
		p := Pointer(binary.BigEndian.Uint32(bytes)) + 2_048
		// fmt.Println("pointer", fmt.Sprintf("%x", bytes), binary.BigEndian.Uint32(bytes))
		// fmt.Println(p, p - 2048, fmt.Sprintf("%x", bytes), fmt.Sprintf("%08b", bytes))
		return p
	case 2:
		bytes := make([]byte, 4)
		// Get last 3 bits of pointer
		// First byte stays zero, so start at 1
		bytes[0] = b[offset] & 0b0000_0111
		bytes[1] = b[offset+1]
		bytes[2] = b[offset+2]
		bytes[3] = b[offset+3]

		p := Pointer(binary.BigEndian.Uint32(bytes)) + 526_336
		//fmt.Println(p, bytes, fmt.Sprintf("%08b", bytes))
		return p
	case 3:
		bytes := make([]byte, 4)
		// Get last 3 bits of pointer
		// First byte stays zero, so start at 1
		bytes[0] = b[offset+1]
		bytes[1] = b[offset+2]
		bytes[2] = b[offset+3]
		bytes[3] = b[offset+4]
		p := Pointer(binary.BigEndian.Uint32(bytes))
		//fmt.Println(p, bytes, fmt.Sprintf("%08b", bytes))
		return p
	default:
		bytes := make([]byte, 2)
		// Get last 3 bits of pointer
		bytes[0] = b[offset] & 0b0000_0111
		bytes[1] = b[offset+1]
		p := Pointer(binary.BigEndian.Uint16(bytes))
		return p
	}

	return Pointer(0)

	//if pointer_size == 1:
	//buf = bytes([size & 0x7]) + buf
	//pointer = struct.unpack(b"!H", buf)[0] + self._pointer_base
	//elif pointer_size == 2:
	//buf = b"\x00" + bytes([size & 0x7]) + buf
	//pointer = struct.unpack(b"!I", buf)[0] + 2048 + self._pointer_base
	//print("p", pointer)
	//elif pointer_size == 3:
	//buf = bytes([size & 0x7]) + buf
	//pointer = struct.unpack(b"!I", buf)[0] + 526336 + self._pointer_base
	//else:
	//pointer = struct.unpack(b"!I", buf)[0] + self._pointer_base

	//ptrSize := (b[offset+0] & 0b0001_1000) >> 3
	//bytes := make([]byte, 4)
	//ptr := uint32(0)
	//fmt.Println("ptrSize", ptrSize)
	//switch ptrSize {
	//case 0:
	//	byte1 := b[offset] & 0b0000_0111
	//	byte2 := b[offset+1]
	//	bytes[2] = byte2
	//	bytes[3] = byte1
	//	ptr += uint32(byte1)
	//	ptr += uint32(byte2)
	//	//ptr = binary.BigEndian.Uint32(bytes)
	//	//fmt.Println("ptr", ptr, bytes, fmt.Sprintf("%x", b[offset:offset+4]), offset)
	//	//fmt.Println("ptr", ptr)
	//	//fmt.Println("bytes", bytes, offset, b[offset])
	//	//fmt.Println(fmt.Sprintf("%x", b[:10]))
	//	//fmt.Println("offset", fp.offset)
	//	//fmt.Println(b[offset], byte1, byte2)
	//	//fmt.Println("ptr", ptr)
	//	//fmt.Println(fmt.Sprintf("%x", b[:20]))
	//	fp.offset += 2
	//	//*data = (*data)[2:]
	//case 1:
	//	byte1 := b[offset] & 0b0000_0111
	//	byte2 := b[offset+1]
	//	byte3 := b[offset+2]
	//	bytes[1] = byte3
	//	bytes[2] = byte2
	//	bytes[3] = byte1
	//	ptr += uint32(byte1)
	//	ptr += uint32(byte2)
	//	ptr += uint32(byte3)
	//	ptr += uint32(2048)
	//
	//	//ptr = binary.BigEndian.Uint32(bytes) + 2048
	//	//*data = (*data)[3:]
	//	//fmt.Println("ptr", ptr)
	//	fp.offset += 3
	//case 2:
	//	byte1 := b[offset] & 0b0000_0111
	//	byte2 := b[offset+1]
	//	byte3 := b[offset+2]
	//	byte4 := b[offset+3]
	//	bytes[0] = byte4
	//	bytes[1] = byte3
	//	bytes[2] = byte2
	//	bytes[3] = byte1
	//	ptr += uint32(byte1)
	//	ptr += uint32(byte2)
	//	ptr += uint32(byte3)
	//	ptr += uint32(byte4)
	//	ptr += 526336
	//	//ptr = binary.BigEndian.Uint32(bytes) + 526336
	//	fp.offset += 4
	//	//*data = (*data)[4:]
	//default:
	//	ptr = binary.BigEndian.Uint32(b[1:5])
	//	fp.offset += 5
	//	//*data = (*data)[5:]
	//}

	//return Pointer(ptr)

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
	pointerSize := (b[fp.offset] & 0b0001_1000) >> 3
	// fmt.Println("pointer is", p)
	fp.SetOffset(uint32(p))
	// fmt.Println("offset", fp.offset, "p", p)
	f := fp.Parse(b)
	// pointerSize = 0, pointer requires 2 bytes
	// pointerSize = 1, pointer requires 3 bytes
	// pointerSize = 2, pointer requires 4 bytes
	// pointerSize = 3, pointer requires 5 bytes
	// fp.SetOffset(uint32(p) + )
	fp.SetOffset(off + uint32(pointerSize) + 2)
	return f
}

func (p Pointer) Bytes() []byte {

	//pointerSize := (uint32(p) & 0b0001_1000) >> 3
	//switch pointerSize {
	//case 1:
	// 11 bit pointer, pointerSize = 0, uses 2 bytes
	if p < 2_048 {

		b := make([]byte, 0)
		b = append(b, 0b0010_0111, 0b1111_1111)
		b2 := make([]byte, 2)
		binary.BigEndian.PutUint16(b2, uint16(p))
		for i, _ := range b2 {
			b2[i] &= b[i]
		}
		b2[0] |= 0b0010_0000
		//fmt.Println(fmt.Sprintf("%08b", b2))
		return b2
	}

	//case 2:
	// 19 bit pointer, pointerSize = 1, uses 3 bytes
	if p >= 2_048 && p < 526_336 {

		b := make([]byte, 0)
		b = append(b, 0, 0b0000_0111, 0b1111_1111, 0b1111_1111)
		b2 := make([]byte, 4)
		binary.BigEndian.PutUint32(b2, uint32(p)-2_048)
		for i, _ := range b2 {
			b2[i] &= b[i]
		}

		b2 = b2[1:]

		// fmt.Println("pointer val", fmt.Sprintf("%x", b2), p)

		b2[0] |= 0b0010_1000
		return b2
		// Remove middle byte
		// return append(b2[:1], b2[2:]...)
	}

	// 27 bit pointer, pointerSize = 2, uses 4 bytes
	if p >= 526_336 && p < 134_217_728 {

		b := make([]byte, 0)
		b = append(b, 0b0000_0111, 0b1111_1111, 0b1111_1111, 0b1111_1111)
		b2 := make([]byte, 4)
		binary.BigEndian.PutUint32(b2, uint32(p)-526_336)
		for i, _ := range b2 {
			b2[i] &= b[i]
		}
		b2[0] |= 0b0011_0000
		return b2
	}

	// 32 bit value, pointerSize = 3, uses 5 bytes
	// 1 byte to hold the pointer type and 4 bytes for the value
	b2 := make([]byte, 5)
	binary.BigEndian.PutUint32(b2[1:], uint32(p))
	b2 = append([]byte{0}, b2...)
	b2[0] |= 0b0011_1000
	return b2

	//default:
	//	fmt.Println("pointer size: ", pointerSize, p)
	//	b := make([]byte, 4)
	//	b[0] = 1
	//	return b
	//b := make([]byte, 0)
	//b = append(b, 0, 0b0010_0111, 0b1111_1111, 0b1111_1111)
	//b2 := make([]byte, 4)
	//binary.BigEndian.PutUint32(b2, uint32(p))
	//for i, _ := range b2 {
	//	b2[i] &= b[i]
	//}
	//return b2
	//}
	//// 11 bit pointer
	//if p <= 2_048 {
	//
	//}
	//// 19 bit pointer
	//if p < 526_336 {
	//
	//}
	//fmt.Println(b)
	//if p <= 255 {
	//	b = append(b, 0)
	//	if p < 7 {
	//		b[0] &= byte(p)
	//		b[0] |= 0b0010_0000
	//	} else {
	//		b[1] = byte(p) - 7
	//		fmt.Println(b[0], b[1])
	//	}
	//} else if p >= 256 && p <= 65_535 {
	//	b = append(b, 0, 0)
	//	binary.BigEndian.PutUint16(b[1:], uint16(p)-7)
	//	b[0] |= 0b0010_1000
	//} else if p >= 65_536 && p <= 16777215 {
	//	b = append(b, 0, 0, 0, 0)
	//	binary.BigEndian.PutUint32(b[1:], uint32(p)-7)
	//	b = append(b[:1], b[2:]...)
	//	b[0] |= 0b0011_0000
	//} else {
	//	b = append(b, 0, 0, 0, 0)
	//	binary.BigEndian.PutUint32(b[1:], uint32(p)-7)
	//	b[0] |= 0b0011_1000
	//}
	//fmt.Println("pointer value", p, "pointer bytes", b)
	//return b
}
