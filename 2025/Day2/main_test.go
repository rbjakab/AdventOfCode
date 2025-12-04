package main

import "testing"

func TestPart1IsIDValid(t *testing.T) {
	tests := []struct {
		id   int
		want bool
	}{
		{101, true},
		{1234, true},
		{1122, true},
		{5678, true},
		{3344, true},
		{123321, true},
		{123456, true},
		{123123, false},
		{11, false},
		{2222, false},
		{4444, false},
		{1188511885, false},
		{446446, false},
		{3577824, true},
	}

	for _, tt := range tests {
		got := part1IsIDValid(tt.id)
		if got != tt.want {
			t.Errorf("isIDValid(%d) = %v; want %v", tt.id, got, tt.want)
		}
	}
}

func TestPart2IsIDValid(t *testing.T) {
	tests := []struct {
		id   int
		want bool
	}{
		{101, true},
		{1234, true},
		{1122, true},
		{5678, true},
		{3344, true},
		{123321, true},
		{123456, true},
		{123123, false},
		{11, false},
		{22, false},
		{999, false},
		{1010, false},
		{446446, false},
		{123123123, false},
		{2121212121, false},
		{12312, true},
		{1231234, true},
	}

	for _, tt := range tests {
		got := part2IsIDValid(tt.id)
		if got != tt.want {
			t.Errorf("isIDValid(%d) = %v; want %v", tt.id, got, tt.want)
		}
	}
}
