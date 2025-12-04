package main

import "testing"

func TestGetRanksSortedKeys(t *testing.T) {
	tests := []struct {
		ranks map[int][]int
		want  []int
	}{
		{map[int][]int{1: {0}, 3: {1}, 2: {2}}, []int{3, 2, 1}},
		{map[int][]int{5: {0}, 7: {1}, 6: {2}}, []int{7, 6, 5}},
		{map[int][]int{9: {0}, 4: {1}, 8: {2}}, []int{9, 8, 4}},
	}

	for _, tt := range tests {
		got := getRanksSortedKeys(tt.ranks)
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("getRanksSortedKeys(%v) = %v; want %v", tt.ranks, got, tt.want)
				break
			}
		}
	}
}

func TestPart1GetMaxJoltageFromBank(t *testing.T) {
	tests := []struct {
		bank string
		want int
	}{
		{"987654321111111", 98},
		{"811111111111119", 89},
		{"234234234234278", 78},
		{"818181911112111", 92},
	}

	for _, tt := range tests {
		got := part1GetMaxJoltageFromBank(tt.bank)
		if got != tt.want {
			t.Errorf("isIDValid(%s) = %d; want %d", tt.bank, got, tt.want)
		}
	}
}

func TestPart2GetMaxJoltageFromBank(t *testing.T) {
	tests := []struct {
		bank string
		want int
	}{
		{"987654321111111", 987654321111},
		{"811111111111119", 811111111119},
		{"234234234234278", 434234234278},
		{"818181911112111", 888911112111},
	}

	for _, tt := range tests {
		got := part2GetMaxJoltageFromBank(tt.bank)
		if got != tt.want {
			t.Errorf("isIDValid(%s) = %d; want %d", tt.bank, got, tt.want)
		}
	}
}
