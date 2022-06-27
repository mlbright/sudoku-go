package sudoku

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	rows       = "ABCDEFGHI"
	digits     = "123456789"
	N      int = 17
)

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

func (p *sudoku) Test() {

	if len(p.squares) != 81 {
		panic("the number of squares is not 81")
	}

	if len(p.unitlist) != 27 {
		panic("the number of units is not 27")
	}

	for s := range p.squares {
		if len(p.units[s]) != 3 {
			panic("bad unit")
		}
	}

	for s := range p.squares {
		if len(p.peers[s]) != 20 {
			panic("bad peer list")
		}
	}

	fmt.Println("All tests pass.")
}

func (p *sudoku) parseGrid(grid string) ([]string, bool) {
	// To start, every square can be any digit; then assign values from the grid.
	solution := p.blankPuzzle()
	for s, d := range gridValues(grid) {
		if strings.Contains(digits, d) {
			if !p.assign(solution, s, d) {
				return solution, false
			}
		}
	}
	return solution, true
}

func gridValues(grid string) []string {
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

func (p *sudoku) assign(puzzle []string, s int, d string) bool {
	otherValues := strings.Replace(puzzle[s], d, "", -1)
	for _, otherValue := range strings.Split(otherValues, "") {
		if !p.eliminate(puzzle, s, otherValue) {
			return false
		}
	}
	return true
}

func (p *sudoku) eliminate(puzzle []string, s int, valueToEliminate string) bool {
	if !strings.Contains(puzzle[s], valueToEliminate) {
		return true // Already eliminated
	}

	// (A)
	puzzle[s] = strings.Replace(puzzle[s], valueToEliminate, "", -1)

	if len(puzzle[s]) == 0 {
		return false // Contradiction, removed last value
	}

	// (1) If a square s is reduced to one value, then eliminate it from the peers.
	if len(puzzle[s]) == 1 {
		lastRemainingValue := puzzle[s]
		for _, peer := range p.peers[s] {
			if !p.eliminate(puzzle, peer, lastRemainingValue) {
				return false
			}
		}
	}

	// (2) After (A), if a unit u has only one spot left to place valueToEliminate, then assign it there.
CheckUnits:
	for _, u := range p.units[s] {
		remainingSquareForValueToEliminate := 82
		numberOfPossibleSquaresForValueToEliminate := 0

		for _, sq := range u {
			if strings.Contains(puzzle[sq], valueToEliminate) {
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
			if !p.assign(puzzle, remainingSquareForValueToEliminate, valueToEliminate) {
				return false
			}
		}
	}
	return true
}

func (p *sudoku) Solve(grid string) ([]string, bool) {
	puzzle, ok := p.parseGrid(grid)
	if ok {
		return p.search(puzzle)
	}
	return puzzle, false
}

func (p *sudoku) search(puzzle []string) ([]string, bool) {
	minSquare := 82
	minSize := 10

	for s := range p.squares {
		size := len(puzzle[s])
		if size > 1 && size < minSize {
			minSquare = s
			minSize = size
		}
	}

	if minSquare == 82 {
		return puzzle, true
	}

	for _, d := range strings.Split(puzzle[minSquare], "") {
		puzzleCopy := make([]string, 81)
		copy(puzzleCopy, puzzle)

		if p.assign(puzzleCopy, minSquare, d) {
			result, ok := p.search(puzzleCopy)
			if ok {
				return result, true
			}
		}
	}

	return puzzle, false
}

func unitSolved(puzzle []string, u unit) bool {
	set := make(map[string]bool)
	for _, s := range u {
		set[puzzle[s]] = true
	}
	for _, d := range strings.Split(digits, "") {
		if !set[d] {
			return false
		}
	}
	return true
}

func (p *sudoku) Solved(puzzle []string) bool {
	for _, u := range p.unitlist {
		if !unitSolved(puzzle, u) {
			return false
		}
	}
	return true
}

func (p *sudoku) Random() string {
	solution := p.blankPuzzle()
	shuffled := make([]string, 81)
	perm := rand.Perm(81)
	for i, v := range perm {
		shuffled[v] = p.squares[i]
	}
	for s := range shuffled {
		elements := strings.Split(solution[s], "")
		if !p.assign(solution, s, elements[rand.Intn(len(elements))]) {
			break
		}
		ds := make([]string, 0)
		for sq := range p.squares {
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
			for sq := range p.squares {
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
	return p.Random()
}

func (p *sudoku) blankPuzzle() []string {
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

func (p *sudoku) ShowSolution(s []string) {
	fmt.Println(strings.Join(s, ""))
}
