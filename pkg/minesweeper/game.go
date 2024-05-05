package minesweeper

import (
	"log/slog"

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

// Interface for playing the game
type Game interface {
	// Check a given field and recursevly reveal all neighboring fields that should be revield.
	// Returns the resulting new status of the game
	CheckField(p Pos) *Status
	// Recursive function to reveal all neighbouring fields that can be safely reveald.
	// Stops when a field has not exactly zero neighbouring mines
	RevealField(p Pos)
	// Check if the given position is out of bounds
	OutOfBounds(p Pos) bool
	// Returns the current status of the game. Only contains the knowledge a player should have.
	Status() *Status
	// Check if Game Over
	Lost() bool
	// Check if the game is won
	Won() bool
	// Reset the current game to be played again
	Replay()
	// Check if the game is a replay
	IsReplay() bool
	// Generate a save from the game
	ToSave() (*Save, error)
}

type LocalGame struct {
	Field      [][]Field
	Difficulty Difficulty

	// Keep these 2 exported for testing in other packages
	GameOver bool
	GameWon  bool

	replay bool
}

// Utility function to create empty game
func blankGame(d Difficulty) *LocalGame {
	return &LocalGame{
		Field:      utils.Make2D[Field](d.Row, d.Col),
		Difficulty: d,
		GameOver:   false,
		GameWon:    false,
	}
}

// Create a new game with the mines in the given positions
func newGame(d Difficulty, mines []Pos) *LocalGame {
	g := blankGame(d)

	for _, mine := range mines {
		g.Field[mine.X][mine.Y].Content = Mine
	}

	g.calculateFieldContent()

	return g
}

// Create a new game with mines seeded randomly in the map, with the exception of the given position.
func NewGameWithSafePos(d Difficulty, p Pos) *LocalGame {
	mines := CreateMines(d, []Pos{p})

	return newGame(d, mines)
}

// Create a new game with mines seeded randomly in the map, with the exception of a 3x3 area around the given position.
func NewGameWithSafeArea(d Difficulty, p Pos) *LocalGame {
	area := make([]Pos, 0, 9)
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			p := NewPos(p.X+x, p.Y+y)
			if !OutOfBounds(p, d) {
				area = append(area, p)
			}
		}
	}
	mines := CreateMines(d, area)

	return newGame(d, mines)
}

// Check a given field and recursevly reveal all neighboring fields that should be revield.
// Returns the resulting new status of the game
func (g *LocalGame) CheckField(p Pos) *Status {
	if g.Lost() || g.Won() {
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
func (g *LocalGame) RevealField(p Pos) {
	slog.Debug("Reveal field", slog.String("pos", p.String()), slog.String("content", g.Field[p.X][p.Y].Content.String()))

	g.Field[p.X][p.Y].Checked = true

	if g.Field[p.X][p.Y].Content != 0 {
		return
	}

	slog.Debug("Revealing fields neigbhours", slog.String("pos", p.String()))

	for m := -1; m < 2; m++ {
		for n := -1; n < 2; n++ {
			if m == 0 && n == 0 {
				continue
			}
			i := NewPos(p.X+m, p.Y+n)
			if g.OutOfBounds(i) {
				continue
			}
			if !g.Field[i.X][i.Y].Checked {
				g.RevealField(i)
			}
		}
	}
}

// Check if the given position is out of bounds
func (g *LocalGame) OutOfBounds(p Pos) bool {
	return OutOfBounds(p, g.Difficulty)
}

// Returns the current status of the game. Only contains the knowledge a player should have.
func (g *LocalGame) Status() *Status {
	d := g.Difficulty
	s := &Status{
		Field:    utils.Make2D[Field](d.Row, d.Col),
		GameOver: g.Lost(),
		GameWon:  g.Won(),
	}

	wasWon := g.Won()
	isWon := true

	g.walkField(func(x, y int) {
		s.Field[x][y].Checked = g.Field[x][y].Checked
		if g.Field[x][y].Checked || g.Lost() || g.Won() {
			s.Field[x][y].Content = g.Field[x][y].Content
		} else {
			s.Field[x][y].Content = Unknown
		}
		if !g.Field[x][y].Checked && g.Field[x][y].Content != Mine {
			isWon = false
		}
	})

	if !wasWon && isWon {
		g.GameWon, s.GameWon = isWon, isWon
		for x := 0; x < d.Row; x++ {
			copy(s.Field[x], g.Field[x])
		}
	}

	return s
}

// Check if Game Over
func (g *LocalGame) Lost() bool {
	return g.GameOver
}

// Check if the game is won
func (g *LocalGame) Won() bool {
	return g.GameWon
}

// Reset the current game to be played again
func (g *LocalGame) Replay() {
	g.replay = true
	g.GameOver = false
	g.GameWon = false

	g.walkField(func(x, y int) {
		g.Field[x][y].Checked = false
	})
}

// Check if the game is a replay
func (g *LocalGame) IsReplay() bool {
	return g.replay
}

// Generate a save from the game
func (g *LocalGame) ToSave() (*Save, error) {
	g.replay = true
	return NewSave(g)
}

// Walk through all fields of the game and call the given function
func (g *LocalGame) walkField(f func(x, y int)) {
	d := g.Difficulty
	for x := 0; x < d.Row; x++ {
		for y := 0; y < d.Col; y++ {
			f(x, y)
		}
	}
}

// Count the the number of mines in the neighboring fields
func (g *LocalGame) countNearbyMines(p Pos) int {
	c := 0
	for m := -1; m < 2; m++ {
		for n := -1; n < 2; n++ {
			if g.OutOfBounds(NewPos(p.X+m, p.Y+n)) {
				continue
			}
			if g.Field[p.X+m][p.Y+n].Content == Mine {
				c++
			}
		}
	}
	return c
}

// Get a list of all mines in the game
func (g *LocalGame) getMines() []Pos {
	mines := make([]Pos, 0, g.Difficulty.Mines)

	g.walkField(func(x, y int) {
		if g.Field[x][y].Content == Mine {
			mines = append(mines, NewPos(x, y))
		}
	})
	return mines
}

// Calculate all fields with the count of neighbouring mines
func (g *LocalGame) calculateFieldContent() {
	g.walkField(func(x, y int) {
		if g.Field[x][y].Content == Mine {
			return
		}

		g.Field[x][y].Content = FieldContent(g.countNearbyMines(NewPos(x, y)))
	})
}
