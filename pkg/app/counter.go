package app

import (
	"strconv"

	"fyne.io/fyne/v2/canvas"
)

// Displays a counter that can be increased and decreased
type Counter struct {
	Label *canvas.Text
	Count int
}

// Create new counter
func NewCounter(count int) *Counter {
	return &Counter{
		Label: newGridLabel(strconv.Itoa(count)),
		Count: count,
	}
}

// Set the count to a specific number
func (m *Counter) SetCount(c int) {
	m.Count = c
	m.refresh()
}

// Redraw the counter from the current count
func (m *Counter) refresh() {
	if m.Count < 0 {
		m.Label.Text = "0"
	} else {
		m.Label.Text = strconv.Itoa(m.Count)
	}
	m.Label.Refresh()
}

// Increase the counter
func (m *Counter) Inc() {
	m.Count++
	m.refresh()
}

// Decrease the counter
func (m *Counter) Dec() {
	m.Count--
	m.refresh()
}
