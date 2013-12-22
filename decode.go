package spinal

import (
    "hash"
    "log"
    "sort"
    "math/big"
)

type decodeState struct {
    cost        uint64
    lastSpline  []byte
    message     []byte
}

func Decode(n int, k int, d int, B int, h hash.Hash, enc[]byte) []byte {
    blockcount := n/k

    rngoutput := make([][]byte, blockcount)

    for i, c := range enc {
        rngoutput[i%blockcount] = append(rngoutput[i%blockcount], c)
    }

    states := make([]decodeState, 0, B*d*(1 << 8*k))
    states = append(states, decodeState{0, make([]byte, h.Size()), nil})

    newstates := make([]decodeState, 0, B*d*(1 << 8*k))

    log.Print("Starting Decode")
    for i:= 0; i < len(rngoutput); i++ {
        log.Printf("#States = %d", len(states))

        for _, state := range states {
            //log.Print("State: ", state)
            max := big.NewInt(0).Exp(big.NewInt(2),big.NewInt(int64(8*k)), nil)
            for edge := big.NewInt(0); edge.Cmp(max) < 0; edge.Add(edge,big.NewInt(1)) {
                //log.Printf("Edge #%d", edge)
                h.Reset()
                h.Write(state.lastSpline)
                h.Write(PadBytes(k, edge.Bytes()))
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
                // TODO: Figure out why the below doesn't work. "[:]" should make it a slice. :(
                //x := append([]byte(state.message), byte(edge))[:]
                n_message := append(state.message, PadBytes(k, edge.Bytes())...)
                x := make([]byte, len(n_message))
                copy(x, n_message)

                newstates = append(newstates, decodeState{d + state.cost, spline, x})
            }
        }

        log.Printf("Exploded to %d newstates", len(newstates))

        childcount := 1 << uint(8*k * (d-1))
        subtrees := make([][]decodeState, 0, len(newstates)/childcount)

        for i:=0;i < len(newstates); i += childcount{
            subtrees = append(subtrees, newstates[i:i+childcount])
        }

        sort.Sort(SubTrees(subtrees))

        if len(subtrees) > B {
            sort.Sort(MinCost(subtrees[B]))
            sort.Sort(MinCost(subtrees[B-1]))
            if subtrees[B][0].cost == subtrees[B-1][0].cost {
                log.Print("Subtree collision. This might be a problem")
            }
            subtrees = subtrees[0:B]
        }

        states = states[:0]

        for _, tree := range subtrees {
            states = append(states, tree...)
        }
        newstates = newstates[:0]
    }

    sort.Sort(MinCost(states))

    log.Print("Final Result:", states[0])
    return states[0].message

}
