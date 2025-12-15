package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Range struct {
	Start int
	End   int
}

func readInput(filePath string) []Range {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	line := strings.TrimSpace(string(data))
	raw := strings.Split(line, ",")

	ranges := make([]Range, 0, len(raw))

	for _, r := range raw {
		r = strings.TrimSpace(r)

		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			log.Fatalf("invalid range format: %s", r)
		}

		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			log.Fatalf("invalid number in range: %s", r)
		}

		ranges = append(ranges, Range{Start: start, End: end})
	}

	return ranges
}

func part1IsIDValid(id int) bool {
	sID := strconv.Itoa(id)

	if len(sID)%2 != 0 {
		return true
	}

	return sID[:len(sID)/2] != sID[len(sID)/2:]
}

func part2IsIDValid(id int) bool {
	sID := strconv.Itoa(id)
	letters := 1

	for letters <= len(sID)/2 {
		tmp := ""

		for len(tmp) < len(sID) {
			tmp += sID[:letters]
		}

		if tmp == sID {
			return false
		}

		letters += 1
	}

	return true
}

func main() {
	start := time.Now()
	ranges := readInput("./input.txt")

	fmt.Println("=== First Problem ===")
	results := 0
	for _, r := range ranges {
		count := 0
		for i := r.Start; i <= r.End; i++ {
			if !part1IsIDValid(i) {
				count += i
			}
		}
		results += count
	}
	fmt.Println("Results:", results)
	fmt.Println("=== First Problem ===")

	fmt.Println("=== Second Problem ===")
	results = 0
	for _, r := range ranges {
		count := 0
		for i := r.Start; i <= r.End; i++ {
			if !part2IsIDValid(i) {
				count += i
			}
		}
		results += count
	}
	fmt.Println("Results:", results)
	fmt.Println("=== Second Problem ===")
	log.Printf("Execution time: %s", time.Since(start))
}
