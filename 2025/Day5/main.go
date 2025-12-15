package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type IntRange struct {
	Start int
	End   int
}

func readInput(filename string) ([]IntRange, []int) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var ranges []IntRange
	var ids []int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				log.Fatalf("invalid range line: %q", line)
			}
			start, err1 := strconv.Atoi(parts[0])
			end, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				log.Fatalf("invalid range numbers: %q", line)
			}
			ranges = append(ranges, IntRange{Start: start, End: end})
			continue
		}

		id, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("invalid ingredient ID: %q", line)
		}
		ids = append(ids, id)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	return ranges, ids
}

func isIngredientFresh(ranges []IntRange, id int) bool {
	for _, r := range ranges {
		if id >= r.Start && id <= r.End {
			return true
		}
	}
	return false
}

func mergeRanges(ranges []IntRange) []IntRange {
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Start < ranges[j].Start
	})

	merged := []IntRange{ranges[0]}

	for _, r := range ranges[1:] {
		last := &merged[len(merged)-1]
		if r.Start <= last.End+1 {
			if r.End > last.End {
				last.End = r.End
			}
		} else {
			merged = append(merged, r)
		}
	}
	return merged
}

func main() {
	start := time.Now()

	ranges, ids := readInput("./input.txt")

	freshCount := 0
	for _, id := range ids {
		if isIngredientFresh(ranges, id) {
			freshCount++
		}
	}
	fmt.Println("Number of available ingredient IDs that are fresh:", freshCount)

	merged := mergeRanges(ranges)
	total := 0
	for _, r := range merged {
		total += r.End - r.Start + 1
	}
	fmt.Println("Number of TOTAL available ingredient IDs that are fresh:", total)
	log.Printf("Execution time: %s", time.Since(start))
}
