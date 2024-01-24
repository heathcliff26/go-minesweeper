package minesweeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositionToString(t *testing.T) {
	s := Pos{1, 1}.String()

	assert.Equal(t, "(1, 1)", s)
}

func TestCreateMines(t *testing.T) {
	tMatrix := Difficulties()
	customDifficulty := Difficulty{
		Name:  "Custom",
		Mines: 1000,
		Row:   1000,
		Col:   1000,
	}
	tMatrix = append(tMatrix, customDifficulty)

	for _, d := range tMatrix {
		t.Run(d.Name, func(t *testing.T) {
			p := RandomPos(d.Row, d.Col)
			mines := CreateMines(d, p)

			assert := assert.New(t)

			assert.Equal(d.Mines, len(mines))
			assert.NotContains(mines, p)
		})
	}
}
