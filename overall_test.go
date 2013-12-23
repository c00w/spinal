package spinal

import (
	"hash/crc32"
	"testing"
)

func TestInOut(t *testing.T) {
    testcases := [][]byte{
        []byte("foo"),
        []byte("Hello World!"),
        []byte("Hello World! My name is betty jones")}
    for _, input := range testcases {
        h := crc32.NewIEEE()
        k := 1
        d := 1
        B := 32
        out, _ := Encode(k, h, input, len(input)*2)
        resp := Decode(len(input), k, d, B, h, out)

        if len(resp) != len(input) {
            t.Fatal("Length is incorrect")
        }

        if string(resp) != string(input) {
            t.Fatal(string(resp), "  !=  ", string(input))
        }
    }

}
func TestInOutBigK(t *testing.T) {
	h := crc32.NewIEEE()
	k := 2
	d := 1
	B := 4
	input := []byte("Hello world! My name is betty jone")
	out, _ := Encode(k, h, input, len(input)*2)
	resp := Decode(len(input), k, d, B, h, out)

	if len(resp) != len(input) {
		t.Fatal("Length is incorrect")
	}

	if string(resp) != string(input) {
		t.Fatal(string(resp), "  !=  ", string(input))
	}

}
