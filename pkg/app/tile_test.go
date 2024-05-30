package app

import (
	"strconv"
	"testing"
	"time"

	"github.com/heathcliff26/go-minesweeper/assets"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/stretchr/testify/assert"
)

func TestNewTile(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	tile := g.Tiles[1][2]

	t.Run("New", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(minesweeper.NewPos(1, 2), tile.pos)

		assert.Equal(g, tile.grid)

		f := minesweeper.Field{
			Checked: false,
			Content: minesweeper.Unknown,
		}
		assert.Equal(f, tile.field)

		assert.False(tile.Flagged())
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
	t.Parallel()
	assert := assert.New(t)

	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}
	tile := g.Tiles[0][0]

	tile.Flag(true)
	tile.Tapped(nil)
	time.Sleep(time.Second)

	assert.Nil(tile.grid.Game, "Flagged tiles should not trigger checks")
	assert.False(tile.field.Checked, "Flagged tiles should not trigger checks")

	tile.Flag(false)
	tile.Tapped(nil)
	time.Sleep(time.Second)

	assert.NotNil(tile.grid.Game, "Game should be started")
	assert.True(tile.field.Checked, "Tile should be checked")
}

func TestTileTappedSecondary(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	tile := g.Tiles[1][2]
	tile.CreateRenderer()

	count := tile.grid.MineCount.Count

	tile.field.Checked = true
	tile.TappedSecondary(nil)
	time.Sleep(time.Second)

	assert.Equal(count, tile.grid.MineCount.Count, "MineCount should not have changed")
	assert.False(tile.Flagged(), "Tile should not be flagged")

	tile.field.Checked = false
	tile.TappedSecondary(nil)
	time.Sleep(time.Second)

	assert.Equal(count-1, tile.grid.MineCount.Count, "MineCount should be decreased")
	assert.True(tile.Flagged(), "Tile should be flagged")

	tile.TappedSecondary(nil)
	time.Sleep(time.Second)

	assert.Equal(count, tile.grid.MineCount.Count, "MineCount should be back to original value")
	assert.False(tile.Flagged(), "Tile should not be flagged")
}

func TestDoubleTapped(t *testing.T) {
	t.Parallel()

	tMatrix := []struct {
		Name    string
		Flagged bool
	}{
		{"Unchecked", false},
		{"Flagged", true},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			t.Parallel()

			g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
			for _, row := range g.Tiles {
				for _, tile := range row {
					tile.CreateRenderer()
				}
			}
			tile := g.Tiles[0][0]
			g.Game = minesweeper.NewGameWithSafePos(g.Difficulty, tile.pos)

			tile.Flag(tCase.Flagged)
			t.Cleanup(func() {
				tile.Flag(false)
			})
			tile.DoubleTapped(nil)
			time.Sleep(time.Second)

			for m := -1; m < 2; m++ {
				for n := -1; n < 2; n++ {
					p := tile.pos
					p.X += m
					p.Y += n
					if g.OutOfBounds(p) {
						continue
					}

					assert.Falsef(t, g.Tiles[p.X][p.Y].Checked(), p.String())
				}
			}
		})
	}

	t.Run("Checked", func(t *testing.T) {
		t.Parallel()
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
		g.testChannel = make(chan string, 2)

		tile.Tapped(nil)
		<-g.testChannel

		g.Tiles[0][0].Flag(true)

		tile.DoubleTapped(nil)
		<-g.testChannel

		assert.False(g.Game.Status().GameOver(), "Game should not be lost")
		assert.True(g.Game.Status().GameWon(), "Game should be won")
		assert.False(g.Tiles[0][0].Checked(), "Flagged field should not be checked")

		for x := 0; x < 2; x++ {
			for y := 0; y < 2; y++ {
				if x == 0 && y == 0 {
					continue
				}
				assert.Truef(g.Tiles[x][y].Checked(), "Field should be checked, tile=(%d, %d)", x, y)
			}
		}
	})
}

func TestTileUpdateContent(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	tile := g.Tiles[1][2]
	tile.CreateRenderer()

	tile.updateContent()

	t.Run("Default", func(t *testing.T) {
		assert := assert.New(t)

		assert.True(tile.icon.Hidden)
		assert.True(tile.label.Hidden)
		assert.Equal(TileDefaultColor, tile.background.FillColor)
	})
	t.Run("Checked", func(t *testing.T) {
		tile.field.Checked = true
		t.Cleanup(func() {
			tile.field.Checked = false
		})
		tile.updateContent()

		assert.Equal(t, TileBackgroundColor, tile.background.FillColor)
	})
	t.Run("Flagged", func(t *testing.T) {

		tile.Flag(true)
		t.Cleanup(func() {
			tile.Flag(false)
		})
		tile.updateContent()

		assert := assert.New(t)

		assert.False(tile.icon.Hidden)
		assert.Equal(assets.ResourceFlagPng, tile.icon.Resource)
		assert.Equal(TileDefaultColor, tile.background.FillColor)
	})
	t.Run("FlaggedSuccess", func(t *testing.T) {
		tile.Flag(true)
		tile.field.Content = minesweeper.Mine
		t.Cleanup(func() {
			tile.Flag(false)
			tile.field.Content = minesweeper.Unknown
		})
		tile.updateContent()

		assert := assert.New(t)

		assert.False(tile.icon.Hidden)
		assert.Equal(assets.ResourceFlagSuccessPng, tile.icon.Resource)
		assert.Equal(TileDefaultColor, tile.background.FillColor)
	})
	t.Run("Mine", func(t *testing.T) {
		tile.field.Content = minesweeper.Mine
		t.Cleanup(func() {
			tile.field.Content = minesweeper.Unknown
		})
		tile.updateContent()

		assert := assert.New(t)

		assert.False(tile.icon.Hidden)
		assert.Equal(assets.ResourceMinePng, tile.icon.Resource)
		assert.Equal(TileDefaultColor, tile.background.FillColor)

		tile.field.Checked = true
		t.Cleanup(func() {
			tile.field.Checked = false
		})
		tile.updateContent()

		assert.False(tile.icon.Hidden)
		assert.Equal(assets.ResourceMinePng, tile.icon.Resource)
		assert.Equal(TileExplodedColor, tile.background.FillColor)
	})
	t.Run("Numbers", func(t *testing.T) {
		t.Cleanup(func() {
			tile.field.Content = minesweeper.Unknown
			tile.field.Checked = false
		})

		assert := assert.New(t)

		for i := 0; i < 10; i++ {
			tile.field.Content = minesweeper.FieldContent(i)
			tile.updateContent()

			assert.True(tile.label.Hidden)
		}

		tile.field.Checked = true

		for i := 1; i < 9; i++ {
			tile.field.Content = minesweeper.FieldContent(i)
			tile.updateContent()

			assert.False(tile.label.Hidden)
			assert.Equal(strconv.Itoa(i), tile.label.Text)
			assert.Equal(TileTextColor[i], tile.label.Color)
		}

		tile.field.Content = minesweeper.FieldContent(0)
		tile.updateContent()

		assert.True(tile.label.Hidden)

		tile.field.Content = minesweeper.FieldContent(9)
		tile.updateContent()

		assert.True(tile.label.Hidden)
	})
}

func TestTileReset(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	tile := g.Tiles[1][2]
	tile.CreateRenderer()

	tile.Flag(true)
	tile.field.Checked = true
	tile.field.Content = minesweeper.Mine
	tile.Mark(HelpMarkingMine)

	tile.updateContent()
	tile.Reset()

	assert := assert.New(t)

	assert.False(tile.Flagged())
	assert.False(tile.field.Checked)
	assert.Equal(minesweeper.Unknown, tile.field.Content)

	assert.Equal(TileDefaultColor, tile.background.FillColor)
	assert.Equal(TileSize, tile.background.MinSize())

	assert.Equal("", tile.label.Text)
	assert.True(tile.label.TextStyle.Bold)
	assert.Equal(TileTextSize, tile.label.TextSize)
	assert.True(tile.label.Hidden)

	assert.True(tile.icon.Hidden)

	assert.Equal(HelpMarkingNone, tile.Marking())
}

func TestTileUntappable(t *testing.T) {
	t.Parallel()

	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}
	tile := g.Tiles[1][2]

	assert := assert.New(t)

	assert.False(tile.untappable())

	tile.field.Checked = true
	assert.True(tile.untappable())
	tile.field.Checked = false

	tile.Tapped(nil)
	time.Sleep(time.Second)

	tile.field.Checked = false
	assert.False(tile.untappable())

	game := tile.grid.Game.(*minesweeper.LocalGame)

	game.GameOver = true
	assert.True(tile.untappable())
	game.GameOver = false

	game.GameWon = true
	assert.True(tile.untappable())
	game.GameWon = false
}
