package gol

import "time"

func startTicker(events chan<- Event, prevWorld [][]uint8, turns chan int) {
	c := time.Tick(2 * time.Second)
	turn := 0
	for range c {
		select {
		case value := <-turns:
			turn = value
		default:
			break
		}
		alive := len(calculateAliveCells(prevWorld))
		events <- AliveCellsCount{turn, alive}
	}
}
