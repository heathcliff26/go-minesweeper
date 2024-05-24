package minesweeper

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/heathcliff26/go-minesweeper/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestAssistedMode(t *testing.T) {
	tMatrix := []string{"assisted_mode_1", "assisted_mode_2"}
	for _, tCase := range tMatrix {
		t.Run(tCase, func(t *testing.T) {
			assert := assert.New(t)

			save, err := LoadSave("testdata/" + tCase + ".sav")
			if !assert.Nil(err, "Should load savegame") {
				t.FailNow()
			}
			game := save.Game()

			buf, err := os.ReadFile("testdata/" + tCase + ".json")
			if !assert.Nil(err, "Should load test config") {
				t.FailNow()
			}
			var testConfig []struct {
				CheckPos Pos
				Mines    []Pos
				SafePos  []Pos
			}
			err = json.Unmarshal(buf, &testConfig)
			if !assert.Nil(err, "Should parse test config") {
				t.FailNow()
			}

			for i, step := range testConfig {
				s := game.CheckField(step.CheckPos)

				fail := false

				if !assert.ElementsMatch(step.Mines, s.ObviousMines(), "Mines should match") {
					fail = true
				}
				if !assert.Equal(len(step.Mines), len(s.actions.Mines), "Should have the same amount of mines") {
					fail = true
				}
				if !assert.Equal(len(step.SafePos), len(s.actions.SafePos), "Should already have safe positions") {
					fail = true
				}
				if !assert.ElementsMatch(step.SafePos, s.ObviousSafePos(), "Safe positions should match") {
					fail = true
				}

				s.actions.SafePos = nil
				if !assert.ElementsMatch(step.SafePos, s.ObviousSafePos(), "Safe positions should match") {
					fail = true
				}

				for _, p := range s.ObviousMines() {
					if !assert.Falsef(s.Field[p.X][p.Y].Checked, "Mine should not be checked, pos="+p.String()) {
						fail = true
					}
				}
				for _, p := range s.ObviousSafePos() {
					if !assert.Falsef(s.Field[p.X][p.Y].Checked, "Safe position should not be checked, pos="+p.String()) {
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

func TestCreateActionsEarlyReturn(t *testing.T) {
	tMatrix := []struct {
		Name   string
		Status Status
	}{
		{
			Name: "FieldNil",
			Status: Status{
				Field: nil,
			},
		},
		{
			Name: "FieldEmpty",
			Status: Status{
				Field: make([][]Field, 0),
			},
		},
		{
			Name: "FieldGameOver",
			Status: Status{
				Field:    utils.Make2D[Field](1, 1),
				gameOver: true,
			},
		},
		{
			Name: "FieldGameWon",
			Status: Status{
				Field:   utils.Make2D[Field](1, 1),
				gameWon: true,
			},
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			tCase.Status.createActions()

			assert.NotNil(tCase.Status.actions.Mines, "Mines should not be nil")
			assert.NotNil(tCase.Status.actions.SafePos, "SafePos should not be nil")
			assert.Empty(tCase.Status.actions.Mines, "Mines should be empty")
			assert.Empty(tCase.Status.actions.SafePos, "SafePos should be empty")
		})
	}
}
