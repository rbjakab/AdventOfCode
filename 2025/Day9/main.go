package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

type Tile struct {
	x, y int
}

func readRedTiles(filename string) ([]Tile, Tile, int, int, int, int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tiles []Tile
	maxX, maxY := 0, 0
	minX, minY := 0, 0
	var topLeftTile Tile

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		var t Tile
		if _, err := fmt.Sscanf(line, "%d,%d", &t.x, &t.y); err == nil {
			tiles = append(tiles, t)
			if t.x > maxX {
				maxX = t.x
			}
			if t.y > maxY {
				maxY = t.y
			}
			if len(tiles) == 1 {
				topLeftTile = t
			} else if t.x < topLeftTile.x && t.y < topLeftTile.y {
				topLeftTile = t
			}
			if t.x < minX {
				minX = t.x
			}
			if t.y < minY {
				minY = t.y
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading %s: %v", filename, err)
	}

	for i := range tiles {
		tiles[i].x -= minX
		tiles[i].y -= minY
	}
	maxX -= minX
	maxY -= minY

	return tiles, topLeftTile, minX, minY, maxX, maxY
}

func abs(a int) int {
	return int(math.Abs(float64(a)))
}

func rectangleArea(a, b Tile) int64 {
	return int64(abs(a.x-b.x)+1) * int64(abs(a.y-b.y)+1)
}

func greenTilesBetween(a, b Tile) []Tile {
	var greens []Tile
	d := max(abs(a.x-b.x), abs(a.y-b.y)) - 1

	if a.x == b.x { // vertical
		x := a.x
		y0 := min(a.y, b.y)
		for i := 1; i <= d; i++ {
			greens = append(greens, Tile{x, y0 + i})
		}
		return greens
	}

	if a.y == b.y { // horizontal
		y := a.y
		x0 := min(a.x, b.x)
		for i := 1; i <= d; i++ {
			greens = append(greens, Tile{x0 + i, y})
		}
		return greens
	}

	return greens
}

func computeGreenLoop(redTiles []Tile) []Tile {
	var greens []Tile
	for i := 0; i < len(redTiles)-1; i++ {
		greens = append(greens, greenTilesBetween(redTiles[i], redTiles[i+1])...)
	}
	greens = append(greens, greenTilesBetween(redTiles[len(redTiles)-1], redTiles[0])...)
	return greens
}

type TilePair struct {
	tileA Tile
	tileB Tile
	area  int64
}

func BuildInsideMask(redGreenSet map[Tile]bool, topLeftTile Tile) map[Tile]struct{} {
	inside := make(map[Tile]struct{}, len(redGreenSet)/4) // guess to reduce realloc
	visited := inside                                     // alias
	empty := struct{}{}

	// BFS queue using indices instead of slice pops
	queue := make([]Tile, 0, 1024)
	queue = append(queue, Tile{topLeftTile.x + 1, topLeftTile.y + 1})

	for i := 0; i < len(queue); i++ {
		curr := queue[i]

		// 4-neighbors inline (avoid slice allocation)
		nx := curr.x + 1
		ny := curr.y
		n := Tile{nx, ny}
		if !redGreenSet[n] {
			if _, ok := visited[n]; !ok {
				visited[n] = empty
				queue = append(queue, n)
			}
		}

		n = Tile{curr.x - 1, curr.y}
		if !redGreenSet[n] {
			if _, ok := visited[n]; !ok {
				visited[n] = empty
				queue = append(queue, n)
			}
		}

		n = Tile{curr.x, curr.y + 1}
		if !redGreenSet[n] {
			if _, ok := visited[n]; !ok {
				visited[n] = empty
				queue = append(queue, n)
			}
		}

		n = Tile{curr.x, curr.y - 1}
		if !redGreenSet[n] {
			if _, ok := visited[n]; !ok {
				visited[n] = empty
				queue = append(queue, n)
			}
		}

		// OPTIONAL: keep logging if needed
		if i%5000000 == 0 && i > 0 {
			log.Printf("Filled %d tiles, queue=%d", len(visited), len(queue))
		}

	}

	return inside
}

func isRectangleInsideInOthers(a, b Tile, rectangles map[Rect]struct{}) bool {
	for rectangle := range rectangles {
		topLeft := Tile{min(rectangle[0].x, rectangle[1].x), min(rectangle[0].y, rectangle[1].y)}
		bottomRight := Tile{max(rectangle[0].x, rectangle[1].x), max(rectangle[0].y, rectangle[1].y)}

		if (topLeft.x <= a.x && a.x <= bottomRight.x) && (topLeft.y <= a.y && a.y <= bottomRight.y) &&
			(topLeft.x <= b.x && b.x <= bottomRight.x) && (topLeft.y <= b.y && b.y <= bottomRight.y) {
			return true
		}
	}

	return false
}

func isRectangleValid(a, b Tile, inside map[Tile]struct{}) bool {
	left := min(a.x, b.x)
	right := max(a.x, b.x)
	top := min(a.y, b.y)
	bottom := max(a.y, b.y)

	// Check top and bottom edges
	for x := left; x <= right; x++ {
		if _, ok := inside[Tile{x, top}]; !ok {
			return false
		}
		if _, ok := inside[Tile{x, bottom}]; !ok {
			return false
		}
	}

	// Check left and right edges
	for y := top; y <= bottom; y++ {
		if _, ok := inside[Tile{left, y}]; !ok {
			return false
		}
		if _, ok := inside[Tile{right, y}]; !ok {
			return false
		}
	}

	// check inside
	for x := left + 1; x < right; x++ {
		for y := top + 1; y < bottom; y++ {
			if _, ok := inside[Tile{x, y}]; !ok {
				return false
			}
		}
	}

	return true
}

type Rect [2]Tile

func main() {
	// === Part 1 ===
	redTiles, topLeftTile, minX, minY, maxX, maxY := readRedTiles("input.txt")
	log.Printf("Read %d red tiles, grid size: %d x %d\n", len(redTiles), maxX, maxY)
	log.Printf("Start position: (%d, %d)\n", minX, minY)

	largestArea := int64(-1)
	for i := range redTiles {
		for j := range redTiles {
			if i != j {
				area := rectangleArea(redTiles[i], redTiles[j])
				if area > largestArea {
					largestArea = area
				}
			}
		}
	}
	log.Println("Part 1:", largestArea)

	// === Part 2 ===
	start := time.Now()

	greenTiles := computeGreenLoop(redTiles)
	redGreenSet := make(map[Tile]bool, len(redTiles)+len(greenTiles))
	for _, t := range redTiles {
		redGreenSet[t] = true
	}
	for _, t := range greenTiles {
		redGreenSet[t] = true
	}

	// Build ordered list of rectangle pairs
	var orderedTilePairs []TilePair
	for i := range redTiles {
		for j := i + 1; j < len(redTiles); j++ {
			orderedTilePairs = append(orderedTilePairs, TilePair{
				tileA: redTiles[i],
				tileB: redTiles[j],
				area:  rectangleArea(redTiles[i], redTiles[j]),
			})
		}
	}
	sort.Slice(orderedTilePairs, func(i, j int) bool {
		return orderedTilePairs[i].area < orderedTilePairs[j].area
	})

	log.Println("Building inside mask (scanline fill)...")
	inside := BuildInsideMask(redGreenSet, topLeftTile)
	log.Println("Inside mask ready, length:", len(inside))
	log.Println("Top Left Tile:", topLeftTile)
	log.Println("Grid size:", maxX, maxY)
	wrongRectangles := make(map[Rect]struct{})

	// Check rectangles using fast perimeter-only test
	largestArea = -1
	for i, pair := range orderedTilePairs {
		if i%100 == 0 {
			log.Printf("Checking rectangle %d / %d with area %d | largest: %d\n", i, len(orderedTilePairs), pair.area, largestArea)
		}

		currRectangle := Rect{pair.tileA, pair.tileB}

		if isRectangleInsideInOthers(pair.tileA, pair.tileB, wrongRectangles) {
			wrongRectangles[currRectangle] = struct{}{}
			continue
		}

		if isRectangleValid(pair.tileA, pair.tileB, inside) {
			largestArea = pair.area
		} else {
			wrongRectangles[currRectangle] = struct{}{}
		}
	}

	log.Println("Part 2:", largestArea)
	log.Printf("Program ran for %s\n", time.Since(start))
}
