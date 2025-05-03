package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMake2D(t *testing.T) {
	tMatrix := []struct{ X, Y int }{
		{-1, -1},
		{-1, 10},
		{10, -1},
		{0, 0},
		{0, 10},
		{10, 0},
		{10, 20},
		{20, 10},
	}

	for _, tCase := range tMatrix {
		name := fmt.Sprintf("%dx%d", tCase.X, tCase.Y)
		t.Run(name, func(t *testing.T) {
			m := Make2D[int](tCase.X, tCase.Y)

			assert := assert.New(t)
			require := require.New(t)

			if tCase.X < 1 || tCase.Y < 1 {
				assert.Empty(m, "Matrix should not exist")
				return
			}
			require.Equal(tCase.X, len(m), "Matrix should have the requested amount of rows")

			i := 1
			for x := 0; x < tCase.X; x++ {
				require.Equalf(tCase.Y, len(m[x]), "Row %d should have the requested amount of columns", x)
				for y := 0; y < tCase.Y; y++ {
					require.Zerof(m[x][y], "(%d, %d) should be zero value", x, y)
					m[x][y] = i
					i++
				}
			}

			assert.Equal(tCase.X*tCase.Y, i-1, "Should have written to every position in matrix")
			i = 1
			for x := 0; x < len(m); x++ {
				for y := 0; y < len(m[x]); y++ {
					require.Equalf(i, m[x][y], "(%d, %d) should equal i", x, y)
					i++
				}
			}
		})
	}
}
