package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Shape [][]bool

type Present struct {
	id    int
	shape Shape
}

type Region struct {
	width  int
	height int
	shapes []int
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func readInput(filename string) ([]Present, []Region) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("open input: %v", err)
	}
	defer file.Close()

	var presents []Present
	var regions []Region

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// ---------- PRESENT SHAPES ----------
		// Format: "0:", "1:", ...
		if strings.HasSuffix(line, ":") && isNumber(strings.TrimSuffix(line, ":")) {
			idStr := strings.TrimSuffix(line, ":")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("invalid shape id %q", line)
			}

			var shape Shape

			// Read next 3 lines (3x3 grid)
			for i := 0; i < 3; i++ {
				if !scanner.Scan() {
					log.Fatal("unexpected EOF while reading shape")
				}
				rowLine := scanner.Text()
				row := make([]bool, len(rowLine))
				for j, c := range rowLine {
					row[j] = (c == '#')
				}
				shape = append(shape, row)
			}

			presents = append(presents, Present{
				id:    id,
				shape: shape,
			})
			continue
		}

		// ---------- REGIONS ----------
		// Format: "47x38: 27 37 25 25 36 29"
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			log.Fatalf("invalid line: %q", line)
		}

		// Parse width x height
		dims := strings.Split(parts[0], "x")
		if len(dims) != 2 {
			log.Fatalf("invalid region dims: %q", parts[0])
		}

		width, err := strconv.Atoi(dims[0])
		if err != nil {
			log.Fatal(err)
		}
		height, err := strconv.Atoi(dims[1])
		if err != nil {
			log.Fatal(err)
		}

		// Parse present IDs
		fields := strings.Fields(parts[1])
		var presentIDs []int
		for _, f := range fields {
			id, err := strconv.Atoi(f)
			if err != nil {
				log.Fatal(err)
			}
			presentIDs = append(presentIDs, id)
		}

		regions = append(regions, Region{
			width:  width,
			height: height,
			shapes: presentIDs,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("scan input: %v", err)
	}

	return presents, regions
}

func sum(a []int) int {
	total := 0
	for _, e := range a {
		total += e
	}
	return total
}
