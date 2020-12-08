package gol

import (
	"net/rpc"
	"sync"
	"time"

	"uk.ac.bris.cs/gameoflife/stubs"
)

// Ticker is used to send AliveCellsCount events every 2 seconds
type Ticker struct {
	stop       chan bool
	done       chan int
	server     *rpc.Client
	mutex      sync.Mutex
	totalTurns int
}

func (t *Ticker) startTicker(events chan<- Event) {
	ticker := time.NewTicker(2 * time.Second)
	running := true
	for running {
		select {
		case <-t.stop:
			ticker.Stop()
			running = false
		case <-ticker.C:
			t.mutex.Lock()
			args := new(stubs.Default)
			reply := new(stubs.Alive)
			t.server.Call(stubs.GetNumAlive, args, reply)
			events <- AliveCellsCount{reply.Turn, reply.Num}
			if reply.Turn > t.totalTurns {
				t.done <- reply.Turn
			}
			t.mutex.Unlock()
		}
	}
}
