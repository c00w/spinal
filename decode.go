package spinal

import (
	"hash"
	"log"
	"sort"
)

type decodeState struct {
	cost       uint64
	lastSpline []byte
	message    []byte
}

func Decode(n int, k int, d int, B int, h hash.Hash, enc []byte) []byte {
	blockcount := n / k

	rngoutput := make([][]byte, blockcount)

	for i, c := range enc {
		rngoutput[i%blockcount] = append(rngoutput[i%blockcount], c)
	}

	states := make([]decodeState, 0, B*d*(1<<8*k))
	states = append(states, decodeState{0, make([]byte, h.Size()), nil})

	newstates := make([]decodeState, 0, B*d*(1<<8*k))

	edge := make([]byte, k)
	bedge := &BufAdd{edge, false}

	for i := 0; i < len(rngoutput); i++ {

		for _, state := range states {
			//log.Print("State: ", state)
			for bedge.Zero(); !bedge.Max(); bedge.Increment() {
				//log.Printf("Edge #%d", edge)
				h.Reset()
				h.Write(state.lastSpline)
				h.Write(edge)
				//h.Write([]byte{byte(edge)})
				spline := h.Sum(make([]byte, 0, h.Size()))
				rng := NewRNG(h, spline)
				d := uint64(0)
				for _, c := range rngoutput[i] {
					n := rng.Next()
					d += uint64(HammingDistance(c, n))
					if c != n && d == 0 {
						log.Print("Shenanigans")
					}
				}
				x := make([]byte, 0, (i+1)*k)
				x = append(x, state.message...)
				x = append(x, edge...)

				newstates = append(newstates, decodeState{d + state.cost, spline, x})
			}
		}

		childcount := 1 << uint(8*k*(d-1))
		subtrees := make([][]decodeState, 0, len(newstates)/childcount)

		for i := 0; i < len(newstates); i += childcount {
			subtrees = append(subtrees, newstates[i:i+childcount])
			sort.Sort(MinCost(subtrees[len(subtrees)-1]))
		}

		sort.Sort(SubTrees(subtrees))

		if len(subtrees) > B {
			subtrees = subtrees[0:B]
		}

		states = states[:0]

		for _, tree := range subtrees {
			states = append(states, tree...)
		}
		newstates = newstates[:0]
	}

	sort.Sort(MinCost(states))

	return states[0].message

}
