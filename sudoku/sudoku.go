package sudoku

import "math/rand/v2"

const (
	N         = 9
	N2        = N * N
	N3        = N / 3
	MinGivens = 17
	MaxGivens = 81
)

type Board [N2]int
type Candidates [N2]uint16

var (
	squares        []int
	units          [][]int
	unitsForSquare [][]int
	peers          [][]int
)

func init() {
	squares = make([]int, N2)
	for i := range squares {
		squares[i] = i
	}
	units = getAllUnits()
	unitsForSquare = getUnitsForSquare()
	peers = getPeers()
}

func index(row, col int) int { return row*N + col }

func getAllUnits() [][]int {
	units := make([][]int, 0, 27)
	for r := 0; r < N; r++ {
		unit := make([]int, N)
		for c := 0; c < N; c++ {
			unit[c] = index(r, c)
		}
		units = append(units, unit)
	}
	for c := 0; c < N; c++ {
		unit := make([]int, N)
		for r := 0; r < N; r++ {
			unit[r] = index(r, c)
		}
		units = append(units, unit)
	}
	for br := 0; br < N; br += N3 {
		for bc := 0; bc < N; bc += N3 {
			unit := make([]int, N)
			idx := 0
			for r := br; r < br+N3; r++ {
				for c := bc; c < bc+N3; c++ {
					unit[idx] = index(r, c)
					idx++
				}
			}
			units = append(units, unit)
		}
	}
	return units
}

func getUnitsForSquare() [][]int {
	result := make([][]int, N2)
	for i := 0; i < N2; i++ {
		var sqUnits []int
		for ui, unit := range units {
			for _, cell := range unit {
				if cell == i {
					sqUnits = append(sqUnits, ui)
					break
				}
			}
		}
		result[i] = sqUnits
	}
	return result
}

func getPeers() [][]int {
	result := make([][]int, N2)
	for i := 0; i < N2; i++ {
		set := make(map[int]bool)
		for _, ui := range unitsForSquare[i] {
			for _, cell := range units[ui] {
				if cell != i {
					set[cell] = true
				}
			}
		}
		ps := make([]int, 0, len(set))
		for p := range set {
			ps = append(ps, p)
		}
		result[i] = ps
	}
	return result
}

func newCandidates() Candidates {
	var c Candidates
	allBits := uint16(0)
	for d := 0; d < N; d++ {
		allBits |= 1 << d
	}
	for i := 0; i < N2; i++ {
		c[i] = allBits
	}
	return c
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func shuffle(seq []int) []int {
	n := len(seq)
	result := make([]int, n)
	for i, p := range rand.Perm(n) {
		result[i] = seq[p]
	}
	return result
}
