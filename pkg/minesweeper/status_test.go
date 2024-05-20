package minesweeper

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/heathcliff26/go-minesweeper/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestAssistedMode(t *testing.T) {
	tCase := "assisted_mode_1"

	assert := assert.New(t)

	save, err := LoadSave("testdata/" + tCase + ".sav")
	if !assert.Nil(err, "Should load savegame") {
		t.FailNow()
	}
	buf, err := os.ReadFile("testdata/" + tCase + ".json")
	if !assert.Nil(err, "Should load test config") {
		t.FailNow()
	}
	var testConfig struct {
		StartPos Pos
		Mines    []Pos
		SafePos  []Pos
	}
	err = json.Unmarshal(buf, &testConfig)
	if !assert.Nil(err, "Should parse test config") {
		t.FailNow()
	}

	s := save.Game().CheckField(testConfig.StartPos)

	assert.ElementsMatch(testConfig.Mines, s.ObviousMines(), "Mines should match")
	assert.Equal(len(testConfig.Mines), len(s.actions.Mines), "Should have the same amount of mines")
	assert.Equal(len(testConfig.SafePos), len(s.actions.SafePos), "Should already have safe positions")
	assert.ElementsMatch(testConfig.SafePos, s.ObviousSafePos(), "Safe positions should match")

	s.actions.SafePos = nil
	assert.ElementsMatch(testConfig.SafePos, s.ObviousSafePos(), "Safe positions should match")
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
				GameOver: true,
			},
		},
		{
			Name: "FieldGameWon",
			Status: Status{
				Field:   utils.Make2D[Field](1, 1),
				GameWon: true,
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
