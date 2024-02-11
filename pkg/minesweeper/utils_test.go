package minesweeper

import (
	"strconv"
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

type tCaseFieldContentString struct {
	Name, Result string
	FC           FieldContent
}

func TestFieldContentString(t *testing.T) {
	tMatrix := []tCaseFieldContentString{
		{
			FC:     Mine,
			Result: "Mine",
		},
		{
			FC:     Unknown,
			Result: "Unknown",
		},
		{
			Name: "Invalid1",
			FC:   9,
		},
		{
			Name: "Invalid2",
			FC:   -3,
		},
	}
	for i := range 9 {
		tMatrix = append(
			tMatrix,
			tCaseFieldContentString{
				FC:     FieldContent(i),
				Result: strconv.Itoa(i),
			},
		)
	}

	for _, tCase := range tMatrix {
		name := tCase.Name
		if name == "" {
			name = tCase.Result
		}

		t.Run(name, func(t *testing.T) {
			res := tCase.FC.String()

			assert := assert.New(t)
			if tCase.Name == "" {
				assert.Equal(tCase.Result, res)
			} else {
				assert.Contains(res, "not a valid FieldContent")
			}
		})
	}
}
