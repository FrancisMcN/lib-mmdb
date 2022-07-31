package field

import "testing"

// Test a new uint16 from bytes is created successfully
func TestNewUint16FromBytes(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(0)

	bytes := []byte{
		0b1010_0001,
		0b0000_0010,
	}
	i := Uint16FromBytes(bytes, 1)
	if i != 2 {
		t.Errorf("i = %d, test failed should be 2", i)
	}

}

// Test a new uint16 to bytes is created successfully
func TestUint16ToBytes(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(0)

	bytes := Uint16(24).Bytes()
	expected := []byte{
		0b1010_0001,
		0b0001_1000,
	}
	for i, _ := range bytes {
		if bytes[i] != expected[i] {
			t.Errorf("got %x expected %x", bytes[i], expected[i])
		}
	}

}
