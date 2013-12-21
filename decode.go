package spinal

import (
    "hash"
    "log"
    "sort"
)

func HammingDistance(a byte, b byte) (d byte) {
    x := a ^ b

    for i := uint(0); i < 8; i++ {
        d += (x >> i) & 1
    }
    return
}

type decodeState struct {
    cost        uint64
    lastSpline  []byte
    message     []byte
}

type SubTrees [][]decodeState

func (s SubTrees) Len() int {
    return len(s)
}

func (s SubTrees) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s SubTrees) Less(i, j int) bool {
    min1 := s[i][0].cost
    for _, c := range s[i] {
        if c.cost < min1 {
            min1 = c.cost
        }
    }

    min2 := s[j][0].cost
    for _, c := range s[j] {
        if c.cost < min2 {
            min2 = c.cost
        }
    }
    return min1 < min2
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

func Decode(n int, k int, d int, B int, h hash.Hash, enc[]byte) []byte {
    blockcount := n/k

    rngoutput := make([][]byte, blockcount)

    for i, c := range enc {
        rngoutput[i%blockcount] = append(rngoutput[i%blockcount], c)
    }

    states := []decodeState{decodeState{0, make([]byte, h.Size()), nil}}

    log.Print("Starting Decode")
    for i:= 0; i < len(rngoutput); i++ {
        log.Printf("#States = %d", len(states))
        newstates := make([]decodeState, 0, len(states)*(1 << 8*k))

        for _, state := range states {
            //log.Print("State: ", state)
            for edge := uint64(0); edge < (1 << uint(8*k)); edge++ {
                //log.Printf("Edge #%d", edge)
                h.Reset()
                h.Write(state.lastSpline)
                h.Write([]byte{byte(edge)})
                spline := h.Sum(nil)
                rng := NewRNG(h, spline)
                d := uint64(0)
                for _, c := range rngoutput[i] {
                    n := rng.Next()
                    d += uint64(HammingDistance(c, n))
                    if c != n && d == 0 {
                        log.Print("Shenanigans")
                    }
                }

                x := make([]byte, len(state.message)+1)
                copy(x, state.message)
                x[len(state.message)] = byte(edge)

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

        states = make([]decodeState, 0)

        for _, tree := range subtrees {
            states = append(states, tree...)
        }
    }

    sort.Sort(MinCost(states))

    log.Print("Final Result:", states[0])
    return states[0].message

}
