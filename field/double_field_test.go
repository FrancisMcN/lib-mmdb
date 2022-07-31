package field

import "testing"

// Test a new double from bytes is created successfully
func TestNewDoubleFromBytes(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(0)

	bytes := []byte{
		0b0110_1000,
		0b10111111,
		0b11110000,
		0b00000000,
		0b00000000,
		0b00000000,
		0b00000000,
		0b00000000,
		0b00000000,
	}
	d := DoubleFromBytes(bytes, 8)
	if d != -1 {
		t.Errorf("d = %f, test failed should be -1", d)
	}

}
