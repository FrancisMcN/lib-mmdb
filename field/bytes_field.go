package field

type Bytes []byte

func BytesFromBytes(b []byte, size uint32) Bytes {

	fp := FieldParserSingleton()
	//fmt.Println("bytes")
	// Skip past the control byte
	fp.offset++

	//maxSize := uint32(4)
	var numBytes = make([]byte, size)

	for i := uint32(0); i < size; i++ {
		numBytes[i] = b[fp.offset+i]
	}

	bytes := Bytes(numBytes)
	fp.offset += size
	return bytes

}

func (b Bytes) String() string {
	return string(b)
}

func (b Bytes) Type() Type {
	return BytesField
}

func (b Bytes) Bytes() []byte {

	bytes := make([]byte, 1)
	bytes[0] = 0b1000_0000
	b = append(b, bytes...)
	return bytes
}
