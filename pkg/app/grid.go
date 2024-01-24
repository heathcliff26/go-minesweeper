package app

import (
	"image/color"
	"log"

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

// Graphical display for a minesweeper game
type MinesweeperGrid struct {
	Tiles      [][]*Tile
	Difficulty minesweeper.Difficulty
	Game       *minesweeper.Game

	Timer     *Timer
	MineCount *Counter
	Reset     *Button
}

// Create a new grid suitable for the give difficulty
func NewMinesweeperGrid(d minesweeper.Difficulty) *MinesweeperGrid {
	tiles := utils.Make2D[*Tile](d.Size.Row, d.Size.Col)
	grid := &MinesweeperGrid{
		Tiles:      tiles,
		Difficulty: d,
		Timer:      NewTimer(),
		MineCount:  NewCounter(d.Mines),
	}
	grid.Reset = NewButton(ResetDefaultText, color.RGBA{}, grid.NewGame)
	grid.Reset.Label.TextSize = ResetTextSize

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
	reset := container.NewCenter(g.Reset)
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
func (g *MinesweeperGrid) TappedTile(x, y int) {
	pos := minesweeper.Pos{X: x, Y: y}
	if g.Game == nil {
		g.Game = minesweeper.NewGameWithSafePos(g.Difficulty, pos)
		g.Timer.Start()
	}

	log.Printf("Checking field (%d, %d)\n", pos.X, pos.Y)

	s := g.Game.CheckField(pos)
	if s.GameOver || s.GameWon {
		switch {
		case s.GameWon:
			log.Println("Win")
			g.Reset.SetText(ResetGameWonText)
		case s.GameOver:
			log.Println("Game Over")
			g.Reset.SetText(ResetGameOverText)
		}
		g.Timer.Stop()
	}

	log.Println("Checked field, updating tiles")

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
	log.Println("Finished Update")
}

// Return the number of rows in the grid
func (g *MinesweeperGrid) Row() int {
	return g.Difficulty.Size.Row
}

// Return the number of columns in the grid
func (g *MinesweeperGrid) Col() int {
	return g.Difficulty.Size.Col
}

// Start a new game
func (g *MinesweeperGrid) NewGame() {
	for x := 0; x < g.Row(); x++ {
		for y := 0; y < g.Col(); y++ {
			g.Tiles[x][y].Reset()
		}
	}
	g.Game = nil
	g.MineCount.SetCount(g.Difficulty.Mines)
	g.Timer.Reset()
	g.Reset.SetText(ResetDefaultText)
	g.Reset.Refresh()
}
