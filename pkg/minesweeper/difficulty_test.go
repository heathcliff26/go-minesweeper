package minesweeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDifficulties(t *testing.T) {
	res := Difficulties()

	assert := assert.New(t)

	assert.Equal(difficulties, res)
	assert.NotSame(difficulties, res)
}
