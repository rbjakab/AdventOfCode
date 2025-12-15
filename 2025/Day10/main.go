package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Puzzle struct {
	A     [][]bool
	b     []bool
	jolts []int
}

func createPuzzle(parts []string) Puzzle {
	var newPuzzle Puzzle

	// [...]
	innerBrackets := parts[0][1 : len(parts[0])-1]
	for i := 0; i < len(innerBrackets); i++ {
		newPuzzle.b = append(newPuzzle.b, innerBrackets[i] == '#')
	}

	// (...)
	rows := len(newPuzzle.b)
	cols := len(parts) - 2
	var A [][]bool
	for range rows {
		var row []bool
		for range cols {
			row = append(row, false)
		}
		A = append(A, row)
	}
	for j := 1; j < len(parts)-1; j++ {
		innerBrackets := parts[j][1 : len(parts[j])-1]
		numbers := strings.Split(innerBrackets, ",")

		for _, n := range numbers {
			i, _ := strconv.Atoi(n)
			A[i][j-1] = true
		}
	}
	newPuzzle.A = A

	// {...}
	innerBrackets = parts[len(parts)-1][1 : len(parts[len(parts)-1])-1]
	jolts := strings.Split(innerBrackets, ",")
	for _, j := range jolts {
		i, _ := strconv.Atoi(j)
		newPuzzle.jolts = append(newPuzzle.jolts, i)
	}

	return newPuzzle
}

func readInput(filename string) []Puzzle {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var puzzles []Puzzle

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		newPuzzle := createPuzzle(parts)
		puzzles = append(puzzles, newPuzzle)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return puzzles
}

func printMatrix(A [][]float64) {
	const eps = 1e-9

	for i := range A {
		for j := 0; j < len(A[0])-1; j++ {
			v := A[i][j]
			if math.Abs(v) < eps {
				v = 0
			}
			fmt.Printf("%4.1f ", v)
		}

		rhs := A[i][len(A[0])-1]
		if math.Abs(rhs) < eps {
			rhs = 0
		}
		fmt.Printf("| %4.1f\n", rhs)
	}
}

func increase(x []bool) {
	i := len(x) - 1
	propagation := true

	for propagation && i >= 0 {
		propagation = x[i]
		x[i] = !x[i]
		i--
	}
}

func check(A [][]bool, b []bool, x []bool) bool {
	rows := len(A)
	cols := len(A[0])

	for i := range rows {
		trueCounts := 0
		for j := range cols {
			if A[i][j] && x[j] {
				trueCounts++
			}
		}
		if (trueCounts%2 == 1) != b[i] {
			return false
		}
	}

	return true
}

func countTrues(a []bool) int {
	n := 0
	for _, i := range a {
		if i {
			n++
		}
	}
	return n
}

func solve(A [][]bool, b []bool) int {
	x := make([]bool, len(A[0]))
	for i := range x {
		x[i] = false
	}

	minCount := math.MaxInt
	maxCombinations := int(math.Pow(2, float64(len(x))))

	for range maxCombinations {
		if check(A, b, x) {
			count := countTrues(x)
			if count < minCount {
				minCount = count
			}
		}
		increase(x)
	}

	return minCount
}

type Var struct {
	pos int
	val float64
}

func buildAugmentedMatrix(A [][]bool, b []int) [][]float64 {
	m := len(A)    // number of rows
	n := len(A[0]) // number of columns (variables)

	if len(b) != m {
		panic("row count of A does not match length of b")
	}

	aug := make([][]float64, m)

	for i := 0; i < m; i++ {
		row := make([]float64, n+1)

		for j := 0; j < n; j++ {
			if A[i][j] {
				row[j] = 1
			} else {
				row[j] = 0
			}
		}

		row[n] = float64(b[i]) // RHS
		aug[i] = row
	}

	return aug
}

func getFreeVariables(A [][]float64) []Var {
	var freeVars []Var

	m := len(A)
	n := len(A[0]) - 1 // exclude RHS

	pivotCols := make(map[int]bool)

	for row := 0; row < m; row++ {
		for col := 0; col < n; col++ {
			if A[row][col] != 0 {
				pivotCols[col] = true
				break
			}
		}
	}

	for col := 0; col < n; col++ {
		if _, ok := pivotCols[col]; !ok {
			freeVars = append(freeVars, Var{pos: col, val: 0})
		}
	}

	return freeVars
}

func extractSolution(A [][]float64, free []Var) []float64 {
	n := len(A[0]) - 1
	x := make([]float64, n)

	for _, v := range free {
		x[v.pos] = v.val
	}

	for i := 0; i < len(A); i++ {
		p := -1
		for j := 0; j < n; j++ {
			if math.Abs(A[i][j]-1) < eps {
				p = j
				break
			}
		}
		if p == -1 {
			continue
		}
		val := A[i][n]
		for j := 0; j < n; j++ {
			if j != p {
				val -= A[i][j] * x[j]
			}
		}
		x[p] = val
	}
	return x
}

func search(A [][]float64, level int, free []Var) {
	if level == len(free) {
		x := extractSolution(A, free)
		sum := 0
		for _, v := range x {
			iv := math.Round(v)

			if math.Abs(v-iv) > eps {
				return // not integer
			}

			if iv < 0 {
				return // truly negative
			}

			sum += int(iv)
		}

		if sum < bestSum {
			bestSum = sum
			bestSolution = append([]float64{}, x...)
		}
		return
	}

	for v := 0; v <= MAX; v++ {
		next := make([]Var, len(free))
		copy(next, free)
		next[level].val = float64(v)
		search(A, level+1, next)
	}
}

func getHighestJolt(jolts []int) int {
	highest := jolts[0]
	for _, j := range jolts {
		absJ := int(math.Abs(float64(j)))
		if absJ > highest {
			highest = absJ
		}
	}
	return highest
}

var MAX int
var bestSum int
var bestSolution []float64

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	inputFile := flag.String("file", "input.txt", "path to the input file")
	flag.Parse()

	start := time.Now()
	log.Println("Starting...")

	// === Part 1 ===
	minCounts := 0
	puzzles := readInput(*inputFile)
	for _, puzzle := range puzzles {
		minCount := solve(puzzle.A, puzzle.b)
		minCounts += minCount
	}
	log.Println("PART 1: Total:", minCounts)

	// === Part 2 ===
	log.Println("PART 2...")
	total := 0

	for i, puzzle := range puzzles {
		MAX = getHighestJolt(puzzle.jolts)
		bestSum = math.MaxInt
		bestSolution = make([]float64, 0)

		log.Printf("%d / %d", i+1, len(puzzles))
		log.Println("MAX:", MAX)

		augmented := buildAugmentedMatrix(puzzle.A, puzzle.jolts)
		printMatrix(augmented)
		fmt.Println("---")
		gA := rref(augmented)
		printMatrix(gA)
		freeVariables := getFreeVariables(gA)
		log.Printf("Free variables: %v", freeVariables)
		free := getFreeVariables(gA)
		search(gA, 0, free)
		minimum := int(bestSum)
		log.Printf("The minimum: %d\n", minimum)
		log.Println("Solution:", bestSolution)

		if bestSum > 5000 {
			log.Println("=========== NOOOOOOOOOOOO ===========")
			break
		}

		total += int(bestSum)
	}

	log.Println("In total:", total)
	log.Printf("Execution time: %s", time.Since(start))
}
