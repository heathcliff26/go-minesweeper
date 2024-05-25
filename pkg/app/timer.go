package app

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"fyne.io/fyne/v2/canvas"
)

// Timer to display an upwards countdown in a fyne app.
type Timer struct {
	Label      *canvas.Text
	Seconds    int
	stopSignal chan bool
	running    bool
	lock       sync.Mutex
}

// Create new timer
func NewTimer() *Timer {
	return &Timer{
		Label:      newGridLabel("0000"),
		Seconds:    0,
		stopSignal: make(chan bool),
	}
}

// Start the timer, runs concurrently
func (t *Timer) Start() {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.running {
		return
	}

	ticker := time.NewTicker(time.Second)
	t.running = true
	go func() {
		slog.Debug("Started timer")
		for {
			select {
			case <-t.stopSignal:
				ticker.Stop()
				return
			case <-ticker.C:
				t.Seconds++
				t.refresh()
			}
		}
	}()
}

// Stop the timer
func (t *Timer) Stop() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.stop()
}

// Reset the timer back to zero
func (t *Timer) Reset() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.stop()
	t.Seconds = 0
	t.refresh()
}

// Actual stop logic, put here so lock can be aquired first by caller
func (t *Timer) stop() {
	if t.Running() {
		t.stopSignal <- true
		t.running = false
		slog.Info("Stopped timer", slog.Int("seconds", t.Seconds))
	}
}

// Check if the timer is running
func (t *Timer) Running() bool {
	return t.running
}

// Refresh the timer from it's current values
func (t *Timer) refresh() {
	t.Label.Text = fmt.Sprintf("%04d", t.Seconds)
	t.Label.Refresh()
}
