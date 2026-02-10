package minesweeper

import (
	"fmt"
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

// Interface for playing the game
type Game interface {
	// Check a given field and recursevly reveal all neighboring fields that should be revield.
	// Returns the resulting new status of the game
	CheckField(p Pos) (*Status, bool)
	// Check if the given position is out of bounds
	OutOfBounds(p Pos) bool
	// Returns the current status of the game. Only contains the knowledge a player should have.
	Status() *Status
	// Return the difficulty of the game
	Difficulty() Difficulty
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

// Ensure that compiler throws error if LocalGame does not implement Game interface
var _ = Game(&LocalGame{})

type LocalGame struct {
	Field      [][]Field
	difficulty Difficulty

	// Keep these 2 exported for testing in other packages
	GameOver bool
	GameWon  bool

	// Cache the last status
	status *Status

	replay bool
}

// Utility function to create empty game
func blankGame(d Difficulty) *LocalGame {
	return &LocalGame{
		Field:      utils.Make2D[Field](d.Row, d.Col),
		difficulty: d,
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
	mines := CreateMines(d, areaAroundPos(d, p))

	return newGame(d, mines)
}

// Create a new game that is solvable without random guesses.
func NewGameSolvable(d Difficulty, p Pos) *LocalGame {
	g, err := NewGameSolvableWithIterations(d, p, 10000)
	if err == nil {
		slog.Info("Created solvable minesweeper game")
	} else {
		slog.Error("Failed to create a solvable minesweeper game", "err", err)
	}
	return g
}

// Create a new game that is solvable without random guesses, within a specified number of iterations.
// Returns a solvable game, or a game that needs random guesses and an error indicating that.
func NewGameSolvableWithIterations(d Difficulty, p Pos, maxIterations int) (*LocalGame, error) {
	if maxIterations < 1 {
		return NewGameWithSafePos(d, p), fmt.Errorf("maxIterations must be greater than zero")
	}
	area := areaAroundPos(d, p)

	var mines []Pos
	var success bool
	var err error
	for i := 0; i < maxIterations; i++ {
		mines = CreateMines(d, area)
		g := newGame(d, mines)
		s := NewSolver(g)

		if success = s.Autosolve(p); success {
			break
		}
	}
	if !success {
		err = NewErrCreateUnsolvableGame(maxIterations)
	}

	return newGame(d, mines), err
}

// Check a given field and recursevly reveal all neighboring fields that should be revield.
// Returns the resulting new status of the game and a boolean indicating if there where changes.
func (g *LocalGame) CheckField(p Pos) (*Status, bool) {
	if g.Lost() || g.Won() || g.Field[p.X][p.Y].Checked {
		return g.Status(), false
	}

	g.Field[p.X][p.Y].Checked = true

	if g.Field[p.X][p.Y].Content == Mine {
		g.GameOver = true
		return g.UpdateStatus(), true
	}

	g.revealField(p)

	return g.UpdateStatus(), true
}

// Recursive function to reveal all neighbouring fields that can be safely reveald.
// Stops when a field has not exactly zero neighbouring mines
func (g *LocalGame) revealField(p Pos) {
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
				g.revealField(i)
			}
		}
	}
}

// Check if the given position is out of bounds
func (g *LocalGame) OutOfBounds(p Pos) bool {
	if g == nil {
		return true
	}
	return OutOfBounds(p, g.Difficulty())
}

// Returns the current status of the game. Only contains the knowledge a player should have.
func (g *LocalGame) Status() *Status {
	return g.status
}

// Return the difficulty of the game
func (g *LocalGame) Difficulty() Difficulty {
	return g.difficulty
}

// Update the status from the current state of the game.
// Returns status for convenience.
func (g *LocalGame) UpdateStatus() *Status {
	if g.status == nil {
		g.status = &Status{
			Field:    utils.Make2D[Field](g.difficulty.Row, g.difficulty.Col),
			gameOver: g.Lost(),
			gameWon:  g.Won(),
		}
	}

	wasWon := g.Won()
	isWon := true

	g.walkField(func(x, y int) {
		g.status.Field[x][y].Checked = g.Field[x][y].Checked
		if g.Field[x][y].Checked || g.Lost() || g.Won() {
			g.status.Field[x][y].Content = g.Field[x][y].Content
		} else {
			g.status.Field[x][y].Content = Unknown
		}
		if !g.Field[x][y].Checked && g.Field[x][y].Content != Mine {
			isWon = false
		}
	})

	if !wasWon && isWon {
		g.GameWon = isWon
		for x := 0; x < g.difficulty.Row; x++ {
			copy(g.status.Field[x], g.Field[x])
		}
	}

	g.status.gameOver = g.GameOver
	g.status.gameWon = g.GameWon

	return g.status
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
	g.status = nil

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
	walkField(f, g.difficulty.Row, g.difficulty.Col)
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
	mines := make([]Pos, 0, g.difficulty.Mines)

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
