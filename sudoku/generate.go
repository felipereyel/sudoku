package sudoku

import "math/bits"

func collectValues(mask uint16) []int {
	vals := make([]int, 0, bits.OnesCount16(mask))
	for mask != 0 {
		lsb := mask & -mask
		vals = append(vals, bits.TrailingZeros16(lsb)+1)
		mask &^= lsb
	}
	return vals
}

func assign(c *Candidates, square, val int) bool {
	bit := uint16(1 << (val - 1))
	other := c[square] & ^bit
	for other != 0 {
		lsb := other & -other
		oVal := bits.TrailingZeros16(lsb) + 1
		if !eliminate(c, square, oVal) {
			return false
		}
		other &^= lsb
	}
	return true
}

func eliminate(c *Candidates, square, val int) bool {
	bit := uint16(1 << (val - 1))
	if c[square]&bit == 0 {
		return true
	}
	c[square] &^= bit
	if c[square] == 0 {
		return false
	}

	if c[square]&(c[square]-1) == 0 {
		onlyVal := bits.TrailingZeros16(c[square]) + 1
		for _, p := range peers[square] {
			if !eliminate(c, p, onlyVal) {
				return false
			}
		}
	}

	for _, ui := range unitsForSquare[square] {
		place := -1
		multiple := false
		for _, cell := range units[ui] {
			if c[cell]&bit != 0 {
				if place == -1 {
					place = cell
				} else {
					multiple = true
					break
				}
			}
		}
		if place == -1 {
			return false
		}
		if !multiple {
			if !assign(c, place, val) {
				return false
			}
		}
	}

	return true
}

func search(c Candidates, reverse bool) *Candidates {
	minN := 10
	minIdx := -1
	for i := 0; i < N2; i++ {
		n := bits.OnesCount16(c[i])
		if n > 1 && n < minN {
			minN = n
			minIdx = i
		}
	}
	if minIdx == -1 {
		return &c
	}

	vals := collectValues(c[minIdx])

	if reverse {
		for i := len(vals) - 1; i >= 0; i-- {
			copyC := c
			if assign(&copyC, minIdx, vals[i]) {
				if result := search(copyC, true); result != nil {
					return result
				}
			}
		}
	} else {
		for _, v := range vals {
			copyC := c
			if assign(&copyC, minIdx, v) {
				if result := search(copyC, false); result != nil {
					return result
				}
			}
		}
	}

	return nil
}

func boardToCandidates(b Board) Candidates {
	c := newCandidates()
	for i, v := range b {
		if v != 0 {
			if !assign(&c, i, v) {
				return Candidates{}
			}
		}
	}
	return c
}

func buildAndVerify(c Candidates, n int, givens []int) (Board, bool) {
	var board Board
	for _, gi := range givens {
		board[gi] = bits.TrailingZeros16(c[gi]) + 1
	}

	if len(givens) > n {
		shuffled := shuffle(givens)
		for i := 0; i < len(givens)-n; i++ {
			board[shuffled[i]] = 0
		}
	}

	fwdC := boardToCandidates(board)
	if fwdC == (Candidates{}) {
		return Board{}, false
	}
	fwdSol := search(fwdC, false)
	if fwdSol == nil {
		return Board{}, false
	}

	revC := boardToCandidates(board)
	revSol := search(revC, true)
	if revSol == nil {
		return Board{}, false
	}

	if *fwdSol != *revSol {
		return Board{}, false
	}

	return board, true
}

func tryGenerate(n int) (Board, bool) {
	c := newCandidates()
	cells := shuffle(squares)

	for _, cell := range cells {
		nC := bits.OnesCount16(c[cell])
		if nC == 0 {
			return Board{}, false
		}

		vals := collectValues(c[cell])
		val := vals[randRange(0, len(vals))]

		if !assign(&c, cell, val) {
			return Board{}, false
		}

		var givens []int
		for i := 0; i < N2; i++ {
			if bits.OnesCount16(c[i]) == 1 {
				givens = append(givens, i)
			}
		}

		if len(givens) >= n {
			uniq := make(map[int]bool)
			for _, gi := range givens {
				v := bits.TrailingZeros16(c[gi]) + 1
				uniq[v] = true
			}
			if len(uniq) >= 8 {
				return buildAndVerify(c, n, givens)
			}
		}
	}

	return Board{}, false
}

// Generate generates a sudoku puzzle with approximately n givens.
// n is clamped to [17, 81].
func Generate(n int) Board {
	if n < MinGivens {
		n = MinGivens
	}
	if n > MaxGivens {
		n = MaxGivens
	}
	for {
		board, ok := tryGenerate(n)
		if ok {
			return board
		}
	}
}
