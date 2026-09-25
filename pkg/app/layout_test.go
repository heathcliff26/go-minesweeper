package app

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/stretchr/testify/assert"
)

func TestFixedRatioLayout(t *testing.T) {
	layout := NewFixedRatioLayout(16, 9, 32)
	t.Run("New", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(float32(16), layout.width, "Should match width")
		assert.Equal(float32(9), layout.height, "Should match height")
		assert.Equal(float32(32), layout.minWidth, "Should match minimum width")
	})
	t.Run("MinSize", func(t *testing.T) {
		assert := assert.New(t)

		size := layout.MinSize([]fyne.CanvasObject{})
		assert.Equal(fyne.NewSize(32, 18), size, "Should return default size")

		l1 := canvas.NewText("Hello", color.White)
		l1.TextSize = 16
		l2 := canvas.NewText("Hello", color.White)
		l2.TextSize = 20

		size = layout.MinSize([]fyne.CanvasObject{l1, l2})
		assert.Equal(l2.MinSize(), size, "Should return max size")

		size = layout.MinSize([]fyne.CanvasObject{l2, l1})
		assert.Equal(l2.MinSize(), size, "Should return max size")
	})
	t.Run("CalculateSize", func(t *testing.T) {
		tests := []struct {
			name string
			size fyne.Size
			want fyne.Size
		}{
			{name: "WiderThanAspectRatio", size: fyne.NewSize(32, 9), want: fyne.NewSize(16, 9)},
			{name: "TallerThanAspectRatio", size: fyne.NewSize(16, 18), want: fyne.NewSize(16, 9)},
			{name: "MatchingAspectRatio", size: fyne.NewSize(32, 18), want: fyne.NewSize(32, 18)},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				assert.Equal(t, test.want, layout.calculateSize(test.size))
			})
		}
	})
	t.Run("Layout", func(t *testing.T) {
		assert := assert.New(t)

		objects := []fyne.CanvasObject{
			canvas.NewRectangle(color.White),
			canvas.NewRectangle(color.White),
		}

		layout.Layout(objects, fyne.NewSize(32, 9))

		for _, object := range objects {
			assert.Equal(fyne.NewSize(16, 9), object.Size(), "Should match size")
			assert.Equal(fyne.NewPos(8, 0), object.Position(), "Should match position")
		}
	})
}
