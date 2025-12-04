package main

import "testing"

func TestSolvePart2Example(t *testing.T) {
	words := []string{
		"L68",
		"L30",
		"R48",
		"L5",
		"R60",
		"L55",
		"L1",
		"L99",
		"R14",
		"L82",
	}

	want := 6
	got := solvePart2(words)

	if got != want {
		t.Errorf("solvePart2(example) = %d; want %d", got, want)
	}
}

func TestBigRotates(t *testing.T) {
	tests := []struct {
		count     int
		val       int
		wantCount int
		wantVal   int
	}{
		{0, 420, 4, 20},
		{0, 350, 3, 50},
		{0, 500, 5, 0},
	}

	for _, tt := range tests {
		gotCount, gotVal := bigRotates(tt.count, tt.val)
		if gotCount != tt.wantCount || gotVal != tt.wantVal {
			t.Errorf("bigRotates(%d, %d) = (count=%d, val=%d); want (count=%d, val=%d)",
				tt.count, tt.val,
				gotCount, gotVal,
				tt.wantCount, tt.wantVal,
			)
		}
	}
}

func TestRotatePart2(t *testing.T) {
	tests := []struct {
		direction byte
		val       int
		curr      int
		wantPos   int
		wantCount int
	}{
		// Simple cases
		{'R', 1, 99, 0, 1},
		{'R', 5, 96, 1, 1},

		// No wrap
		{'R', 10, 50, 60, 0},
		{'L', 10, 50, 40, 0},

		// Wraps
		{'R', 60, 70, 30, 1},
		{'L', 70, 45, 75, 1},

		// Starting from zero
		{'L', 1, 0, 99, 0},
		{'R', 70, 0, 70, 0},
		{'L', 70, 0, 30, 0},
	}

	for _, tt := range tests {
		gotPos, gotCount := rotatePart2(tt.direction, tt.val, tt.curr)
		if gotPos != tt.wantPos || gotCount != tt.wantCount {
			t.Errorf("rotatePart2(%c, %d, %d) = (pos=%d, count=%d); want (pos=%d, count=%d)",
				tt.direction, tt.val, tt.curr,
				gotPos, gotCount,
				tt.wantPos, tt.wantCount,
			)
		}
	}
}
