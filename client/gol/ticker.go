package gol

import (
	"time"
)

// Ticker is used to send AliveCellsCount events every 2 seconds
type Ticker struct {
	stop chan bool
	tick chan bool
}

func (t *Ticker) startTicker(events chan<- Event) {
	ticker := time.NewTicker(2 * time.Second)
	fastTicker := time.NewTicker(10 * time.Millisecond)
	running := true
	for running {
		select {
		case <-t.stop:
			ticker.Stop()
			running = false
		case <-ticker.C:
			t.tick <- true
		case <-fastTicker.C:
			t.tick <- false
		}
	}
}
