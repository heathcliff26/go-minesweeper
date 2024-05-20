package app

import (
	"image/color"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/heathcliff26/go-minesweeper/pkg/utils"
)

var (
	GridLabelColor = color.RGBA{240, 10, 20, alpha}
)

const (
	GridLabelSize float32 = 40

	ResetDefaultText          = "🙂"
	ResetGameOverText         = "☠"
	ResetGameWonText          = "😎"
	ResetTextSize     float32 = 40
)

const (
	GameAlgorithmSafePos = iota
	GameAlgorithmSafeArea
)

// Graphical display for a minesweeper game
type MinesweeperGrid struct {
	Tiles         [][]*Tile
	Difficulty    minesweeper.Difficulty
	Game          minesweeper.Game
	AssistedMode  bool
	GameAlgorithm int

	Timer       *Timer
	MineCount   *Counter
	ResetButton *Button
}

// Create a new grid suitable for the give difficulty
func NewMinesweeperGrid(d minesweeper.Difficulty, assistedMode bool) *MinesweeperGrid {
	tiles := utils.Make2D[*Tile](d.Row, d.Col)
	grid := &MinesweeperGrid{
		Tiles:         tiles,
		Difficulty:    d,
		AssistedMode:  assistedMode,
		GameAlgorithm: DEFAULT_GAME_ALGORITHM,
		Timer:         NewTimer(),
		MineCount:     NewCounter(d.Mines),
	}
	grid.ResetButton = NewButton(ResetDefaultText, color.RGBA{}, grid.NewGame)
	grid.ResetButton.Label.TextSize = ResetTextSize

	for x := 0; x < grid.Row(); x++ {
		for y := 0; y < grid.Col(); y++ {
			grid.Tiles[x][y] = NewTile(x, y, grid)
		}
	}

	return grid
}

// Get the graphical representation of the grid
func (g *MinesweeperGrid) GetCanvasObject() fyne.CanvasObject {
	mineCount := container.NewHBox(layout.NewSpacer(), container.NewCenter(newBorder(g.MineCount.Label)))
	reset := container.NewCenter(g.ResetButton)
	timer := container.NewHBox(container.NewCenter(newBorder(g.Timer.Label)), layout.NewSpacer())

	head := newBorder(container.NewGridWithColumns(3, mineCount, reset, timer))

	rows := make([]fyne.CanvasObject, len(g.Tiles))

	for x := 0; x < g.Row(); x++ {
		col := make([]fyne.CanvasObject, g.Col())
		for y := 0; y < g.Col(); y++ {
			col[y] = g.Tiles[x][y]
		}
		rows[x] = container.NewGridWithColumns(g.Col(), col...)
	}
	body := newBorder(container.NewGridWithRows(g.Row(), rows...))
	return container.NewVBox(head, body)
}

// Called by the child tiles to signal they have been tapped.
// Checks the given tile and then updates the display according to the new state.
// Starts a new game when no game is currently running.
func (g *MinesweeperGrid) TappedTile(pos minesweeper.Pos) {
	if g.Game == nil {
		switch g.GameAlgorithm {
		case GameAlgorithmSafePos:
			g.Game = minesweeper.NewGameWithSafePos(g.Difficulty, pos)
		case GameAlgorithmSafeArea:
			g.Game = minesweeper.NewGameWithSafeArea(g.Difficulty, pos)
		default:
			slog.Error("Unkown Algorithm for creating a new game", slog.Int("algorithm", g.GameAlgorithm))
			return
		}
	}
	if !g.Timer.Running() {
		g.Timer.Start()
	}

	slog.Debug("Checking field", slog.String("pos", pos.String()))

	s := g.Game.CheckField(pos)
	slog.Debug("Checked field, updating tiles")
	g.updateFromStatus(s)
}

// Update the grid from the given status
func (g *MinesweeperGrid) updateFromStatus(s *minesweeper.Status) {
	if s == nil {
		return
	}

	if s.GameOver || s.GameWon {
		switch {
		case s.GameWon:
			slog.Info("Win")
			g.ResetButton.SetText(ResetGameWonText)
		case s.GameOver:
			slog.Info("Game Over")
			g.ResetButton.SetText(ResetGameOverText)
		}
		g.Timer.Stop()
	} else if g.AssistedMode {
		slog.Debug("Creating Markers for Assisted Mode")
		for _, p := range s.ObviousMines() {
			t := g.Tiles[p.X][p.Y]
			t.Marker = HelpMarkingMine
			t.UpdateContent()
		}
		for _, p := range s.ObviousSafePos() {
			t := g.Tiles[p.X][p.Y]
			t.Marker = HelpMarkingSafe
			t.UpdateContent()
		}
	}

	for x := 0; x < g.Row(); x++ {
		for y := 0; y < g.Col(); y++ {
			t := g.Tiles[x][y]
			if s.Field[x][y].Content == minesweeper.Unknown {
				continue
			}
			t.Field = &s.Field[x][y]
			t.UpdateContent()
		}
	}
	slog.Debug("Finished Update")
}

// Return the number of rows in the grid
func (g *MinesweeperGrid) Row() int {
	return g.Difficulty.Row
}

// Return the number of columns in the grid
func (g *MinesweeperGrid) Col() int {
	return g.Difficulty.Col
}

// Start a new game
func (g *MinesweeperGrid) NewGame() {
	slog.Info("Preparing for new game")
	g.Game = nil
	g.Reset()
}

// Replay the current game
func (g *MinesweeperGrid) Replay() {
	slog.Info("Preparing for replay of current game")
	if g.Game != nil {
		g.Game.Replay()
	}
	g.Reset()
}

// Reset Grid
func (g *MinesweeperGrid) Reset() {
	for x := 0; x < g.Row(); x++ {
		for y := 0; y < g.Col(); y++ {
			g.Tiles[x][y].Reset()
		}
	}
	g.MineCount.SetCount(g.Difficulty.Mines)
	g.Timer.Reset()
	g.ResetButton.SetText(ResetDefaultText)
	g.ResetButton.Refresh()
	slog.Debug("Reset grid")
}

// Check if the given position is out of bounds.
// Calls Game.OutOfBounds(Pos)
func (g *MinesweeperGrid) OutOfBounds(p minesweeper.Pos) bool {
	return g.Game.OutOfBounds(p)
}
