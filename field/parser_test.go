package field

import (
	"fmt"
	"testing"
)

func TestNewStringField(t *testing.T) {

	fp := FieldParserSingleton()
	fp.SetOffset(0)

	bytes := []byte{
		0b0100_0001,
		0b_110_0001,
	}
	s := String("net:1.178.112.0/20, asn:AS12975").Bytes()
	fmt.Println(len(s))
	for _, c := range s {
		fmt.Printf("%08b ", c)
	}
	fmt.Println("")
	f := fp.Parse(bytes)
	if f != String("a") {
		t.Errorf("s = %s, test failed should be 'a'", f)
	}

}
