package app

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2/canvas"
)

type Timer struct {
	Label   *canvas.Text
	Seconds int
	stop    chan bool
	running bool
}

func NewTimer() *Timer {
	return &Timer{
		Label:   newGridLabel("0000"),
		Seconds: 0,
		stop:    make(chan bool),
	}
}

func (t *Timer) Start() {
	ticker := time.NewTicker(time.Second)
	t.running = true
	go func() {
		log.Println("Started timer")
		for {
			select {
			case <-t.stop:
				ticker.Stop()
				log.Printf("Stopped timer after %d seconds\n", t.Seconds)
				return
			case <-ticker.C:
				t.Seconds++
				t.refresh()
			}
		}
	}()
}

func (t *Timer) Stop() {
	if t.running {
		t.stop <- true
		t.running = false
	}
}

func (t *Timer) Reset() {
	t.Stop()
	t.Seconds = 0
	t.refresh()
}

func (t *Timer) refresh() {
	t.Label.Text = fmt.Sprintf("%04d", t.Seconds)
	t.Label.Refresh()
}
