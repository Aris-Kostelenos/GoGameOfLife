package gol

import "time"

func startTicker(events chan<- Event, prevWorld *[][]uint8, turns chan int, stop chan bool) {
	ticker := time.NewTicker(2 * time.Second)
	turn := 0
	for {
		select {
		case <-stop:
			ticker.Stop()
		case value := <-turns:
			turn = value
		case <-ticker.C:
			pause <- true
			alive := len(calculateAliveCells(*prevWorld))
			events <- AliveCellsCount{turn, alive}
			pause <- false
		default:
			break
		}
	}
}
