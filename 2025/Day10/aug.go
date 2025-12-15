package main

import (
	"math"
)

const eps = 1e-12

func rref(A [][]float64) [][]float64 {
	m := len(A)
	n := len(A[0])

	r := 0 // current pivot row

	for c := 0; c < n-1 && r < m; c++ {
		// Find pivot row
		pivotRow := r
		maxVal := math.Abs(A[r][c])
		for i := r + 1; i < m; i++ {
			if math.Abs(A[i][c]) > maxVal {
				maxVal = math.Abs(A[i][c])
				pivotRow = i
			}
		}

		// If column is zero, skip it
		if math.Abs(A[pivotRow][c]) < eps {
			continue
		}

		// Swap pivot row into place
		A[r], A[pivotRow] = A[pivotRow], A[r]

		// Normalize pivot row
		pivot := A[r][c]
		for j := c; j < n; j++ {
			A[r][j] /= pivot
			if math.Abs(A[r][j]) < eps {
				A[r][j] = 0
			}
		}

		// Eliminate all other rows
		for i := 0; i < m; i++ {
			if i == r {
				continue
			}
			f := A[i][c]
			if math.Abs(f) < eps {
				continue
			}
			for j := c; j < n; j++ {
				A[i][j] -= f * A[r][j]
				if math.Abs(A[i][j]) < eps {
					A[i][j] = 0
				}
			}
		}

		r++
	}

	return A
}
