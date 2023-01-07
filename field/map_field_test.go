package field

import (
	"testing"
)

// Test a new empty map is created successfully
func TestNewEmptyMap(t *testing.T) {

	fp := NewFieldParser(0, 0)
	fp.SetOffset(0)
	bytes := []byte{
		0b1110_0000,
	}
	m := MapFromBytes(bytes, 0)
	if m.Size() != 0 {
		t.Errorf("map should be empty, size is %d", m.Size())
	}

}

// Test a new map with one element is created successfully
func TestNewMapWithOneElement(t *testing.T) {

	fp := NewFieldParser(0, 0)
	fp.SetOffset(0)
	bytes := []byte{
		0b1110_0001,
		0b0100_0011,
		0b0110_0001,
		0b0110_0010,
		0b0110_0011,
		0b0100_0011,
		0b0110_0100,
		0b0110_0101,
		0b0110_0110,
	}
	m := MapFromBytes(bytes, 1)
	if m.Size() != 1 {
		t.Errorf("map should have one element, size is %d", m.Size())
	}
	if m.Get(String("abc")) == nil {
		t.Error("key can't be nil")
	}

}
