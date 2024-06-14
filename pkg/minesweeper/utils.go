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
// Does not create a mine on the given positions.
func CreateMines(d Difficulty, safe []Pos) []Pos {
	var p Pos
	mines := make([]Pos, 0, d.Mines+len(safe))

	mines = append(mines, safe...)
	for i := 0; i < d.Mines; i++ {
		p = RandomPos(d.Row, d.Col)
		for slices.Contains(mines, p) {
			p = RandomPos(d.Row, d.Col)
		}
		mines = append(mines, p)
	}
	return mines[len(safe):]
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

// Check if a position is out of bounds on the given difficulty
func OutOfBounds(p Pos, d Difficulty) bool {
	return p.X < 0 || p.X > d.Row-1 || p.Y < 0 || p.Y > d.Col-1
}

// Walk through all fields in the given dimension and call the given function
func walkField(f func(x, y int), limitX, limitY int) {
	for x := 0; x < limitX; x++ {
		for y := 0; y < limitY; y++ {
			f(x, y)
		}
	}
}

// Autosolve the given Game
// Returns if the game can be autosolved
func autosolve(g *LocalGame, startPos Pos) bool {
	g.CheckField(startPos)

	for safePos := g.Status().ObviousSafePos(); len(safePos) > 0 && !g.Won(); safePos = g.Status().ObviousSafePos() {
		for _, p := range safePos {
			g.CheckField(p)
		}
	}

	return g.Won()
}

// Create an array of positions with 3x3, centered around the given position
func areaAroundPos(d Difficulty, p Pos) []Pos {
	area := make([]Pos, 0, 9)
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			p := NewPos(p.X+x, p.Y+y)
			if !OutOfBounds(p, d) {
				area = append(area, p)
			}
		}
	}
	return area
}
