package field

import (
	"encoding/binary"
)

type FieldParser struct {
	offset         uint32
	nodeCount      uint32
	searchTreeSize int
}

var fieldParser *FieldParser

func NewFieldParser(nodeCount uint32, searchTreeSize int) *FieldParser {
	fieldParser = &FieldParser{
		nodeCount:      nodeCount,
		searchTreeSize: searchTreeSize,
	}
	return fieldParser
}

func FieldParserSingleton() *FieldParser {
	return fieldParser
}

func (fp *FieldParser) Reset() {
	fp.offset = 0
}

func (fp *FieldParser) SetOffset(o uint32) {
	fp.offset = o
}

func (fp *FieldParser) GetOffset() uint32 {
	return fp.offset
}

func (fp *FieldParser) Parse(b []byte) Field {
	//fmt.Println(fmt.Sprintf("b0: %x, offset: %d", b[fp.offset + 0], fp.offset))
	// fmt.Println("off", fp.offset)
	// fmt.Println("off - nodeCount", fp.offset - fp.nodeCount)
	// fmt.Println("off + nodeCount", fp.offset + fp.nodeCount)
	// fmt.Println("off + nodeCount + 16", fp.offset + fp.nodeCount + 16)
	fieldType, size, off := getFieldTypeAndSize(b[fp.offset:])
	//fmt.Println(fmt.Sprintf("%x", b[fp.offset:fp.offset+10]), fieldType, size, off)
	fp.offset += off
	//fmt.Println("off", off)
	//fmt.Println("offset", fp.offset)
	//fmt.Println("type", fieldType)
	var field Field
	switch fieldType {
	case PointerField:
		//fmt.Println("fp", fp.offset)
		field = PointerFromBytes(b)
		// fmt.Println("pointer", field)
		// fmt.Println("fp", fmt.Sprintf("%x %x %x %x", b[fp.offset], b[fp.offset+1], b[fp.offset+2], b[fp.offset+ 3]))
		// fmt.Println(fmt.Sprintf("%08b hex: %x", b[fp.offset-5:fp.offset+5], b[fp.offset-5:fp.offset+5]))
		// fmt.Println(fmt.Sprintf("%08b", b[fp.offset:fp.offset+5]))
		// fmt.Println(fmt.Sprintf("%x", b[fp.offset:fp.offset+5]))
		//fmt.Println("field", field, uint32(field.(Pointer)))
		// fmt.Println(fp.offset, field)
		// fmt.Println("----")
		// fmt.Println("Pointer is", field)
		field = field.(Pointer).Resolve(b)
		//fmt.Println("resolved pointer is", field)
	case StringField:
		//fmt.Println("string", fmt.Sprintf("%x %d", b[fp.offset:fp.offset+10], size))
		field = StringFromBytes(b, size)
	case DoubleField:
		//fmt.Println("double", fmt.Sprintf("%x %d", b[fp.offset:fp.offset+10], size))
		field = DoubleFromBytes(b, size)
	case BytesField:
		field = BytesFromBytes(b, size)
	case Uint16Field:
		field = Uint16FromBytes(b, size)
	case Uint32Field:
		field = Uint32FromBytes(b, size)
	case MapField:
		field = MapFromBytes(b, size)
	case Int32Field:
		field = Int32FromBytes(b, size)
	case Uint64Field:
		field = Uint64FromBytes(b, size)
	//case Uint128Field:
	//	field = NewUint128(data, fieldSize)
	case ArrayField:
		field = ArrayFromBytes(b, size)
	//case DataCacheContainerField:
	//	field = NewDataCacheContainer(data)
	//case EndMarkerField:
	//	field = NewEndMarker(data)
	case BooleanField:
		field = BoolFromBytes(b, size)
	case FloatField:
		field = FloatFromBytes(b, size)
	}
	//if field == nil {
	//	fmt.Println(fmt.Sprintf("%x %d %d %d %d", b[fp.offset:fp.offset+25], fieldType, size, off, fp.offset))
	//}
	return field

}

func getFieldTypeAndSize(b []byte) (Type, uint32, uint32) {

	fieldType := Type(b[0] >> 5)
	fieldSize := uint32(b[0] & 0b0001_1111)
	extended := fieldType == 0
	offset := uint32(0)

	if fieldType != PointerField && fieldType != MapField {
		if extended {
			// If we reach here then the field's type is in the extended byte,
			// the spec says to subtract 7 from the second byte to find the correct type
			fieldType = Type(b[1] + 7)
			offset += 1
		}
		if fieldSize == 29 {
			offset += 1
			fieldSize = 29 + uint32(b[offset])
			//fmt.Println("fs", fieldSize, fmt.Sprintf("%x", b[offset]))
		} else if fieldSize == 30 {
			offset += 2
			fieldSize = 285 + binary.BigEndian.Uint32(b[offset:2])
		} else if fieldSize == 31 {
			offset += 3
			fieldSize = 65_821 + binary.BigEndian.Uint32(b[offset:3])
		}
	}
	//fmt.Println(fieldType, fieldSize, offset)
	return fieldType, fieldSize, offset
}
