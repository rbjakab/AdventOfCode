package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readInput(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Printf("Successfully read in %d words.\n", len(words))

	return words
}

func rotatePart1(direction byte, val, currPoint int) int {
	if direction == 'L' {
		return (currPoint - val + 100) % 100
	}
	return (currPoint + val) % 100
}

func solvePart1(words []string) int {
	currentPoint := 50
	count := 0

	for _, w := range words {
		direction := w[0]
		val, _ := strconv.Atoi(w[1:])
		currentPoint = rotatePart1(direction, val, currentPoint)
		if currentPoint == 0 {
			count++
		}
	}
	return count
}

func bigRotates(count int, val int) (int, int) {
	if val >= 100 {
		count += val / 100
		val = val % 100
	}
	return count, val
}

func rotatePart2(direction byte, val, currPoint int) (int, int) {
	// Determine movement: left = -val, right = +val
	delta := val
	if direction == 'L' {
		delta = -val
	}

	old := currPoint
	newValue := (currPoint + delta) % 100

	// Fix negative modulo result
	if newValue < 0 {
		newValue += 100
	}

	// Did we wrap past 0?
	plusCount := 0
	if direction == 'R' && old+delta >= 100 {
		plusCount = 1
	} else if direction == 'L' && old+delta <= 0 && old != 0 {
		plusCount = 1
	}

	return newValue, plusCount
}

func solvePart2(words []string) int {
	currentPoint := 50
	count := 0
	var plusCount int

	for _, w := range words {
		direction := w[0]
		val, _ := strconv.Atoi(w[1:])
		count, val = bigRotates(count, val)
		currentPoint, plusCount = rotatePart2(direction, val, currentPoint)
		count += plusCount
		// fmt.Printf("%-4s %-4d %-4d %-4d\n", w, currentPoint, plusCount, count)
	}
	return count
}

func main() {
	words := readInput("./input.txt")

	fmt.Println("Number of times at zero, first problem:", solvePart1(words))
	fmt.Println("Number of times at zero, second problem:", solvePart2(words))
}
