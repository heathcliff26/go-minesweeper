package minesweeper

import (
	"encoding/json"
	"os"
	"testing"
	"time"

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
				s, _ := game.CheckField(step.CheckPos)

				fail := false

				if !assert.False(s.actionsUpdated, "Status should need update for actions") {
					fail = true
				}

				if !assert.ElementsMatch(step.Mines, s.ObviousMines(), "Mines should match") {
					fail = true
				}
				if !assert.True(s.actionsUpdated, "Status should no longer need update for actions") {
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

				s.actionsUpdated = false
				s.actions.SafePos = nil
				if !assert.ElementsMatch(step.SafePos, s.ObviousSafePos(), "Safe positions should still match") {
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

func TestUpdateActionsEarlyReturn(t *testing.T) {
	tMatrix := []struct {
		Name   string
		Status *Status
	}{
		{
			Name: "FieldNil",
			Status: &Status{
				Field: nil,
			},
		},
		{
			Name: "FieldEmpty",
			Status: &Status{
				Field: make([][]Field, 0),
			},
		},
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

			tCase.Status.updateActions()

			assert.Nil(tCase.Status.actions.Mines, "Mines should not be nil")
			assert.Nil(tCase.Status.actions.SafePos, "SafePos should not be nil")
			assert.Empty(tCase.Status.actions.Mines, "Mines should be empty")
			assert.Empty(tCase.Status.actions.SafePos, "SafePos should be empty")
		})
	}
}

func TestStatusConcurrencySafe(t *testing.T) {
	s := NewGameWithSafePos(Difficulties()[0], NewPos(0, 0)).UpdateStatus()

	done := make(chan struct{})
	s.updateActionsCalled = func() {
		done <- struct{}{}
	}

	go s.ObviousMines()
	go s.ObviousSafePos()

	count := 0
	timeout := false
	for !timeout {
		select {
		case <-done:
			count++
		case <-time.After(time.Second):
			timeout = true
		}
	}

	assert.Equal(t, 1, count, "Should have called updateActions once")
}
