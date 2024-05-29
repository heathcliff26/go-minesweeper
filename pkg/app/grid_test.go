package app

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/stretchr/testify/assert"
)

type assistedModeTestConfig struct {
	CheckPos minesweeper.Pos
	Mines    []minesweeper.Pos
	SafePos  []minesweeper.Pos
}

func TestNewGrid(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, true)

	assert := assert.New(t)

	assert.Equal(g.Row(), len(g.Tiles))
	assert.Equal(g.Col(), len(g.Tiles[0]))
	assert.Equal(DEFAULT_DIFFICULTY, g.Difficulty)
	assert.Nil(g.Game)
	assert.True(g.AssistedMode)
	assert.Equal(DEFAULT_GAME_ALGORITHM, g.GameAlgorithm)
	assert.NotNil(g.Timer)
	assert.NotNil(g.MineCount)
	assert.NotNil(g.ResetButton)
}

func TestTappedTile(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}

	p := minesweeper.NewPos(0, 0)

	assert := assert.New(t)

	g.TappedTile(p)
	assert.NotNil(g.Game)
	assert.True(g.Tiles[0][0].Checked())

	game := g.Game.(*minesweeper.LocalGame)

	game.GameWon = true
	game.UpdateStatus()
	g.TappedTile(p)
	assert.Equal(ResetGameWonText, g.ResetButton.Label.Text)
	game.GameWon = false

	game.GameOver = true
	game.UpdateStatus()
	g.TappedTile(p)
	assert.Equal(ResetGameOverText, g.ResetButton.Label.Text)
}

func TestGameAlgorithm(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}

	tMatrix := []string{"SafePos", "safeArea"}
	for i, tCase := range tMatrix {
		t.Run(tCase, func(t *testing.T) {
			t.Cleanup(g.NewGame)
			g.GameAlgorithm = i
			g.TappedTile(minesweeper.NewPos(0, 0))
			assert.NotNil(t, g.Game)
		})
	}
	t.Run("Unknown", func(t *testing.T) {
		t.Cleanup(g.NewGame)
		g.GameAlgorithm = -1
		g.TappedTile(minesweeper.NewPos(0, 0))
		assert.Nil(t, g.Game)
	})
}

func TestNewGame(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
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
			assert.Falsef(g.Tiles[x][y].Checked(), "(%d, %d) All fields should be reset", x, y)
		}
	}

	assert.Nil(g.Game)
	assert.Equal(g.Difficulty.Mines, g.MineCount.Count)
	assert.False(g.Timer.running)
	assert.Equal(ResetDefaultText, g.ResetButton.Label.Text)
}

func TestReplay(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}

	g.TappedTile(minesweeper.NewPos(0, 0))
	g.Replay()

	assert := assert.New(t)

	for x := 0; x < g.Row(); x++ {
		for y := 0; y < g.Col(); y++ {
			assert.Falsef(g.Tiles[x][y].Checked(), "(%d, %d) All fields should be reset", x, y)
		}
	}

	assert.NotNil(g.Game)
	assert.True(g.Game.IsReplay())
	assert.Equal(g.Difficulty.Mines, g.MineCount.Count)
	assert.False(g.Timer.running)
	assert.Equal(ResetDefaultText, g.ResetButton.Label.Text)
}

func TestUpdateFromStatus(t *testing.T) {
	t.Parallel()
	t.Run("NilPointer", func(t *testing.T) {
		t.Parallel()

		g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
		// Should not panic
		g.updateFromStatus(nil)
	})

	difficulties := minesweeper.Difficulties()
	for _, i := range []int{10, 11, 25, 30, 50} {
		d, err := minesweeper.NewCustomDifficulty((i*i)/2, i, i)
		if err != nil {
			t.Fatalf("Failed to create custom difficulty %dx%d: %v", i, i, err)
		}
		d.Name = fmt.Sprintf("%dx%d", i, i)
		difficulties = append(difficulties, d)
	}

	for _, d := range difficulties {
		t.Run(d.Name, func(t *testing.T) {
			if d.Row*d.Col < 1000 {
				t.Parallel()
			}

			g := NewMinesweeperGrid(d, false)
			for _, row := range g.Tiles {
				for _, tile := range row {
					tile.CreateRenderer()
				}
			}

			s := minesweeper.NewGameWithSafePos(g.Difficulty, minesweeper.NewPos(0, 0)).UpdateStatus()

			for x := 0; x < g.Row(); x++ {
				for y := 0; y < g.Col(); y++ {
					s.Field[x][y] = minesweeper.Field{
						Checked: true,
						Content: minesweeper.Mine,
					}
				}
			}

			g.updateFromStatus(s)

			assert := assert.New(t)

			for x := 0; x < g.Row(); x++ {
				for y := 0; y < g.Col(); y++ {
					assert.Truef(g.Tiles[x][y].Checked(), "Field should be checked for tile=(%d, %d)", x, y)
				}
			}
		})
	}
}

func TestAssistedMode(t *testing.T) {
	assert := assert.New(t)

	g, testConfig, err := loadAssistedModeTest("../minesweeper/testdata/assisted_mode_1", true)
	if !assert.Nil(err, "Should have loaded the test") {
		t.Fatal(err)
	}

	for _, step := range testConfig {
		g.TappedTile(step.CheckPos)

		mines := 0
		safePos := 0
		for x := 0; x < g.Row(); x++ {
			for y := 0; y < g.Col(); y++ {
				p := minesweeper.NewPos(x, y)
				if g.Tiles[p.X][p.Y].Checked() {
					continue
				}
				switch g.Tiles[p.X][p.Y].Marking() {
				case HelpMarkingNone:
					assert.NotContains(step.Mines, p, "Should have been marked as a mines")
					assert.NotContains(step.SafePos, p, "Should have been marked as safe")
				case HelpMarkingMine:
					mines++
					assert.Contains(step.Mines, p, "Should not have been marked as a mine")
				case HelpMarkingSafe:
					safePos++
					assert.Contains(step.SafePos, p, "Should not have been marked as safe")
				default:
					assert.Fail("Found unknown Marker %d at %s", g.Tiles[x][y].Marking(), p.String())
				}
			}
		}
		assert.Equal(len(step.Mines), mines, "Should have the same amount of mines as in the config")
		assert.Equal(len(step.SafePos), safePos, "Should have the same amount of safe positions as in the config")
	}
}

func TestHint(t *testing.T) {
	t.Run("GameNil", func(t *testing.T) {
		g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
		assert.False(t, g.Hint(), "Should not give hint when game is nil")
	})
	t.Run("StatusNil", func(t *testing.T) {
		assert := assert.New(t)

		g, err := createGridFromSave("testdata/hint.sav", false)
		if !assert.Nil(err, "Should load savegame") {
			t.Fatal(err)
		}

		assert.False(g.Hint(), "Should not be able to display hint on game without status")
	})
	t.Run("GameOver", func(t *testing.T) {
		assert := assert.New(t)

		g, err := createGridFromSave("testdata/hint.sav", false)
		if !assert.Nil(err, "Should load savegame") {
			t.Fatal(err)
		}

		g.TappedTile(minesweeper.NewPos(15, 6))
		assert.False(g.Hint(), "Should not be able to display hint on failed game")
	})
	t.Run("DisplayHint", func(t *testing.T) {
		assert := assert.New(t)

		g, testConfig, err := loadAssistedModeTest("testdata/hint", false)
		if !assert.Nil(err, "Should have loaded the test") {
			t.Fatal(err)
		}
		step := testConfig[0]

		g.TappedTile(step.CheckPos)

		for _, mine := range step.Mines {
			assert.True(g.Hint(), "Should be able to display hints")
			tile := g.Tiles[mine.X][mine.Y]
			assert.Equalf(HelpMarkingMine, tile.Marking(), "Tile should be marked as mine, tile=%s", mine.String())
			tile.Flag(true)
			for x := 0; x < g.Difficulty.Row; x++ {
				for y := 0; y < g.Difficulty.Col; y++ {
					tile := g.Tiles[x][y]
					if tile.Flagged() || tile.Checked() {
						continue
					}
					if !assert.Equalf(HelpMarkingNone, tile.Marking(), "No other tiles should be marked, tile=%s, mine=%s", minesweeper.NewPos(x, y).String(), mine.String()) {
						t.FailNow()
					}
				}
			}
		}
		assert.True(g.Hint(), "Should be able to display hints")
		assert.Equal(HelpMarkingSafe, g.Tiles[step.SafePos[0].X][step.SafePos[0].Y].Marking(), "Tile should be marked as safe")
		for x := 0; x < g.Difficulty.Row; x++ {
			for y := 0; y < g.Difficulty.Col; y++ {
				p := minesweeper.NewPos(x, y)
				tile := g.Tiles[p.X][p.Y]
				if tile.Flagged() || tile.Checked() || p == step.SafePos[0] {
					continue
				}
				if !assert.Equalf(HelpMarkingNone, tile.Marking(), "No other tiles should be marked, tile=%s", p.String()) {
					t.FailNow()
				}
			}
		}
	})
	t.Run("NoHints", func(t *testing.T) {
		assert := assert.New(t)

		g, err := createGridFromSave("testdata/no-hints.sav", false)
		if !assert.Nil(err, "Should load savegame") {
			t.Fatal(err)
		}
		g.TappedTile(minesweeper.NewPos(0, 0))

		assert.False(g.Hint(), "Should find no hints")
	})
}

func TestAutosolve(t *testing.T) {
	t.Parallel()
	t.Run("GameNil", func(t *testing.T) {
		t.Parallel()

		g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)

		assert.False(t, g.Autosolve(0), "Should not run autosolve")
	})
	t.Run("StatusNil", func(t *testing.T) {
		t.Parallel()

		g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
		for _, row := range g.Tiles {
			for _, tile := range row {
				tile.CreateRenderer()
			}
		}
		g.Game = minesweeper.NewGameWithSafePos(g.Difficulty, minesweeper.NewPos(0, 0))

		assert.False(t, g.Autosolve(0), "Should not run autosolve")
	})
	t.Run("GameWon", func(t *testing.T) {
		t.Parallel()

		g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
		for _, row := range g.Tiles {
			for _, tile := range row {
				tile.CreateRenderer()
			}
		}
		game := minesweeper.NewGameWithSafePos(g.Difficulty, minesweeper.NewPos(0, 0))
		g.Game = game

		game.GameWon = true
		game.UpdateStatus()

		assert.False(t, g.Autosolve(0), "Should not run autosolve")
	})
	t.Run("GameOver", func(t *testing.T) {
		t.Parallel()

		g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
		for _, row := range g.Tiles {
			for _, tile := range row {
				tile.CreateRenderer()
			}
		}
		game := minesweeper.NewGameWithSafePos(g.Difficulty, minesweeper.NewPos(0, 0))
		g.Game = game

		game.GameOver = true
		game.UpdateStatus()

		assert.False(t, g.Autosolve(0), "Should not run autosolve")
	})
	t.Run("RestoresAssistedModeSetting", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		for _, value := range []bool{true, false} {
			g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
			for _, row := range g.Tiles {
				for _, tile := range row {
					tile.CreateRenderer()
				}
			}
			g.TappedTile(minesweeper.NewPos(0, 0))

			g.AssistedMode = value

			assert.True(g.Autosolve(0), "Should run autosolve")
			assert.Equal(value, g.AssistedMode, "Assisted mode setting should be restored")
		}
	})

	tMatrix := []struct {
		Name string
		Won  bool
	}{
		{"win", true},
		{"unfinished", false},
	}
	for _, tCase := range tMatrix {
		t.Run("Solve_"+tCase.Name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			path := "testdata/autosolve_" + tCase.Name
			g, err := createGridFromSave(path+".sav", false)
			if !assert.Nil(err, "Should load savegame") {
				t.Fatal(err)
			}
			g.TappedTile(minesweeper.NewPos(0, 0))

			assert.True(g.Autosolve(0), "Should run autosolve")

			msg1 := "unfinished"
			msg2 := "not be flagged"
			if tCase.Won {
				msg1 = "won"
				msg2 = "be flagged"
			}
			assert.Equal(tCase.Won, g.Game.Status().GameWon(), "Game should be "+msg1)

			for _, p := range g.Game.Status().ObviousMines() {
				assert.Equal(!tCase.Won, g.Tiles[p.X][p.Y].Flagged(), "Tile should "+msg2+", tile="+p.String())
			}
		})
	}

	t.Run("FlaggedMines", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		g, err := createGridFromSave("testdata/autosolve_unfinished.sav", false)
		if !assert.Nil(err, "Should load savegame") {
			t.Fatal(err)
		}
		g.TappedTile(minesweeper.NewPos(0, 0))

		assert.True(g.Autosolve(0), "Should run autosolve")
		// Run twice to check that flagged fields to not get unflagged
		assert.True(g.Autosolve(0), "Should run autosolve a second time")

		for _, p := range g.Game.Status().ObviousMines() {
			assert.True(g.Tiles[p.X][p.Y].Flagged(), "Tile should be flagged, tile="+p.String())
		}
	})
}

func createGridFromSave(path string, assistedMode bool) (*MinesweeperGrid, error) {
	save, err := minesweeper.LoadSave(path)
	if err != nil {
		return nil, err
	}
	g := NewMinesweeperGrid(save.Data.Difficulty, assistedMode)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}
	g.Game = save.Game()

	return g, nil
}

func loadAssistedModeTest(path string, assistedMode bool) (*MinesweeperGrid, []assistedModeTestConfig, error) {
	g, err := createGridFromSave(path+".sav", assistedMode)
	if err != nil {
		return nil, nil, err
	}

	buf, err := os.ReadFile(path + ".json")
	if err != nil {
		return nil, nil, err
	}

	var testConfig []assistedModeTestConfig
	err = json.Unmarshal(buf, &testConfig)
	if err != nil {
		return nil, nil, err
	}
	return g, testConfig, nil
}
