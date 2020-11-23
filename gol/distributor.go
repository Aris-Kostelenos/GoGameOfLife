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

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// TODO: Create a 2D slice to store the world.
	c.ioCommand <- ioInput
	c.ioFilename <- fmt.Sprintf("%vx%v", p.ImageWidth, p.ImageHeight)

	grid := make([][]uint8, p.ImageHeight)
	for row := 0; row < p.ImageHeight; row++ {
		thisRow := make([]uint8, p.ImageWidth)
		for cell := 0; cell < p.ImageWidth; cell++ {
			value := <-c.ioInput
			thisRow[cell] = value
			// grid[row][cell] = value
			c.events <- CellFlipped{0, util.Cell{X: cell, Y: row}}
		}
		grid[row] = thisRow
	}
	// TODO: For all initially alive cells send a CellFlipped Event.

	turn := 0

	// TODO: Execute all turns of the Game of Life.
	ap := golParams{c.events, p.ImageWidth, p.ImageHeight}
	for turn = 0; turn < p.Turns; turn++ {
		grid = calculateNextState(ap, grid)
		c.events <- TurnComplete{turn}
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
