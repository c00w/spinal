package spinal

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

func PadBytes(k int, bytes []byte) []byte {
    if len(bytes) < k {
        add := k-len(bytes)
        //var padding [add]byte
        padding := make([]byte, add)
        n_bytes := append(padding, bytes...)
        return n_bytes
    }
    return bytes
}