package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// readCharGrid reads a text file where each line is a string,
// splits each line into individual characters, and returns a
// 2D slice of strings representing a character grid.
//
// Example line: ".@.@"
// Becomes: []string{".", "@", ".", "@"}
//
// Used to load the puzzle input as a grid.
func readCharGrid(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var grid [][]string

	// Read the file line-by-line
	for scanner.Scan() {
		line := scanner.Text()
		chars := strings.Split(line, "") // split into characters
		grid = append(grid, chars)
	}

	if err := scanner.Err(); err != nil {
		return nil
	}

	return grid
}

// extendArray surrounds the given 2D grid with a border of "."
// characters. This is used to avoid bounds checking when scanning
// neighbors, because all original edge elements now have safe
// padding around them.
//
// Example (3x3 input):
// [ a b c ]
// [ d e f ]
// [ g h i ]
//
// Output (5x5):
// [ . . . . . ]
// [ . a b c . ]
// [ . d e f . ]
// [ . g h i . ]
// [ . . . . . ]
func extendArray(array [][]string) [][]string {
	var extendedArray [][]string
	extendedRowLength := len(array[0]) + 2

	// Create an empty row filled with "."
	emptyRow := make([]string, extendedRowLength)
	for i := range extendedRowLength {
		emptyRow[i] = "."
	}

	// Add top padding row
	extendedArray = append(extendedArray, emptyRow)

	// Add "." to the start and end of every original row
	for _, row := range array {
		extendedRow := append([]string{"."}, row...)
		extendedRow = append(extendedRow, ".")
		extendedArray = append(extendedArray, extendedRow)
	}

	// Add bottom padding row
	extendedArray = append(extendedArray, emptyRow)

	return extendedArray
}

// countNeighbors counts how many '@' characters appear in the
// eight adjacent positions around the cell at (row, col).
//
// It assumes the grid has been padded so that indexing is safe.
// The center cell itself is NOT counted.
func countNeighbors(grid [][]string, row int, col int) int {
	rolls := 0

	// Scan the 3x3 neighborhood around (row, col)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Skip the center cell
			if i == 0 && j == 0 {
				continue
			}
			if grid[row+i][col+j] == "@" {
				rolls++
			}
		}
	}

	return rolls
}

// countAccessibleRolls returns the number of '@' cells that have
// fewer than 4 neighboring '@' cells. These are considered
// "accessible rolls" per the puzzle rules.
func countAccessibleRolls(array [][]string) int {
	rowLength := len(array[0])
	colLength := len(array)

	accesibleRolls := 0

	// Skip the outer border (added by extendArray)
	for i := 1; i < rowLength-1; i++ {
		for j := 1; j < colLength-1; j++ {

			// Only consider '@' cells
			if array[i][j] != "@" {
				continue
			}

			neighborRolls := countNeighbors(array, i, j)

			// Accessible if < 4 neighbors
			if neighborRolls < 4 {
				accesibleRolls++
			}
		}
	}

	return accesibleRolls
}

// removeAccessibleRolls removes all '@' cells that have fewer than
// 4 '@' neighbors and returns the modified grid and the number of
// cells removed in this pass.
func removeAccessibleRolls(array [][]string) ([][]string, int) {
	rowLength := len(array[0])
	colLength := len(array)

	removedRolls := 0

	// Iterate only inside the padded region
	for i := 1; i < rowLength-1; i++ {
		for j := 1; j < colLength-1; j++ {

			if array[i][j] != "@" {
				continue
			}

			neighborRolls := countNeighbors(array, i, j)

			// Remove '@' if it has fewer than 4 neighbors
			if neighborRolls < 4 {
				array[i][j] = "."
				removedRolls++
			}
		}
	}

	return array, removedRolls
}

// main loads the input grid, pads it, counts accessible rolls (Part 1),
// then repeatedly removes accessible rolls until none remain (Part 2).
func main() {
	start := time.Now()

	array := readCharGrid("./input.txt")
	extendedArray := extendArray(array)

	// === PART 1 ===
	accesibleRolls := countAccessibleRolls(extendedArray)
	fmt.Println("Accessible rolls:", accesibleRolls)

	// === PART 2 ===
	totalRemovedRolls := 0
	var removedRolls int

	// Repeatedly delete '@' clusters until stable
	for countAccessibleRolls(extendedArray) > 0 {
		extendedArray, removedRolls = removeAccessibleRolls(extendedArray)
		totalRemovedRolls += removedRolls
	}

	fmt.Println("Removed rolls:", totalRemovedRolls)
	log.Printf("Execution time: %s", time.Since(start))
}
