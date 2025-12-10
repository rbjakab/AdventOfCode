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
	X, Y float64
}

type Rect [2]Tile

func readRedTiles(filename string) []Tile {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tiles []Tile

	// Initialize mins to +∞ so the first tile sets them
	minX, minY := math.Inf(1), math.Inf(1)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		var x, y float64
		if _, err := fmt.Sscanf(line, "%f,%f", &x, &y); err == nil {

			// track minimums
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}

			tiles = append(tiles, Tile{X: x, Y: y})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading %s: %v", filename, err)
	}

	// Normalize all tiles by subtracting minX, minY
	for i := range tiles {
		tiles[i].X -= minX
		tiles[i].Y -= minY
	}

	return tiles
}

func rectangleArea(a, b Tile) float64 {
	return (math.Abs(a.X-b.X) + 1.0) * (math.Abs(a.Y-b.Y) + 1.0)
}

// orientation test helper
func orient(ax, ay, bx, by, cx, cy float64) float64 {
	return (bx-ax)*(cy-ay) - (by-ay)*(cx-ax)
}

// check if P lies on segment A->B
func onSegment(px, py, ax, ay, bx, by float64) bool {
	if orient(ax, ay, bx, by, px, py) != 0 {
		return false
	}
	if px < min(ax, bx) || px > max(ax, bx) {
		return false
	}
	if py < min(ay, by) || py > max(ay, by) {
		return false
	}
	return true
}

// PointInPolygon returns true if point P is inside or on boundary.
func PointInPolygon(p Tile, poly []Tile) bool {
	px, py := p.X, p.Y
	n := len(poly)

	inside := false

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		xi, yi := poly[i].X, poly[i].Y
		xj, yj := poly[j].X, poly[j].Y

		// 1. Check if point lies exactly on this edge → inside
		if onSegment(px, py, xi, yi, xj, yj) {
			return true
		}

		// 2. Ray-casting: does edge cross the ray horizontally?
		intersects := (yi > py) != (yj > py)
		if intersects {
			// compute x coordinate of intersection
			xIntersect := xi + (py-yi)*(xj-xi)/(yj-yi)
			if px <= xIntersect {
				inside = !inside
			}
		}
	}

	return inside
}

func isRectangleInsideInOthers(a, b Tile, rectangles map[Rect]struct{}) bool {
	for rectangle := range rectangles {
		topLeft := Tile{min(rectangle[0].X, rectangle[1].X), min(rectangle[0].Y, rectangle[1].Y)}
		bottomRight := Tile{max(rectangle[0].X, rectangle[1].X), max(rectangle[0].Y, rectangle[1].Y)}

		if (topLeft.X <= a.X && a.X <= bottomRight.X) && (topLeft.Y <= a.Y && a.Y <= bottomRight.Y) &&
			(topLeft.X <= b.X && b.X <= bottomRight.X) && (topLeft.Y <= b.Y && b.Y <= bottomRight.Y) {
			return true
		}
	}

	return false
}

func isRectangleValid(a, b Tile, edges []Tile, insideArray map[Tile]struct{}) bool {
	left := min(a.X, b.X)
	right := max(a.X, b.X)
	top := min(a.Y, b.Y)
	bottom := max(a.Y, b.Y)

	// Check top and bottom edges
	for x := left; x <= right; x++ {
		topTile := Tile{x, top}
		_, okTop := insideArray[topTile]
		if !okTop {
			inside := PointInPolygon(topTile, edges)
			if inside {
				insideArray[topTile] = struct{}{}
			} else {
				return false
			}
		}

		bottomTile := Tile{x, bottom}
		_, okBottom := insideArray[bottomTile]
		if !okBottom {
			inside := PointInPolygon(bottomTile, edges)
			if inside {
				insideArray[bottomTile] = struct{}{}
			} else {
				return false
			}
		}
	}

	// Check left and right edges
	for y := top; y <= bottom; y++ {
		leftTile := Tile{left, y}
		_, okLeft := insideArray[leftTile]
		if !okLeft {
			inside := PointInPolygon(leftTile, edges)
			if inside {
				insideArray[leftTile] = struct{}{}
			} else {
				return false
			}
		}

		rightTile := Tile{right, y}
		_, okRight := insideArray[rightTile]
		if !okRight {
			inside := PointInPolygon(rightTile, edges)
			if inside {
				insideArray[rightTile] = struct{}{}
			} else {
				return false
			}
		}
	}

	// check inside
	for x := left + 1; x < right; x++ {
		for y := top + 1; y < bottom; y++ {
			tile := Tile{x, y}
			if _, ok := insideArray[tile]; !ok {
				inside := PointInPolygon(tile, edges)
				if inside {
					insideArray[tile] = struct{}{}
				} else {
					return false
				}

			}
		}
	}

	return true
}

func main() {
	// === Part 1 ===
	redTiles := readRedTiles("input.txt")
	log.Printf("Read %d red tiles\n", len(redTiles))

	largestArea := float64(-1)

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

	// === Test 2 ===
	start := time.Now()

	redTiles = readRedTiles("input.txt")
	log.Printf("Testing PointInPolygon with %d red tiles\n", len(redTiles))

	var orderedTiles []Tile
	orderedTiles = append(orderedTiles, redTiles...)
	sort.Slice(orderedTiles, func(i, j int) bool {
		return orderedTiles[i].X < orderedTiles[j].X
	})

	insideArray := make(map[Tile]struct{})
	wrongRectangles := make(map[Rect]struct{})
	largestArea = -1

	for i := 0; i < len(orderedTiles); i++ {
		log.Printf("Checking tile %d/%d, largest so far: %d, inside array: %d, wrong rectangles: %d\n",
			i, len(orderedTiles), int(largestArea), len(insideArray), len(wrongRectangles))

		t1 := orderedTiles[i]
		for j := i + 1; j < len(orderedTiles); j++ {
			t2 := orderedTiles[j]

			currRectangle := Rect{t1, t2}

			if isRectangleInsideInOthers(t1, t2, wrongRectangles) {
				wrongRectangles[currRectangle] = struct{}{}
				continue
			}

			if isRectangleValid(t1, t2, redTiles, insideArray) {
				area := rectangleArea(t1, t2)
				largestArea = max(largestArea, area)
			} else {
				wrongRectangles[currRectangle] = struct{}{}
			}
		}
	}

	log.Println("Largest area float64:", largestArea)
	log.Println("Largest area int64:", int64(largestArea))
	log.Println("Largest area int:", int(largestArea))

	log.Printf("Program ran for %s\n", time.Since(start))
}
