package app

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type Button struct {
	widget.BaseWidget

	Label  *canvas.Text
	Action func()
}

func NewButton(text string, color color.Color, action func()) *Button {
	b := &Button{
		Label:  canvas.NewText(text, color),
		Action: action,
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *Button) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(b.Label)
}

func (b *Button) Tapped(_ *fyne.PointEvent) {
	b.Action()
}

func (b *Button) TappedSecondary(_ *fyne.PointEvent) {}

func (b *Button) SetText(text string) {
	b.Label.Text = text
	b.Refresh()
}
