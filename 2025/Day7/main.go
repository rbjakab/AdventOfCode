package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"
)

type Pos struct {
	row int
	col int
}

type PosMap map[Pos]int

func AllIndexesASCII(s string, target []byte) []int {
	var idxs []int
	for i := 0; i < len(s); i++ {
		if slices.Contains(target, s[i]) {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func readAndSolveInput(filename string) int {
	file, _ := os.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	var newLine []byte
	splits := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(lines) == 0 {
			lines = append(lines, line)
			continue
		}

		prevLine := lines[len(lines)-1]
		newLine = []byte(line)

		allSplittersIndexes := AllIndexesASCII(line, []byte{'^'})
		for _, i := range allSplittersIndexes {
			if prevLine[i] == '|' || prevLine[i] == 'S' {
				splits++
				newLine[i-1] = '|'
				newLine[i+1] = '|'
			}
		}

		allBeamIndexes := AllIndexesASCII(prevLine, []byte{'|', 'S'})
		for _, i := range allBeamIndexes {
			if (prevLine[i] == '|' || prevLine[i] == 'S') && newLine[i] != '^' {
				newLine[i] = '|'
			}
		}

		lines = append(lines, string(newLine))
	}

	SaveLines(strings.Split(strings.Split(filename, "/")[1], ".")[0]+"_filled.txt", lines)

	return splits
}

func SaveLines(filename string, lines []string) error {
	data := strings.Join(lines, "\n")
	return os.WriteFile(filename, []byte(data), 0644)
}

func readDiagram(filename string) []string {
	file, _ := os.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lines = append(lines, line)
	}

	return lines
}

func travel(diagram []string, pos Pos) int {
	// base case
	if pos.row == len(diagram)-1 {
		return 1
	}

	// cache check
	if v, ok := cache[pos]; ok {
		return v
	}

	var result int

	switch diagram[pos.row+1][pos.col] {
	case '.':
		result = travel(diagram, Pos{row: pos.row + 1, col: pos.col})
	case '^':
		left := travel(diagram, Pos{pos.row + 1, pos.col - 1})
		right := travel(diagram, Pos{pos.row + 1, pos.col + 1})
		result = left + right
	default:
		panic("unexpected character in diagram")
	}

	// store in cache BEFORE returning
	cache[pos] = result
	return result
}

var cache = make(map[Pos]int)

func main() {
	start := time.Now()

	// === TEST 1 ===
	test_splits := readAndSolveInput("./test.txt")
	fmt.Printf("TEST 1: There were %d splits in the TEST file (wanted: 21).\n", test_splits)

	// === PART 1 ===
	splits := readAndSolveInput("./input.txt")
	fmt.Printf("PART 1: There were %d splits.\n", splits)

	// === TEST 2 ===
	testDiagram := readDiagram("./test.txt")
	testStartingPoint := strings.Index(testDiagram[0], "S")
	testActiveTimelines := travel(testDiagram, Pos{row: 1, col: testStartingPoint})
	fmt.Printf("TEST 2: Active timelines: %d (wanted: 40)\n", testActiveTimelines)

	// === PART 2 ===
	cache = make(map[Pos]int)
	diagram := readDiagram("./input.txt")
	startingPoint := strings.Index(diagram[0], "S")
	activeTimelines := travel(diagram, Pos{row: 1, col: startingPoint})
	fmt.Printf("PART 2: Active timelines: %d\n", activeTimelines)

	log.Printf("Execution time: %s", time.Since(start))
}
