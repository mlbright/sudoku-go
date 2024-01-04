package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/mlbright/sudoku-norvig-go/sudoku"
)

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
	puzzle := sudoku.New()
	start := time.Now().UnixNano()
	solution, _ := puzzle.Solve(grid)
	end := time.Now().UnixNano()
	duration := end - start
	// fmt.Println(grid)
	// puzzle.ShowSolution(solution)
	return duration, puzzle.Solved(solution)
}

func fromFile(filename string) []string {
	dat, _ := ioutil.ReadFile(filename)
	grids := strings.Split(string(dat), "\n")
	return grids[:len(grids)-1]
}

func nanoconv(nanos int64) float64 {
	return float64(nanos) / 1000000000.0
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

func solveAll(grids []string, name string) {
	times := make([]int64, 0)
	results := make([]bool, 0)

	for _, grid := range grids {
		t, result := timeSolve(grid)
		times = append(times, t)
		results = append(results, result)
	}

	n := len(grids)
	if n >= 1 {
		fmt.Printf("Solved %d of %d %s puzzles (avg %.4f secs (%.2f Hz), max %.4f secs).\n",
			sum(bool2int(results)), n, name, float64(nanoconv(sum(times)))/float64(n), float64(n)/float64(nanoconv(sum(times))), nanoconv(max(times)))
	}
}

func main() {
	solveAll(fromFile("puzzles/incredibly-difficult.txt"), "incredibly-difficult")
	solveAll(fromFile("puzzles/easy50.txt"), "easy")
	solveAll(fromFile("puzzles/top95.txt"), "hard")
	solveAll(fromFile("puzzles/hardest.txt"), "hardest")
	solveAll(fromFile("puzzles/hardest20.txt"), "hardest20")
	solveAll(fromFile("puzzles/hardest20x50.txt"), "hardest20x50")
	solveAll(fromFile("puzzles/topn87.txt"), "topn87")
	solveAll(fromFile("puzzles/all.txt"), "all")
}
