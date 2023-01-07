package field

import "testing"

// Test a new string from bytes is created successfully
func TestNewOneCharStringFromBytes(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(0)

	bytes := []byte{
		0b0100_0001,
		0b_110_0001,
	}
	s := StringFromBytes(bytes, 1)
	if s != "a" {
		t.Errorf("s = %s, test failed should be 'a'", s)
	}

}
