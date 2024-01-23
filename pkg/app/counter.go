package app

import (
	"strconv"

	"fyne.io/fyne/v2/canvas"
)

type MineCount struct {
	Label *canvas.Text
	Count int
}

func NewMineCount(count int) *MineCount {
	return &MineCount{
		Label: newGridLabel(strconv.Itoa(count)),
		Count: count,
	}
}

func (m *MineCount) SetCount(c int) {
	m.Count = c
	m.refresh()
}

func (m *MineCount) refresh() {
	if m.Count < 0 {
		m.Label.Text = "0"
	} else {
		m.Label.Text = strconv.Itoa(m.Count)
	}
	m.Label.Refresh()
}

func (m *MineCount) Inc() {
	m.Count++
	m.refresh()
}

func (m *MineCount) Dec() {
	m.Count--
	m.refresh()
}
