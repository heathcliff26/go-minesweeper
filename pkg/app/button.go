package app

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// Custom implementation for a button that is only portraied by text.
// The text is fully configurable by exposing the backing label.
type Button struct {
	widget.BaseWidget

	Label  *canvas.Text
	Action func()
}

// Create a new button from the given text, color and function
func NewButton(text string, color color.Color, action func()) *Button {
	b := &Button{
		Label:  canvas.NewText(text, color),
		Action: action,
	}
	b.ExtendBaseWidget(b)
	return b
}

// Function to create renderer needed to implement widget
func (b *Button) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(b.Label)
}

// Left click action
func (b *Button) Tapped(_ *fyne.PointEvent) {
	b.Action()
}

// Right click action, currently not implemented or exposed
func (b *Button) TappedSecondary(_ *fyne.PointEvent) {}

// Set label text to the given string and refresh widget
func (b *Button) SetText(text string) {
	b.Label.Text = text
	fyne.Do(b.Refresh)
}
