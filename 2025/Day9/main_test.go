package main

import (
	"math"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// small epsilon for float comparisons
const eps = 1e-9

func floatEq(a, b float64) bool {
	return math.Abs(a-b) <= eps
}

func TestRectangleArea(t *testing.T) {
	t.Run("same point", func(t *testing.T) {
		a := Tile{X: 0, Y: 0}
		b := Tile{X: 0, Y: 0}
		got := rectangleArea(a, b)
		want := 1.0 // inclusive of both ends
		if !floatEq(got, want) {
			t.Fatalf("rectangleArea(%v, %v) = %v, want %v", a, b, got, want)
		}
	})

	t.Run("horizontal line", func(t *testing.T) {
		a := Tile{X: 0, Y: 0}
		b := Tile{X: 3, Y: 0}
		got := rectangleArea(a, b)
		// width = |3 - 0| + 1 = 4, height = 1
		want := 4.0
		if !floatEq(got, want) {
			t.Fatalf("rectangleArea(%v, %v) = %v, want %v", a, b, got, want)
		}
	})

	t.Run("vertical line", func(t *testing.T) {
		a := Tile{X: 2, Y: 1}
		b := Tile{X: 2, Y: 4}
		got := rectangleArea(a, b)
		// width = 1, height = |4 - 1| + 1 = 4
		want := 4.0
		if !floatEq(got, want) {
			t.Fatalf("rectangleArea(%v, %v) = %v, want %v", a, b, got, want)
		}
	})

	t.Run("general rectangle", func(t *testing.T) {
		a := Tile{X: 1, Y: 2}
		b := Tile{X: 4, Y: 6}
		got := rectangleArea(a, b)
		// width = |4 - 1| + 1 = 4, height = |6 - 2| + 1 = 5 -> 20
		want := 20.0
		if !floatEq(got, want) {
			t.Fatalf("rectangleArea(%v, %v) = %v, want %v", a, b, got, want)
		}
	})
}

func TestOrient(t *testing.T) {
	tests := []struct {
		name    string
		a, b, c Tile
		sign    int // -1 for negative, 0 for zero, 1 for positive
	}{
		{
			name: "counter clockwise",
			a:    Tile{0, 0},
			b:    Tile{1, 0},
			c:    Tile{1, 1},
			sign: 1,
		},
		{
			name: "clockwise",
			a:    Tile{0, 0},
			b:    Tile{1, 0},
			c:    Tile{1, -1},
			sign: -1,
		},
		{
			name: "collinear along x",
			a:    Tile{0, 0},
			b:    Tile{2, 0},
			c:    Tile{5, 0},
			sign: 0,
		},
		{
			name: "collinear along y",
			a:    Tile{0, 0},
			b:    Tile{0, 2},
			c:    Tile{0, 5},
			sign: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := orient(tt.a.X, tt.a.Y, tt.b.X, tt.b.Y, tt.c.X, tt.c.Y)
			switch tt.sign {
			case 1:
				if !(val > eps) {
					t.Fatalf("orient = %v, want > 0", val)
				}
			case -1:
				if !(val < -eps) {
					t.Fatalf("orient = %v, want < 0", val)
				}
			case 0:
				if !floatEq(val, 0) {
					t.Fatalf("orient = %v, want 0", val)
				}
			}
		})
	}
}

func TestOnSegment(t *testing.T) {
	tests := []struct {
		name string
		p    Tile
		a    Tile
		b    Tile
		want bool
	}{
		{
			name: "point inside segment horizontally",
			p:    Tile{1, 0},
			a:    Tile{0, 0},
			b:    Tile{3, 0},
			want: true,
		},
		{
			name: "point is endpoint a",
			p:    Tile{0, 0},
			a:    Tile{0, 0},
			b:    Tile{3, 0},
			want: true,
		},
		{
			name: "point is endpoint b",
			p:    Tile{3, 0},
			a:    Tile{0, 0},
			b:    Tile{3, 0},
			want: true,
		},
		{
			name: "collinear but outside bounding box",
			p:    Tile{4, 0},
			a:    Tile{0, 0},
			b:    Tile{3, 0},
			want: false,
		},
		{
			name: "not collinear",
			p:    Tile{1, 1},
			a:    Tile{0, 0},
			b:    Tile{3, 0},
			want: false,
		},
		{
			name: "vertical segment, on segment",
			p:    Tile{2, 3},
			a:    Tile{2, 1},
			b:    Tile{2, 5},
			want: true,
		},
		{
			name: "vertical segment, off segment",
			p:    Tile{2, 6},
			a:    Tile{2, 1},
			b:    Tile{2, 5},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := onSegment(tt.p.X, tt.p.Y, tt.a.X, tt.a.Y, tt.b.X, tt.b.Y)
			if got != tt.want {
				t.Fatalf("onSegment(%v,%v,%v) = %v, want %v", tt.p, tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestPointInPolygon_Square(t *testing.T) {
	// Square: (0,0) -> (10,0) -> (10,10) -> (0,10)
	square := []Tile{
		{0, 0},
		{10, 0},
		{10, 10},
		{0, 10},
	}

	tests := []struct {
		name string
		p    Tile
		want bool
	}{
		{"inside center", Tile{5, 5}, true},
		{"inside near edge", Tile{1, 9}, true},
		{"on left edge", Tile{0, 5}, true},
		{"on bottom edge", Tile{5, 0}, true},
		{"on vertex", Tile{0, 0}, true},
		{"outside left", Tile{-1, 5}, false},
		{"outside right", Tile{11, 5}, false},
		{"outside top", Tile{5, 11}, false},
		{"outside bottom", Tile{5, -1}, false},
		{"just outside corner", Tile{-0.1, -0.1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PointInPolygon(tt.p, square)
			if got != tt.want {
				t.Fatalf("PointInPolygon(%v, square) = %v, want %v", tt.p, got, tt.want)
			}
		})
	}
}

func TestPointInPolygon_Concave(t *testing.T) {
	// Simple concave "L" shape:
	// (0,0) -> (4,0) -> (4,1) -> (1,1) -> (1,4) -> (0,4)
	// Visual:
	// ####
	// #
	// #
	// #
	concave := []Tile{
		{0, 0},
		{4, 0},
		{4, 1},
		{1, 1},
		{1, 4},
		{0, 4},
	}

	tests := []struct {
		name string
		p    Tile
		want bool
	}{
		{"inside bottom bar", Tile{2, 0.5}, true},
		{"inside vertical bar", Tile{0.5, 2}, true},
		{"in the concave cutout", Tile{2, 2}, false},
		{"outside far", Tile{10, 10}, false},
		{"on concave edge", Tile{1, 2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PointInPolygon(tt.p, concave)
			if got != tt.want {
				t.Fatalf("PointInPolygon(%v, concave) = %v, want %v", tt.p, got, tt.want)
			}
		})
	}
}

func TestIsRectangleInsideInOthers(t *testing.T) {
	// helper for rect
	r := func(x1, y1, x2, y2 float64) Rect {
		return Rect{Tile{x1, y1}, Tile{x2, y2}}
	}

	rectangles := map[Rect]struct{}{
		r(0, 0, 10, 10): {}, // big one
		r(2, 2, 5, 5):   {}, // smaller
	}

	tests := []struct {
		name       string
		a, b       Tile
		wantInside bool
	}{
		{
			name:       "inside big rectangle",
			a:          Tile{3, 3},
			b:          Tile{4, 4},
			wantInside: true,
		},
		{
			name:       "matches smaller rectangle exactly",
			a:          Tile{2, 2},
			b:          Tile{5, 5},
			wantInside: true,
		},
		{
			name:       "outside all rectangles",
			a:          Tile{-1, -1},
			b:          Tile{1, 1},
			wantInside: false,
		},
		{
			name:       "partially overlapping big, but not fully inside",
			a:          Tile{-1, 5},
			b:          Tile{3, 8},
			wantInside: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRectangleInsideInOthers(tt.a, tt.b, rectangles)
			if got != tt.wantInside {
				t.Fatalf("isRectangleInsideInOthers(%v,%v) = %v, want %v",
					tt.a, tt.b, got, tt.wantInside)
			}
		})
	}
}

func TestIsRectangleValid_SimpleSquarePolygon(t *testing.T) {
	// Polygon is a 5x5 square: (0,0) -> (4,0) -> (4,4) -> (0,4)
	poly := []Tile{
		{0, 0},
		{4, 0},
		{4, 4},
		{0, 4},
	}

	// helper to call with fresh insideArray each time
	check := func(a, b Tile) bool {
		insideArray := make(map[Tile]struct{})
		return isRectangleValid(a, b, poly, insideArray)
	}

	t.Run("rectangle fully inside", func(t *testing.T) {
		a := Tile{1, 1}
		b := Tile{3, 3}
		if !check(a, b) {
			t.Fatalf("isRectangleValid(%v,%v) = false, want true", a, b)
		}
	})

	t.Run("rectangle exactly matches polygon bounds", func(t *testing.T) {
		a := Tile{0, 0}
		b := Tile{4, 4}
		if !check(a, b) {
			t.Fatalf("isRectangleValid(%v,%v) = false, want true", a, b)
		}
	})

	t.Run("rectangle partially outside to the left", func(t *testing.T) {
		a := Tile{-1, 1}
		b := Tile{2, 3}
		if check(a, b) {
			t.Fatalf("isRectangleValid(%v,%v) = true, want false", a, b)
		}
	})

	t.Run("rectangle partially outside to the top", func(t *testing.T) {
		a := Tile{1, 3}
		b := Tile{3, 5}
		if check(a, b) {
			t.Fatalf("isRectangleValid(%v,%v) = true, want false", a, b)
		}
	})

	t.Run("thin horizontal strip inside", func(t *testing.T) {
		a := Tile{0, 2}
		b := Tile{4, 2}
		if !check(a, b) {
			t.Fatalf("isRectangleValid(%v,%v) = false, want true", a, b)
		}
	})
}

func TestIsRectangleValid_CachesInsideArray(t *testing.T) {
	// Reuse the same polygon as previous test
	poly := []Tile{
		{0, 0},
		{4, 0},
		{4, 4},
		{0, 4},
	}
	a := Tile{1, 1}
	b := Tile{3, 3}

	insideArray := make(map[Tile]struct{})

	// First call populates cache
	if !isRectangleValid(a, b, poly, insideArray) {
		t.Fatalf("first isRectangleValid(%v,%v) = false, want true", a, b)
	}

	// insideArray should contain some entries
	if len(insideArray) == 0 {
		t.Fatalf("insideArray is still empty after valid rectangle check, cache not used")
	}

	// Second call uses pre-populated cache, should still be valid
	if !isRectangleValid(a, b, poly, insideArray) {
		t.Fatalf("second isRectangleValid(%v,%v) = false, want true", a, b)
	}
}

func TestReadRedTiles_Normalization(t *testing.T) {
	// Create a temporary directory
	dir := t.TempDir()
	filename := filepath.Join(dir, "tiles.txt")

	// Input has negative and positive coordinates, and unordered
	content := "10,10\n5,7\n-2,3\n8,-1\n"
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	tiles := readRedTiles(filename)

	// We expect 4 tiles
	if len(tiles) != 4 {
		t.Fatalf("readRedTiles() len = %d, want 4", len(tiles))
	}

	// Compute minX, minY from original content:
	// x: 10,5,-2,8 -> minX=-2; y:10,7,3,-1 -> minY=-1
	// After normalization, each tile = (x-minX, y-minY)
	// So expected:
	expected := []Tile{
		{10 - (-2), 10 - (-1)}, // (12,11)
		{5 - (-2), 7 - (-1)},   // (7,8)
		{-2 - (-2), 3 - (-1)},  // (0,4)
		{8 - (-2), -1 - (-1)},  // (10,0)
	}

	// tiles order should match reading order, normalization applied
	if len(tiles) != len(expected) {
		t.Fatalf("len(tiles) = %d, expected %d", len(tiles), len(expected))
	}

	for i := range tiles {
		if !floatEq(tiles[i].X, expected[i].X) || !floatEq(tiles[i].Y, expected[i].Y) {
			t.Fatalf("tile[%d] = %v, expected %v", i, tiles[i], expected[i])
		}
	}
}

// Optional: quick sanity check that Rect can be used as map key as expected
func TestRectKeyEquality(t *testing.T) {
	r1 := Rect{Tile{0, 0}, Tile{4, 4}}
	r2 := Rect{Tile{0, 0}, Tile{4, 4}}
	r3 := Rect{Tile{1, 1}, Tile{4, 4}}

	m := map[Rect]string{
		r1: "foo",
	}
	if got, ok := m[r2]; !ok || got != "foo" {
		t.Fatalf("Rect key equality failed: got %q, ok=%v, want %q,true", got, ok, "foo")
	}
	if _, ok := m[r3]; ok {
		t.Fatalf("Rect key mismatch: r3 should not match r1")
	}
}

// Sanity check to ensure Tile type comparisons behave as expected
func TestTileEquality(t *testing.T) {
	a := Tile{X: 1.0, Y: 2.0}
	b := Tile{X: 1.0, Y: 2.0}
	c := Tile{X: 1.0, Y: 2.000000001}

	if !reflect.DeepEqual(a, b) {
		t.Fatalf("Tiles a and b should be equal, got %v vs %v", a, b)
	}
	if reflect.DeepEqual(a, c) {
		t.Fatalf("Tiles a and c should not be exactly equal")
	}
}
