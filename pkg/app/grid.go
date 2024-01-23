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

type MinesweeperGrid struct {
	Tiles      [][]*Tile
	Difficulty minesweeper.Difficulty
	Game       *minesweeper.Game

	Timer     *Timer
	MineCount *MineCount
	Reset     *Button
}

func NewMinesweeperGrid(d minesweeper.Difficulty) *MinesweeperGrid {
	tiles := utils.Make2D[*Tile](d.Size.Row, d.Size.Col)
	grid := &MinesweeperGrid{
		Tiles:      tiles,
		Difficulty: d,
		Timer:      NewTimer(),
		MineCount:  NewMineCount(d.Mines),
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

func (g *MinesweeperGrid) GetCanvasObject() fyne.CanvasObject {
	mineCount := newBorder(g.MineCount.Label)
	timer := newBorder(g.Timer.Label)
	head := newBorder(container.NewHBox(layout.NewSpacer(), mineCount, g.Reset, timer, layout.NewSpacer()))

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

func (g *MinesweeperGrid) Row() int {
	return g.Difficulty.Size.Row
}

func (g *MinesweeperGrid) Col() int {
	return g.Difficulty.Size.Col
}

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
