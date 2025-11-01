package minesweeper

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/heathcliff26/go-minesweeper/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSolver(t *testing.T) {
	assert := assert.New(t)

	game := NewGameWithSafePos(Difficulties()[0], NewPos(0, 0))
	solver := NewSolver(game)
	assert.NotNil(solver, "Solver should not be nil")
	assert.NotNil(solver.mines, "Mines slice should not be nil")
	assert.NotNil(solver.nextSteps, "Next steps slice should not be nil")
	assert.Equal(game, solver.game, "Solver should contain reference to the game")
}

func TestSolverAutosolve(t *testing.T) {
	for i := 1; i < 6; i++ {
		t.Run("Solvable-"+strconv.Itoa(i), func(t *testing.T) {
			s, err := LoadSave("testdata/autosolve-solvable-" + strconv.Itoa(i) + ".sav")
			require.NoError(t, err, "Failed to load savegame")

			g := s.Game()
			solver := NewSolver(g)
			assert.True(t, solver.Autosolve(NewPos(0, 0)))
			assert.True(t, g.Won())
		})
	}
	for i := 1; i < 6; i++ {
		t.Run("Unsolvable-"+strconv.Itoa(i), func(t *testing.T) {
			s, err := LoadSave("testdata/autosolve-unsolvable-" + strconv.Itoa(i) + ".sav")
			require.NoError(t, err, "Failed to load savegame")

			g := s.Game()
			solver := NewSolver(g)
			assert.False(t, solver.Autosolve(NewPos(0, 0)))
			assert.False(t, g.Won())
		})
	}
}

func TestAssistedMode(t *testing.T) {
	tMatrix := []string{"assisted_mode_1", "assisted_mode_2"}
	for _, tCase := range tMatrix {
		t.Run(tCase, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			save, err := LoadSave("testdata/" + tCase + ".sav")
			require.NoError(err, "Should load savegame")
			game := save.Game()

			buf, err := os.ReadFile("testdata/" + tCase + ".json")
			require.NoError(err, "Should load test config")
			var testConfig []struct {
				CheckPos Pos
				Mines    []Pos
				SafePos  []Pos
			}
			err = json.Unmarshal(buf, &testConfig)
			require.NoError(err, "Should parse test config")

			s := NewSolver(game)

			for i, step := range testConfig {
				status, _ := game.CheckField(step.CheckPos)
				s.Update()

				fail := false

				if !assert.ElementsMatch(step.Mines, s.KnownMines(), "Mines should match") {
					fail = true
				}
				if !assert.Equal(len(step.Mines), len(s.KnownMines()), "Should have the same amount of mines") {
					fail = true
				}
				if !assert.Equal(len(step.SafePos), len(s.NextSteps()), "Should already have safe positions") {
					fail = true
				}
				if !assert.ElementsMatch(step.SafePos, s.NextSteps(), "Safe positions should match") {
					fail = true
				}

				for _, p := range s.KnownMines() {
					if !assert.Falsef(status.Field[p.X][p.Y].Checked, "Mine should not be checked, pos="+p.String()) {
						fail = true
					}
				}
				for _, p := range s.NextSteps() {
					if !assert.Falsef(status.Field[p.X][p.Y].Checked, "Safe position should not be checked, pos="+p.String()) {
						fail = true
					}
				}

				if fail {
					t.Fatalf("Failed on step %d, checkPos=%s", i, step.CheckPos.String())
				}
			}
		})
	}
}

func TestSolverUpdateEarlyReturn(t *testing.T) {
	tMatrix := []struct {
		Name   string
		Status *Status
	}{
		{
			Name: "FieldGameOver",
			Status: &Status{
				Field:    utils.Make2D[Field](1, 1),
				gameOver: true,
			},
		},
		{
			Name: "FieldGameWon",
			Status: &Status{
				Field:   utils.Make2D[Field](1, 1),
				gameWon: true,
			},
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			g := NewGameWithSafePos(Difficulties()[0], NewPos(0, 0))
			s := NewSolver(g)
			g.UpdateStatus()

			g.GameOver = tCase.Status.gameOver
			g.GameWon = tCase.Status.gameWon

			s.Update()

			assert.NotNil(s.mines, "mines should not be nil")
			assert.NotNil(s.nextSteps, "nextSteps should not be nil")
			assert.Empty(s.mines, "mines should be empty")
			assert.Empty(s.nextSteps, "nextSteps should be empty")
		})
	}
	t.Run("GameIsNil", func(t *testing.T) {
		s := &Solver{}
		assert.NotPanics(t, func() {
			s.Update()
		}, "Update should not panic when game is nil")
	})
}
