package main

import "testing"

func TestIncrease(t *testing.T) {
	tests := []struct {
		name string
		in   []bool
		want []bool
	}{
		{
			name: "simple increment",
			in:   []bool{false, false, false, true},
			want: []bool{false, false, true, false},
		},
		{
			name: "carry propagation",
			in:   []bool{false, false, true, true},
			want: []bool{false, true, false, false},
		},
		{
			name: "multiple carries",
			in:   []bool{false, true, true, true},
			want: []bool{true, false, false, false},
		},
		{
			name: "overflow wraps through all bits",
			in:   []bool{true, true, true, true},
			want: []bool{false, false, false, false},
		},
	}

	for _, tt := range tests {
		// copy input because the function mutates the slice
		inCopy := append([]bool(nil), tt.in...)

		increase(inCopy)

		for i := range tt.want {
			if inCopy[i] != tt.want[i] {
				t.Errorf("%s: got %v, want %v", tt.name, inCopy, tt.want)
				break
			}
		}
	}
}
