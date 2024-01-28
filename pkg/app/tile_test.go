package app

import (
	"strconv"
	"testing"

	"github.com/heathcliff26/go-minesweeper/assets"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/stretchr/testify/assert"
)

func TestNewTile(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY)
	tile := g.Tiles[1][2]

	t.Run("New", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(1, tile.x)
		assert.Equal(2, tile.y)

		assert.Equal(g, tile.grid)

		f := &minesweeper.Field{
			Checked: false,
			Content: minesweeper.Unknown,
		}
		assert.Equal(f, tile.Field)

		assert.False(tile.Flagged)
	})
	t.Run("CreateRender", func(t *testing.T) {
		tile.CreateRenderer()

		assert := assert.New(t)

		assert.Equal(TileDefaultColor, tile.background.FillColor)
		assert.Equal(TileSize, tile.background.MinSize())

		assert.Equal("", tile.label.Text)
		assert.True(tile.label.TextStyle.Bold)
		assert.Equal(TileTextSize, tile.label.TextSize)
		assert.True(tile.label.Hidden)

		assert.Equal(TileSize, tile.icon.Size())
		assert.True(tile.icon.Hidden)
	})
}

func TestTileTapped(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}
	tile := g.Tiles[0][0]

	assert := assert.New(t)

	tile.Flagged = true
	tile.Tapped(nil)

	assert.Nil(tile.grid.Game, "Flagged tiles should not trigger checks")
	assert.False(tile.Field.Checked, "Flagged tiles should not trigger checks")

	tile.Flagged = false
	tile.Tapped(nil)

	assert.NotNil(tile.grid.Game, "Game should be started")
	assert.True(tile.Field.Checked, "Tile should be checked")
}

func TestTileTappedSecondary(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY)
	tile := g.Tiles[1][2]
	tile.CreateRenderer()

	assert := assert.New(t)

	count := tile.grid.MineCount.Count

	tile.Field.Checked = true
	tile.TappedSecondary(nil)

	assert.Equal(count, tile.grid.MineCount.Count, "MineCount should not have changed")
	assert.False(tile.Flagged, "Tile should not be flagged")

	tile.Field.Checked = false
	tile.TappedSecondary(nil)

	assert.Equal(count-1, tile.grid.MineCount.Count, "MineCount should be decreased")
	assert.True(tile.Flagged, "Tile should be flagged")

	tile.TappedSecondary(nil)

	assert.Equal(count, tile.grid.MineCount.Count, "MineCount should be back to original value")
	assert.False(tile.Flagged, "Tile should not be flagged")
}

func TestTileUpdateContent(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY)
	tile := g.Tiles[1][2]
	tile.CreateRenderer()

	tile.UpdateContent()

	t.Run("Default", func(t *testing.T) {
		assert := assert.New(t)

		assert.True(tile.icon.Hidden)
		assert.True(tile.label.Hidden)
		assert.Equal(TileDefaultColor, tile.background.FillColor)
	})
	t.Run("Checked", func(t *testing.T) {
		tile.Field.Checked = true
		t.Cleanup(func() {
			tile.Field.Checked = false
		})
		tile.UpdateContent()

		assert.Equal(t, TileBackgroundColor, tile.background.FillColor)
	})
	t.Run("Flagged", func(t *testing.T) {
		tile.Flagged = true
		t.Cleanup(func() {
			tile.Flagged = false
		})
		tile.UpdateContent()

		assert := assert.New(t)

		assert.False(tile.icon.Hidden)
		assert.Equal(assets.ResourceFlagPng, tile.icon.Resource)
	})
	t.Run("Mine", func(t *testing.T) {
		tile.Field.Content = minesweeper.Mine
		t.Cleanup(func() {
			tile.Field.Content = minesweeper.Unknown
		})
		tile.UpdateContent()

		assert := assert.New(t)

		assert.False(tile.icon.Hidden)
		assert.Equal(assets.ResourceMinePng, tile.icon.Resource)
		assert.Equal(TileDefaultColor, tile.background.FillColor)

		tile.Field.Checked = true
		t.Cleanup(func() {
			tile.Field.Checked = false
		})
		tile.UpdateContent()

		assert.False(tile.icon.Hidden)
		assert.Equal(assets.ResourceMinePng, tile.icon.Resource)
		assert.Equal(TileExplodedColor, tile.background.FillColor)
	})
	t.Run("Numbers", func(t *testing.T) {
		t.Cleanup(func() {
			tile.Field.Content = minesweeper.Unknown
		})

		assert := assert.New(t)

		tile.Field.Content = minesweeper.FieldContent(0)
		tile.UpdateContent()

		assert.True(tile.label.Hidden)

		for i := 1; i < 9; i++ {
			tile.Field.Content = minesweeper.FieldContent(i)
			tile.UpdateContent()

			assert.False(tile.label.Hidden)
			assert.Equal(strconv.Itoa(i), tile.label.Text)
			assert.Equal(TileTextColor[i], tile.label.Color)
		}

		tile.Field.Content = minesweeper.FieldContent(9)
		tile.UpdateContent()

		assert.True(tile.label.Hidden)
	})
}

func TestTileReset(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY)
	tile := g.Tiles[1][2]
	tile.CreateRenderer()

	tile.Flagged = true
	tile.Field.Checked = true
	tile.Field.Content = minesweeper.Mine

	tile.UpdateContent()
	tile.Reset()

	assert := assert.New(t)

	assert.False(tile.Flagged)
	assert.False(tile.Field.Checked)
	assert.Equal(minesweeper.Unknown, tile.Field.Content)

	assert.Equal(TileDefaultColor, tile.background.FillColor)
	assert.Equal(TileSize, tile.background.MinSize())

	assert.Equal("", tile.label.Text)
	assert.True(tile.label.TextStyle.Bold)
	assert.Equal(TileTextSize, tile.label.TextSize)
	assert.True(tile.label.Hidden)

	assert.True(tile.icon.Hidden)
}

func TestTileUntappable(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}
	tile := g.Tiles[1][2]

	assert := assert.New(t)

	assert.False(tile.untappable())

	tile.Field.Checked = true
	assert.True(tile.untappable())
	tile.Field.Checked = false

	tile.Tapped(nil)
	tile.Field.Checked = false
	assert.False(tile.untappable())

	tile.grid.Game.GameOver = true
	assert.True(tile.untappable())
	tile.grid.Game.GameOver = false

	tile.grid.Game.GameWon = true
	assert.True(tile.untappable())
	tile.grid.Game.GameWon = false
}
