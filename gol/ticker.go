package gol

import (
	"sync"
	"time"
)

// Ticker is used to send AliveCellsCount events every 2 seconds
type Ticker struct {
	turns     chan int
	prevWorld *[][]uint8
	stop      chan bool
	mutex     sync.Mutex
}

func (t *Ticker) startTicker(events chan<- Event) {
	ticker := time.NewTicker(2 * time.Second)
	turn := 0
	running := true
	for running {
		select {
		case <-t.stop:
			ticker.Stop()
			running = false
		case value := <-t.turns:
			turn = value + 1
		case <-ticker.C:
			t.mutex.Lock()
			alive := len(calculateAliveCells(*t.prevWorld))
			events <- AliveCellsCount{turn, alive}
			t.mutex.Unlock()
		default:
			break
		}
	}
}
