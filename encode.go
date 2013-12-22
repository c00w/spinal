package spinal

import (
    "hash"
)

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
