package spinal

import (
    "testing"
    "hash/crc32"
)

func TestInOut(t *testing.T) {
    h := crc32.NewIEEE()
    k := 1
    d := 1
    B := 4
    input := []byte("Hello world! My name is betty jone")
    out, _ := Encode(k, h, input, len(input))
    resp := Decode(len(input), k, d, B, h, out)


    if len(resp) !=len(input) {
        t.Fatal("Length is incorrect")
    }

    if string(resp) != string(input) {
        t.Fatal(string(resp), "  !=  ", string(input))
    }

}
func TestInOutBigK(t *testing.T) {
    h := crc32.NewIEEE()
    k := 2
    d := 1
    B := 4
    input := []byte("Hello world! My name is betty jone")
    out, _ := Encode(k, h, input, len(input))
    resp := Decode(len(input), k, d, B, h, out)


    if len(resp) !=len(input) {
        t.Fatal("Length is incorrect")
    }

    if string(resp) != string(input) {
        t.Fatal(string(resp), "  !=  ", string(input))
    }

}
