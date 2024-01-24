package minesweeper

import (
	"fmt"
	"math/rand"
	"slices"
)

// Represent a position in the minefield
type Pos struct {
	X, Y int
}

// Returns a string representation of the position
func (p Pos) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

// Randomly create mines for the given difficulty.
// Does not create a mine on the given position.
func CreateMines(d Difficulty, p Pos) []Pos {
	mines := make([]Pos, 0, d.Mines+1)

	mines = append(mines, p)
	for i := 0; i < d.Mines; i++ {
		p = RandomPos(d.Row, d.Col)
		for slices.Contains(mines, p) {
			p = RandomPos(d.Row, d.Col)
		}
		mines = append(mines, p)
	}
	return mines[1:]
}

// Returns a random position inside the provided limits
func RandomPos(maxX, maxY int) Pos {
	return Pos{
		X: rand.Intn(maxX),
		Y: rand.Intn(maxY),
	}
}
