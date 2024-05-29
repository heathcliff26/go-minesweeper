//go:build race

package app

import "testing"

func TestTileFlagRace(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	tile := g.Tiles[0][0]
	tile.CreateRenderer()

	go tile.Flag(true)
	go tile.Flagged()
}
