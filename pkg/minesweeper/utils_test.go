package minesweeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMines(t *testing.T) {
	tMatrix := Difficulties()
	customDifficulty := Difficulty{
		Name:  "Custom",
		Mines: 1000,
		Size:  GridSize{1000, 1000},
	}
	tMatrix = append(tMatrix, customDifficulty)

	for _, d := range tMatrix {
		t.Run(d.Name, func(t *testing.T) {
			p := RandomPos(d.Size.Row, d.Size.Col)
			mines := CreateMines(d, p)

			assert := assert.New(t)

			assert.Equal(d.Mines, len(mines))
			assert.NotContains(mines, p)
		})
	}
}
