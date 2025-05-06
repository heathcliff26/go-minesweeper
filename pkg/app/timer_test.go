package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	t.Parallel()
	t.Run("New", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		timer := NewTimer()

		assert.Equal(0, timer.Seconds)
		assert.Equal("0000", timer.Label.Text)
	})
	t.Run("Counting", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		timer := NewTimer()

		timer.Start()
		assert.True(timer.running)

		time.Sleep(10*time.Second + 500*time.Millisecond)
		timer.Stop()

		assert.False(timer.running)
		assert.Equal(10, timer.Seconds)
		assert.Equal("0010", timer.Label.Text)
	})
	t.Run("Reset", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		timer := NewTimer()

		timer.Start()
		assert.True(timer.running)

		time.Sleep(2 * time.Second)

		timer.Reset()
		assert.False(timer.running)
		assert.Equal(0, timer.Seconds)
		assert.Equal("0000", timer.Label.Text)
	})
	t.Run("MultipleCallsToStart", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		timer := NewTimer()

		timer.Start()
		assert.True(timer.running)
		timer.Start()
		assert.True(timer.running)

		time.Sleep(10*time.Second + 500*time.Millisecond)
		timer.Stop()

		assert.Equal(10, timer.Seconds)
	})
}
