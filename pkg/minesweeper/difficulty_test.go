package minesweeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDifficulties(t *testing.T) {
	res := Difficulties()

	assert := assert.New(t)

	assert.Equal(difficulties, res)
	assert.NotSame(&difficulties[0], &res[0])
}

func TestNewCustomDifficulty(t *testing.T) {
	tMatrix := []struct {
		Difficulty Difficulty
		Error      error
	}{
		{
			Difficulty: Difficulty{
				Name:  "RowTooSmall",
				Row:   DifficultyRowColMin - 1,
				Col:   12,
				Mines: DifficultyMineMin,
			},
			Error: NewErrDifficultyDimension(DifficultyRowColMin-1, 12),
		},
		{
			Difficulty: Difficulty{
				Name:  "RowTooBig",
				Row:   DifficultyRowColMax + 1,
				Col:   12,
				Mines: DifficultyMineMin,
			},
			Error: NewErrDifficultyDimension(DifficultyRowColMax+1, 12),
		},
		{
			Difficulty: Difficulty{
				Name:  "ColTooSmall",
				Row:   12,
				Col:   DifficultyRowColMin - 1,
				Mines: DifficultyMineMin,
			},
			Error: NewErrDifficultyDimension(12, DifficultyRowColMin-1),
		},
		{
			Difficulty: Difficulty{
				Name:  "ColTooBig",
				Row:   12,
				Col:   DifficultyRowColMax + 1,
				Mines: DifficultyMineMin,
			},
			Error: NewErrDifficultyDimension(12, DifficultyRowColMax+1),
		},
		{
			Difficulty: Difficulty{
				Name:  "NotEnoughMines",
				Row:   12,
				Col:   12,
				Mines: DifficultyMineMin - 1,
			},
			Error: NewErrDifficultyMineCount(DifficultyMineMin - 1),
		},
		{
			Difficulty: Difficulty{
				Name:  "TooManyMines",
				Row:   12,
				Col:   12,
				Mines: 116,
			},
			Error: NewErrDifficultyMineCount(116),
		},
		{
			Difficulty: Difficulty{
				Name:  "LowerLimit",
				Row:   DifficultyRowColMin,
				Col:   DifficultyRowColMin,
				Mines: DifficultyMineMin,
			},
			Error: nil,
		},
		{
			Difficulty: Difficulty{
				Name:  "UpperLimit",
				Row:   DifficultyRowColMax,
				Col:   DifficultyRowColMax,
				Mines: 7840,
			},
			Error: nil,
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Difficulty.Name, func(t *testing.T) {
			d, err := NewCustomDifficulty(tCase.Difficulty.Mines, tCase.Difficulty.Row, tCase.Difficulty.Col)

			assert := assert.New(t)

			if tCase.Error == nil {
				assert.NoError(err)
				res := tCase.Difficulty
				res.Name = "Custom"
				assert.Equal(res, d)
			} else {
				assert.Empty(d)
				assert.Equal(tCase.Error, err)
			}
		})
	}
}
