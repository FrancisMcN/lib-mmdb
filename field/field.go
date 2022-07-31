package field

type Type int

const (
	PointerField            Type = 1
	StringField                  = 2
	DoubleField                  = 3
	BytesField                   = 4
	Uint16Field                  = 5
	Uint32Field                  = 6
	MapField                     = 7
	Int32Field                   = 8
	Uint64Field                  = 9
	Uint128Field                 = 10
	ArrayField                   = 11
	DataCacheContainerField      = 12
	EndMarkerField               = 13
	BooleanField                 = 14
	FloatField                   = 15
)

type Field interface {
	//FromBytes(b []byte) Field
	String() string
	Type() Type
	Bytes() []byte
}
