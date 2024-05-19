package minesweeper

import (
	"encoding/json"
	"os"
	"testing"

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
	assert.ElementsMatch(testConfig.SafePos, s.ObviousSafePos(), "Safe positions should match")
}
