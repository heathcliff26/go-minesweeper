package minesweeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutOfBounds(t *testing.T) {
	g := NewGameWithSafePos(difficulties[DifficultyExpert], Pos{2, 2})
	s := g.Difficulty.Size

	assert := assert.New(t)

	for x := 0; x < s.Row; x++ {
		for y := 0; y < s.Col; y++ {
			p := Pos{x, y}
			assert.Falsef(g.outOfBounds(p), "%v should be within the field", p)

			p.X = -1
			assert.Truef(g.outOfBounds(p), "%v should be out of bounds", p)
			p.X = s.Row
			assert.Truef(g.outOfBounds(p), "%v should be out of bounds", p)
			p.X = x

			p.Y = -1
			assert.Truef(g.outOfBounds(p), "%v should be out of bounds", p)
			p.Y = s.Col
			assert.Truef(g.outOfBounds(p), "%v should be out of bounds", p)
		}
	}
}
