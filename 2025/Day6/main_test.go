package main

import (
	"testing"
)

func TestApplyOperator(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		operator string
		want     int
	}{
		{
			name:     "Additon",
			numbers:  []int{1, 2},
			operator: "+",
			want:     3,
		},
		{
			name:     "Multiplication",
			numbers:  []int{1, 2},
			operator: "*",
			want:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyOperator(tt.numbers, tt.operator)
			if got != tt.want {
				t.Errorf("applyOperator(%v, %s) = %v; want %d", tt.numbers, tt.operator, got, tt.want)
			}
		})
	}
}

func TestCalculateRow(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		operator string
		want     int
	}{
		{
			name:     "Addition",
			numbers:  []int{1, 2, 3},
			operator: "+",
			want:     6,
		},
		{
			name:     "Multiplication",
			numbers:  []int{1, 2, 3, 4},
			operator: "*",
			want:     24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateRow(tt.numbers, tt.operator)
			if got != tt.want {
				t.Errorf("calculateRow(%v, %s) = %v; want %d", tt.numbers, tt.operator, got, tt.want)
			}
		})
	}
}
