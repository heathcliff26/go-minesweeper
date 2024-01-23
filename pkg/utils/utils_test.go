package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

			if tCase.X < 1 || tCase.Y < 1 {
				assert.Empty(m, "Matrix should not exist")
				return
			}
			if !assert.Equal(tCase.X, len(m), "Matrix should have the requested amount of rows") {
				t.FailNow()
			}

			i := 1
			for x := 0; x < tCase.X; x++ {
				if !assert.Equalf(tCase.Y, len(m[x]), "Row %d should have the requested amount of columns", x) {
					t.FailNow()
				}
				for y := 0; y < tCase.Y; y++ {
					if !assert.Zerof(m[x][y], "(%d, %d) should be zero value", x, y) {
						t.FailNow()
					}
					m[x][y] = i
					i++
				}
			}

			assert.Equal(tCase.X*tCase.Y, i-1, "Should have written to every position in matrix")
			i = 1
			for x := 0; x < len(m); x++ {
				for y := 0; y < len(m[x]); y++ {
					if !assert.Equalf(i, m[x][y], "(%d, %d) should equal i", x, y) {
						t.FailNow()
					}
					i++
				}
			}
		})
	}
}
