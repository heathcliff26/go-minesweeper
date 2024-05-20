package app

import (
	"strconv"
	"testing"

	"github.com/heathcliff26/go-minesweeper/assets"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/stretchr/testify/assert"
)

func TestNewTile(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	tile := g.Tiles[1][2]

	t.Run("New", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(minesweeper.NewPos(1, 2), tile.Pos)

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
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
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
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
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

func TestDoubleTapped(t *testing.T) {
	tMatrix := []struct {
		Name    string
		Flagged bool
	}{
		{"Unchecked", false},
		{"Flagged", true},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
			for _, row := range g.Tiles {
				for _, tile := range row {
					tile.CreateRenderer()
				}
			}
			tile := g.Tiles[0][0]
			g.Game = minesweeper.NewGameWithSafePos(g.Difficulty, tile.Pos)

			tile.Flagged = tCase.Flagged
			t.Cleanup(func() {
				tile.Flagged = false
			})
			tile.DoubleTapped(nil)

			for m := -1; m < 2; m++ {
				for n := -1; n < 2; n++ {
					p := tile.Pos
					p.X += m
					p.Y += n
					if g.OutOfBounds(p) {
						continue
					}

					assert.Falsef(t, g.Tiles[p.X][p.Y].Field.Checked, p.String())
				}
			}
		})
	}

	t.Run("Checked", func(t *testing.T) {
		assert := assert.New(t)

		save, err := minesweeper.LoadSave("testdata/double-tapped_checked.sav")
		if !assert.Nil(err, "Should load savegame") {
			t.FailNow()
		}
		game := save.Game()

		g := NewMinesweeperGrid(game.Difficulty, false)
		for _, row := range g.Tiles {
			for _, tile := range row {
				tile.CreateRenderer()
			}
		}
		tile := g.Tiles[1][1]
		g.Game = game
		tile.Tapped(nil)
		g.Tiles[0][0].Flagged = true

		tile.DoubleTapped(nil)

		assert.False(g.Game.Status().GameOver, "Game should not be lost")
		assert.True(g.Game.Status().GameWon, "Game should be won")
		assert.False(g.Tiles[0][0].Field.Checked, "Flagged field should not be checked")

		for x := 0; x < 2; x++ {
			for y := 0; y < 2; y++ {
				if x == 0 && y == 0 {
					continue
				}
				assert.Truef(g.Tiles[x][y].Field.Checked, "Field should be checked, tile=(%d, %d)", x, y)
			}
		}
	})
}

func TestTileUpdateContent(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
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
		assert.Equal(TileDefaultColor, tile.background.FillColor)
	})
	t.Run("FlaggedSuccess", func(t *testing.T) {
		tile.Flagged = true
		tile.Field.Content = minesweeper.Mine
		t.Cleanup(func() {
			tile.Flagged = false
			tile.Field.Content = minesweeper.Unknown
		})
		tile.UpdateContent()

		assert := assert.New(t)

		assert.False(tile.icon.Hidden)
		assert.Equal(assets.ResourceFlagSuccessPng, tile.icon.Resource)
		assert.Equal(TileDefaultColor, tile.background.FillColor)
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
			tile.Field.Checked = false
		})

		assert := assert.New(t)

		for i := 0; i < 10; i++ {
			tile.Field.Content = minesweeper.FieldContent(i)
			tile.UpdateContent()

			assert.True(tile.label.Hidden)
		}

		tile.Field.Checked = true

		for i := 1; i < 9; i++ {
			tile.Field.Content = minesweeper.FieldContent(i)
			tile.UpdateContent()

			assert.False(tile.label.Hidden)
			assert.Equal(strconv.Itoa(i), tile.label.Text)
			assert.Equal(TileTextColor[i], tile.label.Color)
		}

		tile.Field.Content = minesweeper.FieldContent(0)
		tile.UpdateContent()

		assert.True(tile.label.Hidden)

		tile.Field.Content = minesweeper.FieldContent(9)
		tile.UpdateContent()

		assert.True(tile.label.Hidden)
	})
}

func TestTileReset(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	tile := g.Tiles[1][2]
	tile.CreateRenderer()

	tile.Flagged = true
	tile.Field.Checked = true
	tile.Field.Content = minesweeper.Mine
	tile.Marker = HelpMarkingMine

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

	assert.Equal(HelpMarkingNone, tile.Marker)
}

func TestTileUntappable(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
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

	game := tile.grid.Game.(*minesweeper.LocalGame)

	game.GameOver = true
	assert.True(tile.untappable())
	game.GameOver = false

	game.GameWon = true
	assert.True(tile.untappable())
	game.GameWon = false
}
