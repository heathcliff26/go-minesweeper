package minesweeper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	d := Difficulties()[DifficultyIntermediate]
	g := NewGameWithSafePos(d, NewPos(0, 0))

	var s *Save
	t.Run("New", func(t *testing.T) {
		tmp, err := g.ToSave()
		s = tmp

		assert := assert.New(t)

		assert.NoError(err)

		assert.NotNil(s)
		assert.NotEmpty(s.ID)
		assert.Equal(d.Mines, len(s.Data.Mines))
		assert.Equal(d, s.Data.Difficulty)

		assert.True(g.IsReplay())
	})
	t.Run("SaveAndLoad", func(t *testing.T) {
		path := "tmp-test" + SaveFileExtension
		err := s.Save(path)

		assert := assert.New(t)

		require.NoError(t, err, "Failed to save to file")
		t.Cleanup(func() {
			_ = os.Remove(path)
		})

		s2, err := LoadSave(path)

		assert.NoError(err)
		assert.Equal(s, s2)
	})
	t.Run("SaveAddsExtension", func(t *testing.T) {
		path := "tmp-test2"
		err := s.Save(path)

		assert := assert.New(t)

		require.NoError(t, err, "Failed to save to file")
		t.Cleanup(func() {
			_ = os.Remove(path + SaveFileExtension)
		})

		_, err = os.Stat(path + SaveFileExtension)
		assert.NoError(err, "Warning, can't delete save file, as name is not known")
	})
	t.Run("Game", func(t *testing.T) {
		g2 := s.Game()

		assert := assert.New(t)

		assert.True(g2.IsReplay())

		g2.replay = false

		assert.Equal(d, g2.Difficulty())
		assert.ElementsMatch(g.Field, g2.Field)
	})
}
