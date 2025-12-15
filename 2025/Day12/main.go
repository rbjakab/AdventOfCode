package main

import (
	"log"
	"time"
)

func countFittableRegions(regions []Region) int {
	count := 0
	for _, region := range regions {
		maxBoundWidth := region.width / 3
		maxBoundHeight := region.height / 3

		ok := sum(region.shapes) <= maxBoundWidth*maxBoundHeight

		if ok {
			count++
		}
	}
	return count
}

func main() {
	start := time.Now()
	_, regions := readInput("input.txt")

	log.Println("=== PART 1 ===")
	count := countFittableRegions(regions)
	log.Printf("Fittable regions: %d", count)

	log.Printf("Execution time: %s", time.Since(start))
}
