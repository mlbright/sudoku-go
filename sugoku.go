package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type iunit []string
type unit []int
type unitgroup []unit
type peerlist []int

var (
	rows        string
	cols        string
	digits      string
	squares     []string
	unitlist    []unit
	iunitlist   []iunit
	units       []unitgroup
	peers       []peerlist
	blankPuzzle []string
	valid       string
)

const puzzleN int = 17

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

func test() {

	if len(squares) != 81 {
		panic("the number of squares is not 81")
	}

	if len(iunitlist) != 27 {
		panic("the number of iunits is not 27")
	}

	if len(unitlist) != 27 {
		panic("the number of units is not 27")
	}

	for s := range squares {
		if len(units[s]) != 3 {
			panic("bad unit")
		}
	}

	for s := range squares {
		if len(peers[s]) != 20 {
			panic("bad peer list")
		}
	}

	fmt.Println("All tests pass.")
}

// Parse a grid

func parseGrid(grid string) ([]string, bool) {
	puzzle := gridValues(grid)
	// To start, every square can be any digit; then assign values from the grid.
	solution := make([]string, len(puzzle))
	copy(solution, blankPuzzle)
	for s, d := range puzzle {
		if strings.Contains(digits, d) {
			if !assign(solution, s, d) {
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
		if strings.Contains(valid, c) {
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

func assign(puzzle []string, s int, d string) bool {
	otherValues := strings.Replace(puzzle[s], d, "", -1)
	for _, otherValue := range strings.Split(otherValues, "") {
		if !eliminate(puzzle, s, otherValue) {
			return false
		}
	}
	return true
}

func eliminate(puzzle []string, s int, valueToEliminate string) bool {
	if !strings.Contains(puzzle[s], valueToEliminate) {
		return true // Already eliminated
	}

	puzzle[s] = strings.Replace(puzzle[s], valueToEliminate, "", -1)

	if len(puzzle[s]) == 0 {
		return false // Contradiction, removed last value
	}

	// (1) If a square s is reduced to one value, then eliminate it from the peers.
	if len(puzzle[s]) == 1 {
		d2 := puzzle[s]
		for _, peer := range peers[s] {
			if !eliminate(puzzle, peer, d2) {
				return false
			}
		}
	}

	// (2) If a unit u is reduced to only one place for a value d, then put it there.
	for _, u := range units[s] {
		dplaces := make([]int, 0)
		for _, sq := range u {
			if strings.Contains(puzzle[sq], valueToEliminate) {
				dplaces = append(dplaces, sq)
			}
		}

		numSpots := len(dplaces)
		if numSpots == 0 {
			return false // Contradiction: no place for this value
		}

		if numSpots == 1 {
			if !assign(puzzle, dplaces[0], valueToEliminate) {
				return false
			}
		}
	}
	return true
}

// Search

func solve(grid string) ([]string, bool) {
	puzzle, ok := parseGrid(grid)
	if ok {
		return search(puzzle)
	}
	return puzzle, ok
}

func search(puzzle []string) ([]string, bool) {
	minSquare := 82
	minSize := 10

	for s := range squares {
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

		if assign(puzzleCopy, minSquare, d) {
			result, ok := search(puzzleCopy)
			if ok {
				return result, true
			}
		}
	}

	return puzzle, false
}

func max(values []int64) int64 {
	var max int64
	max = 0
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

func sum(items []int64) int64 {
	var accum int64
	accum = 0
	for _, b := range items {
		accum += b
	}
	return accum
}

func bool2int(booleans []bool) []int64 {
	ints := make([]int64, 0)
	for _, b := range booleans {
		if b {
			ints = append(ints, 1)
		} else {
			ints = append(ints, 0)
		}
	}
	return ints
}

func timeSolve(grid string) (int64, bool) {
	start := time.Now()
	nanosStart := start.UnixNano()
	puzzle, _ := solve(grid)
	end := time.Now()
	nanosEnd := end.UnixNano()
	duration := nanosEnd - nanosStart
	/*
	   fmt.Println(grid)
	   var solved_ []string
	   for _,sq := range squares {
	       solved_ = append(solved_,puzzle[sq])
	   }
	   fmt.Println(strings.Join(solved_,""))
	*/
	return duration, solved(puzzle)
}

func fromFile(filename string) []string {
	dat, _ := ioutil.ReadFile(filename)
	grids := strings.Split(string(dat), "\n")
	return grids[:len(grids)-1]
}

func nanoconv(nanos int64) float64 {
	return float64(nanos) / 1000000000.0
}

func solveAll(grids []string, name string) {
	times := make([]int64, 0)
	results := make([]bool, 0)

	for _, grid := range grids {
		t, result := timeSolve(grid)
		times = append(times, t)
		results = append(results, result)
	}

	n := len(grids)
	if n > 1 {
		fmt.Printf("Solved %d of %d %s puzzles (avg %.4f secs (%.2f Hz), max %.4f secs).\n",
			sum(bool2int(results)), n, name, float64(nanoconv(sum(times)))/float64(n), float64(n)/float64(nanoconv(sum(times))), nanoconv(max(times)))
	}
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

func solved(puzzle []string) bool {
	for _, u := range unitlist {
		if !unitSolved(puzzle, u) {
			return false
		}
	}
	return true
}

func randomPuzzle() string {
	puzzle := make([]string, 81)
	copy(puzzle, blankPuzzle)
	shuffled := make([]string, len(squares))
	perm := rand.Perm(len(squares))
	for i, v := range perm {
		shuffled[v] = squares[i]
	}
	for s := range shuffled {
		elements := strings.Split(puzzle[s], "")
		if !assign(puzzle, s, elements[rand.Intn(len(elements))]) {
			break
		}
		ds := make([]string, 0)
		for sq := range squares {
			if len(puzzle[sq]) == 1 {
				ds = append(ds, puzzle[sq])
			}
		}
		set := make(map[string]bool)
		for _, sq := range ds {
			set[sq] = true
		}
		if len(ds) >= puzzleN && len(set) >= 8 {
			out := make([]string, 0)
			for sq := range squares {
				if len(puzzle[sq]) == 1 {
					out = append(out, puzzle[sq])
				} else {
					out = append(out, ".")
				}
			}
			puzzle := strings.Join(out, "")
			return puzzle
		}
	}
	return randomPuzzle()
}

func main() {
	rows = "ABCDEFGHI"
	digits = "123456789"
	valid = digits + "."
	cols = digits
	squares = cross(rows, cols)
	squaresDict := make(map[string]int)
	for i, sq := range squares {
		squaresDict[sq] = i
	}

	for range squares {
		blankPuzzle = append(blankPuzzle, digits)
	}

	iunitlist = make([]iunit, 0)

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

	unitlist = make([]unit, 0)
	for _, u := range iunitlist {
		squareList := make(unit, 0)
		for _, sq := range u {
			squareList = append(squareList, squaresDict[sq])
		}
		unitlist = append(unitlist, squareList)
	}

	units = make([]unitgroup, 0)
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

	peers = make([]peerlist, 81)

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

	test()
	solveAll(fromFile("easy50.txt"), "easy")
	solveAll(fromFile("top95.txt"), "hard")
	solveAll(fromFile("hardest.txt"), "hardest")
	randomPuzzles := make([]string, 99)
	for j := 0; j < 99; j++ {
		randomPuzzles[j] = randomPuzzle()
	}
	solveAll(randomPuzzles, "random")
}
