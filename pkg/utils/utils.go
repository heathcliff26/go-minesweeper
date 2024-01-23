package utils

func Make2D[T any](x, y int) [][]T {
	if x < 1 || y < 1 {
		return [][]T{}
	}

	matrix := make([][]T, x)
	rows := make([]T, x*y)

	for i := 0; i < x; i++ {
		matrix[i] = rows[i*y : (i+1)*y]
	}
	return matrix
}
