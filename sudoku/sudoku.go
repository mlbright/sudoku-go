package sudoku

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dropbox/godropbox/container/bitvector"
)

const (
	rows   = "ABCDEFGHI"
	digits = "123456789"
)

type puzzle []bitvector.BitVector

func newPuzzleElement() *bitvector.BitVector {
	return bitvector.NewBitVector([]byte{byte(0b11111111), byte(0b11111111)}, 9)
}

func (p puzzle) lengthAndRemainingValue(i int) (int, int) {
	n := 0
	lastValueSet := 10
	for j := 0; j <= 8; j++ {
		if p.isSet(i, j) {
			n++
			lastValueSet = j
		}
	}
	return n, lastValueSet
}

func (p puzzle) length(i int) int {
	n := 0
	for j := 0; j <= 8; j++ {
		if p.isSet(i, j) {
			n++
		}
	}
	return n
}

func (p puzzle) isSet(i, j int) bool {
	return p[i].Element(j) == 1
}

func (p puzzle) unset(i, j int) {
	p[i].Set(0, j)
}

func (p puzzle) Duplicate() puzzle {
	tmp := make([]bitvector.BitVector, 81)
	for i := 0; i < 81; i++ {
		// for j := 0; j < p[i].Length(); j++ {
		// 	tmp[i].Append(p[i].Element(j))
		// }
		backing := make([]byte, 2)
		copy(backing, p[i].Bytes())
		tmp[i] = *bitvector.NewBitVector(backing, p[i].Length())
	}
	return tmp
}

type iunit []string
type unit []int
type unitgroup []unit
type peerlist []int

type sudoku struct {
	unitlist []unit
	units    []unitgroup
	peers    []peerlist
}

func cross(x string, y string) []string {
	result := make([]string, 0)
	a := strings.Split(x, "")
	b := strings.Split(y, "")
	for _, i := range a {
		for _, j := range b {
			s := []string{i, j}
			result = append(result, strings.Join(s, ""))
		}
	}
	return result
}

func (s *sudoku) parseGrid(grid string) (puzzle, bool) {
	// To start, every square can be any digit; then assign values from the grid.
	solution := s.BlankPuzzle()
	for sq, d := range gridValues(grid) {
		if strings.Contains(digits, d) {
			n, err := strconv.Atoi(d)
			if err != nil {
				panic(err)
			}
			if !s.assign(solution, sq, n-1) {
				return solution, false
			}
		}
	}
	return solution, true
}

func gridValues(grid string) []string {
	p := make([]string, 81)
	i := 0
	for _, c := range strings.Split(grid, "") {
		if strings.Contains(digits, c) || strings.Contains(".", c) {
			p[i] = c
			i++
		}
	}
	if len(p) != 81 {
		panic("invalid puzzle")
	}
	return p
}

// Constraint Propagation

func (s *sudoku) assign(p puzzle, sq int, valueToAssign int) bool {
	for j := 0; j <= 8; j++ {
		if p.isSet(sq, j) && j != valueToAssign {
			if !s.eliminate(p, sq, j) {
				return false
			}
		}
	}
	return true
}

func (s *sudoku) eliminate(p puzzle, sq int, valueToEliminate int) bool {
	if !p.isSet(sq, valueToEliminate) {
		return true // already eliminated
	}

	// (A)
	p.unset(sq, valueToEliminate)

	numberOfRemainingValues, remainingValue := p.lengthAndRemainingValue(sq)

	if numberOfRemainingValues == 0 {
		return false // Contradiction: removed last value
	} else if numberOfRemainingValues == 1 {
		// (1) If the square sq is reduced to one value, then eliminate the value from its peers.
		for _, peer := range s.peers[sq] {
			if !s.eliminate(p, peer, remainingValue) {
				return false
			}
		}
	}

	// (2) After (A), if a unit u has only one spot left to place valueToEliminate, then assign it there.
CheckUnits:
	for _, u := range s.units[sq] {
		remainingSquareForValueToEliminate := 82
		numberOfPossibleSquaresForValueToEliminate := 0

		for _, sq := range u {
			if p.isSet(sq, valueToEliminate) {
				remainingSquareForValueToEliminate = sq
				numberOfPossibleSquaresForValueToEliminate++
			}

			if numberOfPossibleSquaresForValueToEliminate > 1 {
				continue CheckUnits
			}
		}

		if numberOfPossibleSquaresForValueToEliminate == 0 {
			return false // Contradiction: no valid square for valueToEliminate
		} else if numberOfPossibleSquaresForValueToEliminate == 1 {
			if !s.assign(p, remainingSquareForValueToEliminate, valueToEliminate) {
				return false
			}
		}
	}

	return true
}

func (s *sudoku) Solve(grid string) (puzzle, bool) {
	p, ok := s.parseGrid(grid)
	if ok {
		return s.search(p)
	}
	return p, false
}

func (s *sudoku) search(p puzzle) (puzzle, bool) {
	squareWithFewestPossibilities := 82
	minSize := 10

	for sq := 0; sq < 81; sq++ {
		l := p.length(sq)

		if l > 1 && l < minSize {
			minSize = l
			squareWithFewestPossibilities = sq
			if minSize == 2 {
				break
			}
		}
	}

	if squareWithFewestPossibilities == 82 {
		return p, true // solved
	}

	for j := 0; j <= 8; j++ {
		if p.isSet(squareWithFewestPossibilities, j) {
			copied := p.Duplicate()
			if s.assign(copied, squareWithFewestPossibilities, j) {
				solution, ok := s.search(copied)
				if ok {
					return solution, true
				}
			}
		}
	}

	return p, false
}

func unitSolved(p puzzle, u unit) bool {
	set := make(map[int]bool)
	for _, sq := range u {
		value := 0
		for j := 0; j <= 8; j++ {
			if p.isSet(sq, j) {
				value = j
				break
			}
		}
		set[value] = true
	}
	for d := 0; d <= 8; d++ {
		if !set[d] {
			return false
		}
	}
	return true
}

func (s *sudoku) Solved(p puzzle) bool {
	for _, u := range s.unitlist {
		if !unitSolved(p, u) {
			return false
		}
	}
	return true
}

func (p *sudoku) BlankPuzzle() puzzle {
	puzzle := make([]bitvector.BitVector, 81)
	for i := range puzzle {
		puzzle[i] = *newPuzzleElement()
	}
	return puzzle
}

func New() *sudoku {
	solver := sudoku{}

	cols := digits
	squares := cross(rows, cols)

	squaresDict := make(map[string]int)
	for i, sq := range squares {
		squaresDict[sq] = i
	}

	iunitlist := make([]iunit, 0)

	for _, c := range cols {
		iunitlist = append(iunitlist, cross(rows, string(c)))
	}

	for _, r := range rows {
		iunitlist = append(iunitlist, cross(string(r), cols))
	}

	rs := []string{"ABC", "DEF", "GHI"}
	cs := []string{"123", "456", "789"}

	for _, r := range rs {
		for _, c := range cs {
			iunitlist = append(iunitlist, cross(r, c))
		}
	}

	unitlist := make([]unit, 0)
	for _, u := range iunitlist {
		squareList := make(unit, 0)
		for _, sq := range u {
			squareList = append(squareList, squaresDict[sq])
		}
		unitlist = append(unitlist, squareList)
	}

	units := make([]unitgroup, 0)
	for s := range squares {
		group := make(unitgroup, 0)
		for _, unit := range unitlist {
			for _, square := range unit {
				if square == s {
					group = append(group, unit)
					break
				}
			}
		}
		units = append(units, group)
	}

	peers := make([]peerlist, 81)

	for s := range squares {
		peerSet := make(map[int]bool)
		for _, unit := range units[s] {
			for _, square := range unit {
				if square != s {
					peerSet[square] = true
				}
			}
		}
		peerList := make([]int, 0)
		for k := range peerSet {
			peerList = append(peerList, k)
		}
		peers[s] = peerList
	}

	solver.peers = peers
	solver.units = units
	solver.unitlist = unitlist
	return &solver
}

func (s *sudoku) ShowSolution(p puzzle) {
	var b strings.Builder
	for i := range p {
		l := p.length(i)
		if l != 1 {
			b.WriteString("[")
		}
		for j := 0; j <= 8; j++ {
			if p.isSet(i, j) {
				b.WriteString(strconv.Itoa(j + 1))
			}
		}
		if l != 1 {
			b.WriteString("]")
		}
	}
	fmt.Println(b.String())
}
