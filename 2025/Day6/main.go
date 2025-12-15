package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func readInput(filename string) ([][]int, []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file %q: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var rows [][]int
	var ops []string

	const numberLines = 4
	lineIndex := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)

		if lineIndex < numberLines {
			ints := make([]int, len(fields))
			for i, v := range fields {
				n, err := strconv.Atoi(v)
				if err != nil {
					log.Fatalf("invalid integer %q: %v", v, err)
				}
				ints[i] = n
			}
			rows = append(rows, ints)
			lineIndex++
			continue
		}

		// last line: operators
		ops = fields
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("scanner error: %v", err)
	}

	return rows, ops
}

func readInputLeftToRight(filename string) [][]int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file %q: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("scanner error: %v", err)
	}

	// last line contains operators -> ignore
	lines = lines[:len(lines)-1]

	var result [][]int
	var current []int

	cols := len(lines[0])

	for col := 0; col < cols; col++ {
		var sb strings.Builder

		for _, line := range lines {
			if col < len(line) { // safe indexing
				sb.WriteByte(line[col])
			}
		}

		str := strings.TrimSpace(sb.String())
		if str == "" {
			// end of a number block
			if len(current) > 0 {
				result = append(result, current)
				current = nil
			}
			continue
		}

		n, err := strconv.Atoi(str)
		if err != nil {
			log.Fatalf("invalid number %q: %v", str, err)
		}

		current = append(current, n)
	}

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}

func applyOperator(nums []int, op string) int {
	if len(nums) != 2 {
		log.Fatalf("applyOperator expects exactly 2 numbers, got %v", nums)
	}

	switch op {
	case "+":
		return nums[0] + nums[1]
	case "*":
		return nums[0] * nums[1]
	default:
		log.Fatalf("unsupported operator %q", op)
	}
	return 0 // unreachable
}

func calculateRow(nums []int, op string) int {
	result := nums[0]
	for _, n := range nums[1:] {
		result = applyOperator([]int{result, n}, op)
	}
	return result
}

func solve() {
	start := time.Now()

	grid, ops := readInput("./input.txt")

	// === PART 1 ===
	part1 := 0
	for col := range grid[0] {
		var colValues []int
		for _, row := range grid {
			colValues = append(colValues, row[col])
		}
		part1 += calculateRow(colValues, ops[col])
	}
	fmt.Println("Results:", part1)

	// === PART 2 ===
	part2 := 0
	leftToRight := readInputLeftToRight("./input.txt")
	for i, nums := range leftToRight {
		part2 += calculateRow(nums, ops[i])
	}
	fmt.Println("Left To Right Results:", part2)

	log.Printf("Execution time: %s", time.Since(start))
}

const runs = 1000

func main() {
	start := time.Now()
	for i := 0; i < runs; i++ {
		solve()
	}
	elapsed := time.Since(start)
	fmt.Printf("Avg per run: %v\n", elapsed/time.Duration(runs))
}
