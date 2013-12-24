package spinal

import (
	"hash"
)

type BufAdd struct {
	b       []byte
	started bool
}

func (b *BufAdd) Zero() {
	b.started = false
	for i, _ := range b.b {
		b.b[i] = 0
	}
}

func (b *BufAdd) Max() bool {
	if !b.started {
		return false
	}
	for _, c := range b.b {
		if c != 0 {
			return false
		}
	}
	return true
}

func (b *BufAdd) Increment() {
	b.started = true
	for i, _ := range b.b {
		b.b[i] += 1
		if b.b[i] != 0 {
			break
		}
	}
}

type RNG struct {
	i      uint32
	seed   []byte
	h      hash.Hash
	buffer []byte
}

func uint32tob(u uint32, b []byte) []byte {
	for i := 0; i < 4; i++ {
		b = append(b, byte(u>>uint(i)*8))
	}
	return b
}

func NewRNG(h hash.Hash, seed []byte) (r *RNG) {
	r = new(RNG)
	r.i = 3610617884
	r.h = h
	r.buffer = make([]byte, 0, h.Size())
	r.seed = seed
	return
}

func (r *RNG) Next() (b byte) {
	if len(r.buffer) > 0 {
		b = r.buffer[len(r.buffer)-1]
		r.buffer = r.buffer[:len(r.buffer)-1]
		return
	} else {
		r.h.Reset()
		r.h.Write(r.seed)
		r.h.Write(uint32tob(r.i, r.buffer))
		r.buffer = r.h.Sum(r.buffer[:0])
		r.i += 3243335647
		return r.Next()
	}
}

func HammingDistance(a byte, b byte) (d byte) {
	x := a ^ b

	for i := uint(0); i < 8; i++ {
		d += (x >> i) & 1
	}
	return
}

type SubTrees [][]decodeState

func (s SubTrees) Len() int {
	return len(s)
}

func (s SubTrees) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SubTrees) Less(i, j int) bool {
	return s[i][0].cost < s[j][0].cost
}

type MinCost []decodeState

func (m MinCost) Len() int {
	return len(m)
}

func (m MinCost) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MinCost) Less(i, j int) bool {
	return m[i].cost < m[j].cost
}
