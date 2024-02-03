package app

import (
	"testing"

	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/stretchr/testify/assert"
)

func TestNewGrid(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY)

	assert := assert.New(t)

	assert.Equal(g.Row(), len(g.Tiles))
	assert.Equal(g.Col(), len(g.Tiles[0]))
	assert.Equal(DEFAULT_DIFFICULTY, g.Difficulty)
	assert.Nil(g.Game)
	assert.NotNil(g.Timer)
	assert.NotNil(g.MineCount)
	assert.NotNil(g.Reset)
}

func TestTappedTile(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}

	p := minesweeper.NewPos(0, 0)

	assert := assert.New(t)

	g.TappedTile(p)
	assert.NotNil(g.Game)
	assert.True(g.Tiles[0][0].Field.Checked)

	g.Game.GameWon = true
	g.TappedTile(p)
	assert.Equal(ResetGameWonText, g.Reset.Label.Text)
	g.Game.GameWon = false

	g.Game.GameOver = true
	g.TappedTile(p)
	assert.Equal(ResetGameOverText, g.Reset.Label.Text)
}

func TestNewGame(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}

	g.TappedTile(minesweeper.NewPos(0, 0))
	g.NewGame()

	assert := assert.New(t)

	for x := 0; x < g.Row(); x++ {
		for y := 0; y < g.Col(); y++ {
			assert.Falsef(g.Tiles[x][y].Field.Checked, "(%d, %d) All fields should be reset", x, y)
		}
	}

	assert.Nil(g.Game)
	assert.Equal(g.Difficulty.Mines, g.MineCount.Count)
	assert.False(g.Timer.running)
	assert.Equal(ResetDefaultText, g.Reset.Label.Text)
}
