package minesweeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutOfBounds(t *testing.T) {
	g := NewGameWithSafePos(difficulties[DifficultyExpert], Pos{2, 2})
	d := g.Difficulty

	assert := assert.New(t)

	for x := 0; x < d.Row; x++ {
		for y := 0; y < d.Col; y++ {
			p := Pos{x, y}
			assert.Falsef(g.outOfBounds(p), "%v should be within the field", p)

			p.X = -1
			assert.Truef(g.outOfBounds(p), "%v should be out of bounds", p)
			p.X = d.Row
			assert.Truef(g.outOfBounds(p), "%v should be out of bounds", p)
			p.X = x

			p.Y = -1
			assert.Truef(g.outOfBounds(p), "%v should be out of bounds", p)
			p.Y = d.Col
			assert.Truef(g.outOfBounds(p), "%v should be out of bounds", p)
		}
	}
}
