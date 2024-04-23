package app

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	count := 10
	c := NewCounter(count)

	t.Run("Basic", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(strconv.Itoa(count), c.Label.Text)
		assert.Equal(count, c.Count)
	})
	t.Run("Inc", func(t *testing.T) {
		count++
		c.Inc()

		assert := assert.New(t)

		assert.Equal(strconv.Itoa(count), c.Label.Text)
		assert.Equal(count, c.Count)
	})
	t.Run("Dec", func(t *testing.T) {
		count--
		c.Dec()

		assert := assert.New(t)

		assert.Equal(strconv.Itoa(count), c.Label.Text)
		assert.Equal(count, c.Count)
	})
	t.Run("SetCount", func(t *testing.T) {
		count = 20
		c.SetCount(count)

		assert := assert.New(t)

		assert.Equal(strconv.Itoa(count), c.Label.Text)
		assert.Equal(count, c.Count)
	})
	t.Run("Refresh", func(t *testing.T) {
		c.SetCount(-1)

		assert := assert.New(t)

		assert.Equal("00", c.Label.Text, "Should not display negative numbers")
		assert.Equal(-1, c.Count)
	})
}
