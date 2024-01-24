package minesweeper

import (
	"log"

	"github.com/heathcliff26/go-minesweeper/pkg/utils"
)

// Represents the content of a field.
// Can be a mine, unknown or the number of mines in the neighboring fields
type FieldContent int

const (
	Mine    FieldContent = -1
	Unknown FieldContent = -2
)

// Represents a single field in a minefield
type Field struct {
	Checked bool
	Content FieldContent
}

// Status contains the current state of the game known to the player.
// As such it will always be a copy and needs to be it's own type, despite the
// overlapping similarities.
// It does not support any of the functions that Game does.
// It is safe to write to Status, as it is merely a copy.
type Status struct {
	Field    [][]Field
	GameOver bool
	GameWon  bool
}

// Struct for holding the game
type Game struct {
	Field      [][]Field
	Difficulty Difficulty
	GameOver   bool
	GameWon    bool
}

// Create a new game with mines seeded randomly in the map, with the exception of the given position.
func NewGameWithSafePos(d Difficulty, p Pos) *Game {
	g := &Game{
		Field:      utils.Make2D[Field](d.Row, d.Col),
		Difficulty: d,
		GameOver:   false,
		GameWon:    false,
	}

	mines := CreateMines(d, p)
	for _, mine := range mines {
		g.Field[mine.X][mine.Y].Content = Mine
	}

	for x := 0; x < d.Row; x++ {
		for y := 0; y < d.Col; y++ {
			if g.Field[x][y].Content == Mine {
				continue
			}
			c := 0
			for m := -1; m < 2; m++ {
				for n := -1; n < 2; n++ {
					if x+m < 0 || x+m >= d.Row || y+n < 0 || y+n >= d.Row {
						continue
					}
					if g.Field[x+m][y+n].Content == Mine {
						c++
					}
				}
			}
			g.Field[x][y].Content = FieldContent(c)
		}
	}

	return g
}

// Check a given field and recursevly reveal all neighboring fields that should be revield.
// Returns the resulting new status of the game
func (g *Game) CheckField(p Pos) *Status {
	if g.GameOver || g.GameWon {
		return g.Status()
	}

	g.Field[p.X][p.Y].Checked = true

	if g.Field[p.X][p.Y].Content == Mine {
		g.GameOver = true
		return g.Status()
	}

	g.RevealField(p)

	return g.Status()
}

// Recursive function to reveal all neighbouring fields that can be safely reveald.
// Stops when a field has not exactly zero neighbouring mines
func (g *Game) RevealField(p Pos) {
	log.Printf("Reveal tile (%d, %d), content: %d\n", p.X, p.Y, g.Field[p.X][p.Y].Content)

	g.Field[p.X][p.Y].Checked = true

	if g.Field[p.X][p.Y].Content != 0 {
		return
	}

	log.Printf("Revealing neigbhours of (%d, %d)\n", p.X, p.Y)

	for m := -1; m < 2; m++ {
		for n := -1; n < 2; n++ {
			if m == 0 && n == 0 {
				continue
			}
			i := Pos{p.X + m, p.Y + n}
			if g.outOfBounds(i) {
				continue
			}
			if !g.Field[i.X][i.Y].Checked {
				g.RevealField(i)
			}
		}
	}
}

// Check if the given position is out of bounds
func (g *Game) outOfBounds(p Pos) bool {
	d := g.Difficulty
	return p.X < 0 || p.X > d.Row-1 || p.Y < 0 || p.Y > d.Col-1
}

// Returns the current status of the game. Only contains the knowledge a player should have.
func (g *Game) Status() *Status {
	d := g.Difficulty
	s := &Status{
		Field:    utils.Make2D[Field](d.Row, d.Col),
		GameOver: g.GameOver,
	}

	wasWon := g.GameWon
	isWon := true

	for x := 0; x < d.Row; x++ {
		for y := 0; y < d.Col; y++ {
			s.Field[x][y].Checked = g.Field[x][y].Checked
			if g.Field[x][y].Checked || g.GameOver || wasWon {
				s.Field[x][y].Content = g.Field[x][y].Content
			} else {
				s.Field[x][y].Content = Unknown
			}
			if !g.Field[x][y].Checked && g.Field[x][y].Content != Mine {
				isWon = false
			}
		}
	}

	g.GameWon, s.GameWon = isWon, isWon

	if g.GameWon != wasWon {
		for x := 0; x < d.Row; x++ {
			copy(s.Field[x], g.Field[x])
		}
	}

	return s
}
