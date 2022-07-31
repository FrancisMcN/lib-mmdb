package field

import (
	"testing"
)

// Test a new size 0 pointer is created successfully
func TestNewSizeZeroPointer(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(0)
	bytes := []byte{
		0b0010_0000,
		0b0000_1010,
	}
	p := PointerFromBytes(bytes)
	if p != 10 {
		t.Errorf("p = %d, test failed should be 10", p)
	}

}

// Test a new size one pointer is created successfully
func TestNewSizeOnePointer(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(0)
	bytes := []byte{
		0b0010_1000,
		0b0000_0001,
		0b0000_0010,
	}
	p := PointerFromBytes(bytes)
	if p != 2_306 {
		t.Errorf("p = %d, test failed should be 2,306", p)
	}

}

// Test a new size two pointer is created successfully
func TestNewSizeTwoPointer(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(0)
	bytes := []byte{
		0b0011_0000,
		0b0000_0001,
		0b0000_0001,
		0b0000_0001,
	}
	p := PointerFromBytes(bytes)
	if p != 592_129 {
		t.Errorf("p = %d, test failed should be 592,129", p)
	}

}

// Test a new size three pointer is created successfully
func TestNewSizeThreePointer(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(0)
	bytes := []byte{
		0b0011_1000,
		0b0000_0001,
		0b0000_0001,
		0b0000_0001,
		0b0000_0001,
	}
	p := PointerFromBytes(bytes)
	if p != 16_843_009 {
		t.Errorf("p = %d, test failed should be 16,843,009", p)
	}

}
