// Package utils
package utils

func Transpose[T any](a [][]T) [][]T {
	row := len(a)
	if row == 0 {
		return [][]T{}
	}

	col := len(a[0])
	if col == 0 {
		return [][]T{}
	}

	b := make([][]T, col)
	for c := range col {
		b[c] = make([]T, row)
		for r := range row {
			b[c][r] = a[r][c]
		}
	}

	return b
}
