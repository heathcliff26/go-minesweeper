//go:build race

package app

import (
	"testing"

	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
)

func TestTileFlagRace(t *testing.T) {
	g := NewMinesweeperGrid(minesweeper.Difficulties()[DEFAULT_DIFFICULTY], false)
	tile := g.Tiles[0][0]
	tile.CreateRenderer()

	go tile.Flag(true)
	go tile.Flagged()
}
