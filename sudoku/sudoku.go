package sudoku

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	rows   = "ABCDEFGHI"
	digits = "123456789"
	N      = 17
)

type puzzle []string
type iunit []string
type unit []int
type unitgroup []unit
type peerlist []int

type sudoku struct {
	squares  []string
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

func (s *sudoku) Test() {

	if len(s.squares) != 81 {
		panic("the number of squares is not 81")
	}

	if len(s.unitlist) != 27 {
		panic("the number of units is not 27")
	}

	for sq := range s.squares {
		if len(s.units[sq]) != 3 {
			panic("bad unit")
		}
	}

	for sq := range s.squares {
		if len(s.peers[sq]) != 20 {
			panic("bad peer list")
		}
	}

	fmt.Println("All tests pass.")
}

func (s *sudoku) parseGrid(grid string) (puzzle, bool) {
	// To start, every square can be any digit; then assign values from the grid.
	solution := s.blankPuzzle()
	for sq, d := range gridValues(grid) {
		if strings.Contains(digits, d) {
			if !s.assign(solution, sq, d) {
				return solution, false
			}
		}
	}
	return solution, true
}

func gridValues(grid string) puzzle {
	puzzle := make([]string, 81)
	i := 0
	for _, c := range strings.Split(grid, "") {
		if strings.Contains(digits, c) || strings.Contains("0.", c) {
			puzzle[i] = c
			i++
		}
	}
	if len(puzzle) != 81 {
		panic("invalid puzzle")
	}
	return puzzle
}

// Constraint Propagation

func (s *sudoku) assign(p puzzle, sq int, d string) bool {
	otherValues := strings.Replace(p[sq], d, "", -1)
	for _, otherValue := range strings.Split(otherValues, "") {
		if !s.eliminate(p, sq, otherValue) {
			return false
		}
	}
	return true
}

func (s *sudoku) eliminate(p puzzle, sq int, valueToEliminate string) bool {
	if !strings.Contains(p[sq], valueToEliminate) {
		return true // Already eliminated
	}

	// (A)
	p[sq] = strings.Replace(p[sq], valueToEliminate, "", -1)

	if len(p[sq]) == 0 {
		return false // Contradiction, removed last value
	}

	// (1) If a square s is reduced to one value, then eliminate it from the peers.
	if len(p[sq]) == 1 {
		lastRemainingValue := p[sq]
		for _, peer := range s.peers[sq] {
			if !s.eliminate(p, peer, lastRemainingValue) {
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
			if strings.Contains(p[sq], valueToEliminate) {
				remainingSquareForValueToEliminate = sq
				numberOfPossibleSquaresForValueToEliminate++
			}

			if numberOfPossibleSquaresForValueToEliminate > 1 {
				continue CheckUnits
			}
		}

		if numberOfPossibleSquaresForValueToEliminate == 0 {
			return false // Contradiction: no valid square for valueToEliminate
		}

		if numberOfPossibleSquaresForValueToEliminate == 1 {
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
	minSquare := 82
	minSize := 10

	for sq := range s.squares {
		size := len(p[sq])
		if size > 1 && size < minSize {
			minSquare = sq
			minSize = size
			if minSize == 2 {
				break
			}
		}
	}

	if minSquare == 82 {
		return p, true
	}

	for _, d := range strings.Split(p[minSquare], "") {
		puzzleCopy := make([]string, 81)
		copy(puzzleCopy, p)

		if s.assign(puzzleCopy, minSquare, d) {
			result, ok := s.search(puzzleCopy)
			if ok {
				return result, true
			}
		}
	}

	return p, false
}

func unitSolved(p puzzle, u unit) bool {
	set := make(map[string]bool)
	for _, sq := range u {
		set[p[sq]] = true
	}
	for _, d := range strings.Split(digits, "") {
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

func (s *sudoku) Random() string {
	solution := s.blankPuzzle()
	shuffled := make([]string, 81)
	perm := rand.Perm(81)
	for i, v := range perm {
		shuffled[v] = s.squares[i]
	}
	for sq := range shuffled {
		elements := strings.Split(solution[sq], "")
		if !s.assign(solution, sq, elements[rand.Intn(len(elements))]) {
			break
		}
		ds := make([]string, 0)
		for sq := range s.squares {
			if len(solution[sq]) == 1 {
				ds = append(ds, solution[sq])
			}
		}
		set := make(map[string]bool)
		for _, sq := range ds {
			set[sq] = true
		}
		if len(ds) >= N && len(set) >= 8 {
			out := make([]string, 0)
			for sq := range s.squares {
				if len(solution[sq]) == 1 {
					out = append(out, solution[sq])
				} else {
					out = append(out, ".")
				}
			}
			puzzle := strings.Join(out, "")
			return puzzle
		}
	}
	return s.Random()
}

func (p *sudoku) blankPuzzle() puzzle {
	solution := make([]string, 81)
	for i := range solution {
		solution[i] = digits
	}
	return solution
}

func New() *sudoku {
	solver := sudoku{}
	solver.blankPuzzle()

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
	solver.squares = squares
	solver.units = units
	solver.unitlist = unitlist
	return &solver
}

func (s *sudoku) ShowSolution(p puzzle) {
	fmt.Println(strings.Join(p, ""))
}
