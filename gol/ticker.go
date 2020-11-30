package gol

import (
	"time"
)

type state struct {
	turns         chan int
	previousWorld chan [][]uint8
	stop          chan bool
	mutex         chan bool
}

func startTicker(events chan<- Event, state state) {
	ticker := time.NewTicker(2 * time.Second)
	turn := 0
	var prevWorld [][]uint8
	x := true
	for x {
		select {
		case <-state.stop:
			ticker.Stop()
			x = false
		case value := <-state.turns:

			turn = value + 1
			prevWorld = <-state.previousWorld

		case <-ticker.C:
			state.mutex <- true
			//fmt.Println("hi!")
			alive := len(calculateAliveCells(prevWorld))
			events <- AliveCellsCount{turn, alive}
			//fmt.Println("hi again!")
			state.mutex <- false
		default:
			break
		}
	}
}
