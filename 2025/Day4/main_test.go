package main

import "testing"

func TestExtendArray(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]string
		expected [][]string
	}{
		{
			name: "3x3 grid",
			input: [][]string{
				{"a", "b", "c"},
				{"d", "e", "f"},
				{"g", "h", "i"},
			},
			expected: [][]string{
				{".", ".", ".", ".", "."},
				{".", "a", "b", "c", "."},
				{".", "d", "e", "f", "."},
				{".", "g", "h", "i", "."},
				{".", ".", ".", ".", "."},
			},
		},
		{
			name: "2x4 grid",
			input: [][]string{
				{"1", "2", "3", "4"},
				{"5", "6", "7", "8"},
			},
			expected: [][]string{
				{".", ".", ".", ".", ".", "."},
				{".", "1", "2", "3", "4", "."},
				{".", "5", "6", "7", "8", "."},
				{".", ".", ".", ".", ".", "."},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extendArray(tt.input)
			if !equal(result, tt.expected) {
				t.Errorf("extendArray() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// equal compares two [][]string slices for equality.
func equal(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func TestCountNeighbors(t *testing.T) {
	tests := []struct {
		name     string
		grid     [][]string
		row      int
		col      int
		expected int
	}{
		{
			name: "row",
			grid: [][]string{
				{".", ".", ".", ".", ".", "."},
				{".", ".", ".", ".", ".", "."},
				{".", "@", "x", "@", ".", "."},
				{".", ".", ".", ".", ".", "."},
			},
			row:      2,
			col:      2,
			expected: 2,
		},
		{
			name: "column",
			grid: [][]string{
				{".", ".", ".", ".", ".", "."},
				{".", ".", "@", ".", ".", "."},
				{".", ".", "x", ".", ".", "."},
				{".", ".", "@", ".", ".", "."},
			},
			row:      2,
			col:      2,
			expected: 2,
		},
		{
			name: "edges",
			grid: [][]string{
				{".", ".", ".", ".", ".", "."},
				{".", "@", ".", "@", ".", "."},
				{".", ".", "x", ".", ".", "."},
				{".", "@", ".", "@", ".", "."},
			},
			row:      2,
			col:      2,
			expected: 4,
		},
		{
			name: "middle",
			grid: [][]string{
				{".", ".", ".", ".", ".", "."},
				{".", ".", "@", ".", ".", "."},
				{".", "@", "@", "@", ".", "."},
				{".", ".", "@", ".", ".", "."},
			},
			row:      2,
			col:      2,
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countNeighbors(tt.grid, tt.row, tt.col)
			if result != tt.expected {
				t.Errorf("countNeighbors() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCountAccessibleRolls(t *testing.T) {
	tests := []struct {
		name     string
		array    [][]string
		expected int
	}{
		{
			name: "test",
			array: [][]string{
				{".", ".", "@", "@", ".", "@", "@", "@", "@", "."},
				{"@", "@", "@", ".", "@", ".", "@", ".", "@", "@"},
				{"@", "@", "@", "@", "@", ".", "@", ".", "@", "@"},
				{"@", ".", "@", "@", "@", "@", ".", ".", "@", "."},
				{"@", "@", ".", "@", "@", "@", "@", ".", "@", "@"},
				{".", "@", "@", "@", "@", "@", "@", "@", ".", "@"},
				{".", "@", ".", "@", ".", "@", ".", "@", "@", "@"},
				{"@", ".", "@", "@", "@", ".", "@", "@", "@", "@"},
				{".", "@", "@", "@", "@", "@", "@", "@", "@", "."},
				{"@", ".", "@", ".", "@", "@", "@", ".", "@", "."},
			},
			expected: 13,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extendedArray := extendArray(tt.array)
			result := countAccessibleRolls(extendedArray)
			if result != tt.expected {
				t.Errorf("countAccessibleRolls() = %v, want %v", result, tt.expected)
			}
		})
	}
}
