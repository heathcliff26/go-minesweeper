package minesweeper

import (
	"fmt"
	"math/rand"
	"slices"
)

type Pos struct {
	X, Y int
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func CreateMines(d Difficulty, p Pos) []Pos {
	mines := make([]Pos, 0, d.Mines+1)

	mines = append(mines, p)
	for i := 0; i < d.Mines; i++ {
		p = RandomPos(d.Size.Row, d.Size.Col)
		for slices.Contains(mines, p) {
			p = RandomPos(d.Size.Row, d.Size.Col)
		}
		mines = append(mines, p)
	}
	return mines[1:]
}

func RandomPos(maxX, maxY int) Pos {
	return Pos{
		X: rand.Intn(maxX),
		Y: rand.Intn(maxY),
	}
}
