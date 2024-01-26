package app

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestButton(t *testing.T) {
	tapped := false

	b := NewButton("test", color.White, func() { tapped = true })

	t.Run("Render", func(t *testing.T) {
		r := test.WidgetRenderer(b)
		items := r.Objects()

		assert.Equal(t, []fyne.CanvasObject{b.Label}, items, "Button should only consist of a label")
	})
	t.Run("Tappable", func(t *testing.T) {
		test.Tap(b)

		assert.True(t, tapped, "Button should execute action on tap")
	})
	t.Run("Text", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal("test", b.Label.Text)
		assert.Equal(color.White, b.Label.Color)

		b.SetText("test2")
		assert.Equal("test2", b.Label.Text)
	})
}
