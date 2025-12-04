package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func getRanksSortedKeys(ranks map[int][]int) []int {
	keys := make([]int, 0, len(ranks))

	for k := range ranks {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	return keys
}

func part1GetMaxJoltageFromBank(bank string) int {
	ranks := map[int][]int{}

	for i := 0; i < len(bank); i++ {
		digit := int(bank[i] - '0')

		if len(ranks[digit]) == 2 {
			continue
		}

		ranks[digit] = append(ranks[digit], i)
	}

	sortedKeys := getRanksSortedKeys(ranks)
	largest := sortedKeys[0]

	if len(ranks[largest]) == 1 && ranks[largest][0] == len(bank)-1 {
		return sortedKeys[1]*10 + sortedKeys[0]
	}

	if len(ranks[largest]) == 2 {
		return largest * 11
	}

	largestToRight := -1

	for i := ranks[largest][0] + 1; i < len(bank); i++ {
		digit := int(bank[i] - '0')

		if digit > largestToRight {
			largestToRight = digit
		}
	}

	return largest*10 + largestToRight

}

// getIndexOfMaxValue returns the index of the largest value in the slice.
// Assumes slice length > 0.
func getIndexOfMaxValue(nums []int) int {
	maxIdx := 0
	maxVal := nums[0]

	for i := range len(nums) {
		if nums[i] > maxVal {
			maxVal = nums[i]
			maxIdx = i
		}
	}
	return maxIdx
}

// part2GetMaxJoltageFromBank selects exactly k digits that form the
// largest possible number, preserving original order.
func part2GetMaxJoltageFromBank(bank string) int {
	// Number of digits to keep
	const k = 12

	// Convert string -> slice of digits
	digits := make([]int, len(bank))
	for i := range bank {
		digits[i] = int(bank[i] - '0')
	}

	// results holds indices of chosen digits
	// Start with -1 to simplify "from" boundary
	results := []int{-1}

	// Choose k digits, one at a time
	for i := range k {
		from := results[len(results)-1] + 1
		to := len(digits) - (k - i) // inclusive range end

		// find max digit in digits[from : to+1]
		selectedIdx := getIndexOfMaxValue(digits[from : to+1])
		results = append(results, from+selectedIdx)
	}

	// Remove the initial -1
	results = results[1:]

	// Build the final number from picked digits
	resultNumber := 0
	for i := 0; i < len(results); i++ {
		power := float64(len(results) - i - 1)
		resultNumber += digits[results[i]] * int(math.Pow(10, power))
	}

	return resultNumber
}

func main() {
	// Open file
	f, err := os.Open("./input.txt")
	if err != nil {
		fmt.Println("failed to open input:", err)
		return
	}
	defer f.Close()

	// Store total output
	total := 0

	// Scan its content line-by-line and work with it
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		joltage := part2GetMaxJoltageFromBank(line)
		total += joltage
	}

	fmt.Println("Total joltage:", total)
}
