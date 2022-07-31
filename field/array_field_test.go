package field

import (
	"testing"
)

// Test a new size 0 array is created successfully
func TestNewSize0Array(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(1)
	bytes := []byte{
		0b0000_0000,
		0b0000_1011,
	}
	array := ArrayFromBytes(bytes, 0)
	if array.Size() != 0 {
		t.Errorf("array should have size 0")
	}

}

// Test a new size 1 array is created successfully
func TestNewSize1Array(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(1)
	bytes := []byte{
		0b0000_0001,
		0b0000_0100,
		0b0100_0011,
		0b0110_0001,
		0b0110_0010,
		0b0110_0011,
	}
	array := ArrayFromBytes(bytes, 1)
	if array.Size() != 1 {
		t.Errorf("array should have size 1")
	}
	//fmt.Println(array)
	if array.Get(0).String() != "abc" {
		t.Errorf("array element should be 'abc'")
	}

}

// Test a new size 2 array is created successfully
func TestNewSize2Array(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(1)
	bytes := []byte{
		0b0000_0010,
		0b0000_1011,
		0b0100_0011,
		0b0110_0001,
		0b0110_0010,
		0b0110_0011,
		0b0100_0011,
		0b0111_1000,
		0b0111_1001,
		0b0111_1010,
	}

	array := ArrayFromBytes(bytes, 2)

	if array.Size() != 2 {
		t.Errorf("array should have size 2")
	}

	if array.Get(0).String() != "abc" {
		t.Errorf("array element should be 'abc'")
	}

	if array.Get(1).String() != "xyz" {
		t.Errorf("array element should be 'xyz'")
	}

}
