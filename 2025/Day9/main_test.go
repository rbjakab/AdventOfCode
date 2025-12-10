package main

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestRectangleArea(t *testing.T) {
	tests := []struct {
		a, b Point
		want int64
	}{
		{Point{0, 0}, Point{0, 0}, 1},
		{Point{0, 0}, Point{1, 0}, 2},
		{Point{0, 0}, Point{0, 3}, 4},
		{Point{1, 2}, Point{4, 6}, 20},
		{Point{-3, -2}, Point{-1, -1}, 6},
	}
	for _, tt := range tests {
		got := rectangleArea(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("rectangleArea(%v,%v) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestOrient(t *testing.T) {
	tests := []struct {
		a, b, c Point
		want    int // -1 neg, 0 collinear, 1 pos
	}{
		{Point{0, 0}, Point{1, 0}, Point{1, 1}, 1},
		{Point{0, 0}, Point{1, 0}, Point{1, -1}, -1},
		{Point{0, 0}, Point{2, 2}, Point{4, 4}, 0},
		{Point{3, 3}, Point{5, 5}, Point{6, 6}, 0},
	}
	for _, tt := range tests {
		o := orient(tt.a, tt.b, tt.c)
		got := 0
		if o > 0 {
			got = 1
		} else if o < 0 {
			got = -1
		}
		if got != tt.want {
			t.Errorf("orient(%v,%v,%v) = %d, want %d", tt.a, tt.b, tt.c, got, tt.want)
		}
	}
}

func TestOnSegment(t *testing.T) {
	tests := []struct {
		p, a, b Point
		want    bool
	}{
		{Point{1, 0}, Point{0, 0}, Point{3, 0}, true},
		{Point{0, 0}, Point{0, 0}, Point{3, 0}, true},
		{Point{3, 0}, Point{0, 0}, Point{3, 0}, true},
		{Point{4, 0}, Point{0, 0}, Point{3, 0}, false},
		{Point{1, 1}, Point{0, 0}, Point{3, 0}, false},
		{Point{2, 3}, Point{2, 1}, Point{2, 5}, true},
		{Point{2, 6}, Point{2, 1}, Point{2, 5}, false},
	}
	for _, tt := range tests {
		got := onSegment(tt.p, tt.a, tt.b)
		if got != tt.want {
			t.Errorf("onSegment(%v,%v,%v)= %v, want %v", tt.p, tt.a, tt.b, got, tt.want)
		}
	}
}

func TestPointInPolygon_ConvexSquare(t *testing.T) {
	sq := []Point{
		{0, 0},
		{10, 0},
		{10, 10},
		{0, 10},
	}
	tests := []struct {
		p    Point
		want bool
	}{
		{Point{5, 5}, true},
		{Point{0, 5}, true},
		{Point{10, 10}, true},
		{Point{11, 5}, false},
		{Point{5, -1}, false},
		{Point{-1, -1}, false},
	}
	for _, tt := range tests {
		got := pointInPolygon(tt.p, sq)
		if got != tt.want {
			t.Errorf("pointInPolygon(%v) = %v, want %v", tt.p, got, tt.want)
		}
	}
}

func TestPointInPolygon_Concave(t *testing.T) {
	// L-shape
	poly := []Point{
		{0, 0},
		{4, 0},
		{4, 1},
		{1, 1},
		{1, 4},
		{0, 4},
	}
	tests := []struct {
		p    Point
		want bool
	}{
		{Point{2, 0}, true},
		{Point{0, 2}, true},
		{Point{2, 2}, false}, // concave pocket
		{Point{10, 10}, false},
		{Point{1, 2}, true}, // touching edge
	}
	for _, tt := range tests {
		got := pointInPolygon(tt.p, poly)
		if got != tt.want {
			t.Errorf("pointInPolygon(%v)= %v, want %v", tt.p, got, tt.want)
		}
	}
}

func TestProperSegmentIntersect(t *testing.T) {
	tests := []struct {
		a, b, c, d Point
		want       bool
	}{
		// proper cross
		{Point{0, 0}, Point{3, 3}, Point{0, 3}, Point{3, 0}, true},

		// touching at endpoint → NOT proper
		{Point{0, 0}, Point{2, 2}, Point{2, 2}, Point{4, 2}, false},

		// overlapping collinear → NOT proper
		{Point{0, 0}, Point{5, 0}, Point{2, 0}, Point{4, 0}, false},

		// parallel no-touch
		{Point{0, 0}, Point{0, 3}, Point{1, 0}, Point{1, 3}, false},
	}
	for _, tt := range tests {
		got := properSegmentIntersect(tt.a, tt.b, tt.c, tt.d)
		if got != tt.want {
			t.Errorf("properSegmentIntersect(%v,%v,%v,%v) = %v, want %v",
				tt.a, tt.b, tt.c, tt.d, got, tt.want)
		}
	}
}

func TestRectangleInsidePolygon_SimpleSquare(t *testing.T) {
	poly := []Point{
		{0, 0},
		{10, 0},
		{10, 10},
		{0, 10},
	}

	tests := []struct {
		a, b Point
		want bool
	}{
		{Point{1, 1}, Point{3, 3}, true},
		{Point{0, 0}, Point{10, 10}, true},
		{Point{-1, 1}, Point{3, 3}, false},
		{Point{1, 1}, Point{3, 11}, false},
		{Point{5, 5}, Point{5, 5}, true},
		{Point{0, 5}, Point{10, 5}, true}, // thin strip inside
	}
	for _, tt := range tests {
		got := rectangleInsidePolygon(tt.a, tt.b, poly)
		if got != tt.want {
			t.Errorf("rectInsidePoly(%v,%v)= %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestRectangleInsidePolygon_Concave(t *testing.T) {
	poly := []Point{
		{0, 0},
		{6, 0},
		{6, 2},
		{2, 2},
		{2, 6},
		{0, 6},
	}

	tests := []struct {
		a, b Point
		want bool
	}{
		// this rectangle fits inside both "bars"
		{Point{1, 1}, Point{5, 1}, true},

		// rectangle crossing the concave hole
		{Point{3, 3}, Point{4, 4}, false},

		// big rectangle covering the hole
		{Point{1, 1}, Point{5, 5}, false},
	}
	for _, tt := range tests {
		got := rectangleInsidePolygon(tt.a, tt.b, poly)
		if got != tt.want {
			t.Errorf("concave rectangleInside(%v,%v)= %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

// This test reconstructs the example from the problem statement.
func TestExampleFromAoC(t *testing.T) {
	// Example red tiles (in order) from the puzzle text:
	poly := []Point{
		{7, 1},
		{11, 1},
		{11, 7},
		{9, 7},
		{9, 5},
		{2, 5},
		{2, 3},
		{7, 3},
	}

	// Normalize like code does:
	minX, minY := int64(math.MaxInt64), int64(math.MaxInt64)
	for _, p := range poly {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
	}
	for i := range poly {
		poly[i].X -= minX
		poly[i].Y -= minY
	}

	// After translation we get the same shape.
	// Now test the known valid and invalid rectangles from description.

	// A valid rectangle area=24 between (2,5) & (9,7)
	a := Point{2 - minX, 5 - minY}
	b := Point{9 - minX, 7 - minY}
	if !rectangleInsidePolygon(a, b, poly) {
		t.Fatalf("expected example rectangle (area24) to be valid")
	}

	// A valid rectangle area=35 between (7,1) & (11,7)
	a = Point{7 - minX, 1 - minY}
	b = Point{11 - minX, 7 - minY}
	if !rectangleInsidePolygon(a, b, poly) {
		t.Fatalf("expected example rectangle (area35) to be valid")
	}

	// A valid rectangle area=50 between (2,5) & (11,1)
	a = Point{2 - minX, 5 - minY}
	b = Point{11 - minX, 1 - minY}
	if !rectangleInsidePolygon(a, b, poly) {
		t.Fatalf("expected example rectangle (area50) to be valid")
	}

	// A rectangle that goes outside the loop (should be invalid)
	a = Point{0 - minX, 0 - minY}
	b = Point{12 - minX, 8 - minY}
	if rectangleInsidePolygon(a, b, poly) {
		t.Fatalf("big rectangle should NOT be inside polygon")
	}
}

// Stress-test: random rectangle corners inside bounding box, polygon is a square.
func TestFuzzRectangleInsidePolygon(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	poly := []Point{
		{0, 0},
		{50, 0},
		{50, 50},
		{0, 50},
	}

	for it := 0; it < 5000; it++ {
		x1 := rand.Int63n(60) - 5 // allow some outside
		y1 := rand.Int63n(60) - 5
		x2 := rand.Int63n(60) - 5
		y2 := rand.Int63n(60) - 5

		a := Point{x1, y1}
		b := Point{x2, y2}

		insideCorners := pointInPolygon(a, poly) && pointInPolygon(b, poly)
		if insideCorners {
			// rectangleInsidePolygon must accept
			if !rectangleInsidePolygon(a, b, poly) {
				t.Fatalf("rectangleInsidePolygon failed on rectangle %v,%v", a, b)
			}
		} else {
			// rectangleInsidePolygon must reject
			if rectangleInsidePolygon(a, b, poly) {
				t.Fatalf("rectangleInsidePolygon accepted outside rectangle %v,%v", a, b)
			}
		}
	}
}
