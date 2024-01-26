package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	timer := NewTimer()

	assert := assert.New(t)

	assert.Equal(0, timer.Seconds)
	assert.Equal("0000", timer.Label.Text)

	timer.Start()
	assert.True(timer.running)

	time.Sleep(10 * time.Second)

	timer.Stop()

	assert.False(timer.running)
	assert.Equal(10, timer.Seconds)
	assert.Equal("0010", timer.Label.Text)

	timer.Start()
	assert.True(timer.running)

	timer.Reset()
	assert.False(timer.running)
	assert.Equal(0, timer.Seconds)
	assert.Equal("0000", timer.Label.Text)
}
