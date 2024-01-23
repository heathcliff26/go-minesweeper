package app

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/heathcliff26/go-minesweeper/assets"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
)

var (
	TileDefaultColor    = color.Gray16{32767}
	TileBackgroundColor = color.Gray16{^uint16(0)}
	TileExplodedColor   = color.RGBA{240, 10, 20, alpha}

	TileSize             = fyne.NewSize(32, 32)
	TileTextSize float32 = 23 // Biggest we can go with TileSize of 32^2
)

const alpha = ^uint8(0)

var TileTextColor = []color.Color{
	color.White,
	color.RGBA{20, 15, 220, alpha},  // 1, blue
	color.RGBA{5, 110, 20, alpha},   // 2, green
	color.RGBA{240, 10, 20, alpha},  // 3, red
	color.RGBA{5, 5, 100, alpha},    // 4, dark blue
	color.RGBA{90, 38, 42, alpha},   // 5, brown
	color.RGBA{25, 230, 230, alpha}, // 6, cyan
	color.RGBA{10, 10, 10, alpha},   // 7, black
	color.RGBA{64, 64, 64, alpha},   // 8, grey
}

type Tile struct {
	widget.BaseWidget

	background *canvas.Rectangle
	label      *canvas.Text
	icon       *widget.Icon

	x, y    int
	Field   *minesweeper.Field
	grid    *MinesweeperGrid
	Flagged bool
}

func NewTile(x, y int, grid *MinesweeperGrid) *Tile {
	t := &Tile{
		x: x,
		y: y,
		Field: &minesweeper.Field{
			Checked: false,
			Content: minesweeper.Unknown,
		},
		grid: grid,
	}
	t.ExtendBaseWidget(t)

	return t
}

func (t *Tile) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)

	t.background = canvas.NewRectangle(TileDefaultColor)
	t.background.SetMinSize(TileSize)

	t.label = canvas.NewText("", color.White)
	t.label.TextStyle.Bold = true
	t.label.Alignment = fyne.TextAlignCenter
	t.label.TextSize = TileTextSize
	t.label.Hidden = true

	t.icon = widget.NewIcon(nil)
	t.icon.Resize(TileSize)
	t.icon.Hidden = true

	content := container.NewStack(t.background, t.icon, t.label)
	return widget.NewSimpleRenderer(content)
}

func (t *Tile) Tapped(_ *fyne.PointEvent) {
	if t.Field.Checked || t.Flagged {
		return
	}
	t.grid.TappedTile(t.x, t.y)
}

func (t *Tile) TappedSecondary(_ *fyne.PointEvent) {
	if t.Field.Checked {
		return
	}
	if t.Flagged {
		t.grid.MineCount.Inc()
	} else {
		t.grid.MineCount.Dec()
	}
	t.Flagged = !t.Flagged
	t.UpdateContent()
}

func (t *Tile) UpdateContent() {
	t.icon.Hidden = true
	t.label.Hidden = true
	defer t.Refresh()
	if t.Flagged && !t.Field.Checked {
		t.icon.SetResource(assets.ResourceFlagPng)
		t.icon.Hidden = false
		return
	}
	if t.Field.Checked {
		t.background.FillColor = TileBackgroundColor
	} else {
		t.background.FillColor = TileDefaultColor
	}
	if t.Field.Content == minesweeper.Unknown {
		return
	}

	if t.Field.Content == minesweeper.Mine {
		t.icon.SetResource(assets.ResourceMinePng)
		t.icon.Hidden = false
		if t.Field.Checked {
			t.background.FillColor = TileExplodedColor
		}
	}
	if t.Field.Content > 0 && t.Field.Content < 9 {
		t.label.Text = strconv.Itoa(int(t.Field.Content))
		t.label.Color = TileTextColor[t.Field.Content]
		t.label.Hidden = false
	}
}

func (t *Tile) Reset() {
	t.Flagged = false
	t.Field.Checked = false
	t.Field.Content = minesweeper.Unknown
	t.UpdateContent()
}
