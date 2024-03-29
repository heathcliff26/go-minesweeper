package minesweeper

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
)

// Represent a position in the minefield
type Pos struct {
	X, Y int
}

// Create a new Position from the given coordinates
func NewPos(x, y int) Pos {
	return Pos{x, y}
}

// Returns a random position inside the provided limits
func RandomPos(maxX, maxY int) Pos {
	return NewPos(rand.IntN(maxX), rand.IntN(maxY))
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

// Convert FieldContent to string for logging
func (fc FieldContent) String() string {
	switch {
	case fc == Mine:
		return "Mine"
	case fc == Unknown:
		return "Unknown"
	case fc >= 0 && fc < 9:
		return strconv.Itoa(int(fc))
	default:
		return fmt.Sprintf("%d is not a valid FieldContent", fc)
	}
}
