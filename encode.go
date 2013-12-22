package spinal

import (
    "bytes"
    "encoding/binary"
    "hash"
)

type RNG struct {
    i    uint32
    seed []byte
    h    hash.Hash
    buffer []byte
}

func uint32tob(u uint32) []byte {
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.LittleEndian, u)
    return buf.Bytes()
}

func NewRNG(h hash.Hash, seed []byte) (r *RNG) {
    r = new(RNG)
    r.i = 3610617884
    r.h = h
    r.buffer = nil
    r.seed = seed
    return
}

func (r *RNG) Next() (b byte) {
    if len(r.buffer) > 0 {
        b = r.buffer[0]
        r.buffer = r.buffer[1:len(r.buffer)]
        return
    } else {
        r.h.Reset()
        r.h.Write(r.seed)
        r.h.Write(uint32tob(r.i))
        r.buffer = r.h.Sum(nil)
        r.i += 3243335647
        return r.Next()
    }
}

func Encode(k int, h hash.Hash, value []byte, output int) (enc []byte, s [][]byte) {
    s0 := make([]byte, h.Size())

    n := len(value)
    blockcount := n/k

    s = make([][]byte, blockcount)

    for i:= 0; i < blockcount; i++ {
        h.Reset()
        if i == 0 {
            h.Write(s0)
        } else {
            h.Write(s[i-1])
        }

        h.Write(value[i * k: (i+1)*k])
        s[i] = h.Sum(nil)
    }

    rng := make([](*RNG), len(s))
    for i := 0; i < len(s); i++ {
        rng[i] = NewRNG(h, s[i])
    }

    enc = make([]byte, output)

    for i := 0; i < output; i++ {
        enc[i] = rng[i%blockcount].Next()
    }

    return
}
