package main

import (
	"testing"
)

func TestIsIngredientFresh(t *testing.T) {
	tests := []struct {
		name   string
		ranges []IntRange
		id     int
		want   bool
	}{
		{
			name: "ID below all ranges",
			ranges: []IntRange{
				{3, 5}, {10, 14}, {16, 20}, {12, 18},
			},
			id:   1,
			want: false,
		},
		{
			name: "ID at start of first range",
			ranges: []IntRange{
				{3, 5}, {10, 14}, {16, 20}, {12, 18},
			},
			id:   3,
			want: true,
		},
		{
			name: "ID inside first range",
			ranges: []IntRange{
				{3, 5}, {10, 14}, {16, 20}, {12, 18},
			},
			id:   5,
			want: true,
		},
		{
			name: "ID between ranges",
			ranges: []IntRange{
				{3, 5}, {10, 14}, {16, 20}, {12, 18},
			},
			id:   8,
			want: false,
		},
		{
			name: "ID inside overlapping ranges",
			ranges: []IntRange{
				{3, 5}, {10, 14}, {16, 20}, {12, 18},
			},
			id:   11,
			want: true,
		},
		{
			name: "ID in last range",
			ranges: []IntRange{
				{3, 5}, {10, 14}, {16, 20}, {12, 18},
			},
			id:   17,
			want: true,
		},
		{
			name: "ID above all ranges",
			ranges: []IntRange{
				{3, 5}, {10, 14}, {16, 20}, {12, 18},
			},
			id:   32,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIngredientFresh(tt.ranges, tt.id)
			if got != tt.want {
				t.Errorf("isIngredientFresh(%d) = %v; want %v", tt.id, got, tt.want)
			}
		})
	}
}

func TestMergeRanges(t *testing.T) {
	tests := []struct {
		name   string
		ranges []IntRange
		want   []IntRange
	}{
		{
			name: "Simple overlap merge",
			ranges: []IntRange{
				{3, 5},
				{10, 14},
				{16, 20},
				{12, 18},
			},
			want: []IntRange{
				{3, 5},
				{10, 20},
			},
		},
		{
			name: "Multiple overlaps collapse into one",
			ranges: []IntRange{
				{3, 5},
				{10, 14},
				{14, 24},
				{12, 18},
			},
			want: []IntRange{
				{3, 5},
				{10, 24},
			},
		},
		{
			name: "No merges",
			ranges: []IntRange{
				{3, 5},
				{10, 14},
				{16, 18},
				{20, 20},
			},
			want: []IntRange{
				{3, 5},
				{10, 14},
				{16, 18},
				{20, 20},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeRanges(tt.ranges)
			if len(got) != len(tt.want) {
				t.Fatalf("mergeRanges() returned %d ranges; want %d", len(got), len(tt.want))
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("range %d = %v; want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
