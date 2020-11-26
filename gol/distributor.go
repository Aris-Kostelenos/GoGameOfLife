package gol

import (
	"fmt"

	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioInput    <-chan uint8
	ioOutput   chan<- uint8
}

type workerParams struct {
	id                       int
	events                   chan<- Event
	imagePartWidthStartpoint int
	imagePartWidth           int
	imagePartHeight          int
	turns                    int
	threads                  int
	//	answer          chan<- [][]uint8
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// TODO: Create a 2D slice to store the world.
	c.ioCommand <- ioInput
	c.ioFilename <- fmt.Sprintf("%vx%v", p.ImageWidth, p.ImageHeight)

	grid := make([][]uint8, p.ImageHeight)

	for row := 0; row < p.ImageHeight; row++ {
		grid[row] = make([]uint8, p.ImageWidth)
		for cell := 0; cell < p.ImageWidth; cell++ {
			grid[row][cell] = <-c.ioInput
			c.events <- CellFlipped{0, util.Cell{X: cell, Y: row}}
		}
	}

	// TODO: For all initially alive cells send a CellFlipped Event.

	turn := 0

	// TODO: Execute all turns of the Game of Life.
	ap := workerParams{0, c.events, 0, p.ImageWidth, p.ImageHeight, p.Turns, p.Threads}

	wp := workerParams{0, c.events, 0, 0, p.ImageHeight, p.Turns, p.Threads}

	workerChannels := make([]chan uint8, 2*p.Threads+1)
	for i := range workerChannels {
		workerChannels[i] = make(chan uint8, p.ImageHeight)
	}
	statusChannels := make([]chan bool, p.Threads)
	for i := range statusChannels {
		statusChannels[i] = make(chan bool)
	}

	wp.imagePartWidth = p.ImageWidth / p.Threads

	if p.ImageWidth%p.Threads == 0 {

		// basically for each go routine we need two channels where it is going to send the edges of the results it gets

		for i := 0; i < p.Threads; i++ {
			wp.id = i
			wp.imagePartWidthStartpoint = i * wp.imagePartWidth
			go workerGoroutine(wp, grid, workerChannels, statusChannels)
		}
	} else {
		fmt.Println("uuuhhh fuk...........")
	}

	for turn = 0; turn < p.Turns; turn++ {
		for i := 0; i < p.Threads; i++ {
			h := <-statusChannels[i]
			if h == true {
			}
		}
		c.events <- TurnComplete{turn}

	}

	/*
		for turn = 0; turn < p.Turns; turn++ {
			grid = calculateNextState(ap, grid)
			c.events <- TurnComplete{turn}
		}
	*/
	for i := 0; i < p.Threads; i++ {
		h := <-statusChannels[i]
		if h == false {
		}
	}

	c.events <- FinalTurnComplete{turn, calculateAliveCells(ap, grid)}

	// TODO: Send correct Events when required, e.g. CellFlipped, TurnComplete and FinalTurnComplete.
	//		 See event.go for a list of all events.

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
