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
	id                        int
	events                    chan<- Event
	imagePartHeightStartpoint int
	imagePartWidth            int
	imagePartHeight           int
	turns                     int
	threads                   int
}

func makeImmutableMatrix(matrix [][]uint8) func(row, cell int) uint8 {
	return func(row, cell int) uint8 {
		if row >= 0 && row < len(matrix) {
			return matrix[row][cell]
		} else if row == -1 {
			return matrix[len(matrix)-1][cell]
		} else if row == len(matrix) {
			return matrix[0][cell]
		} else {
			return 0
		}
	}
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// TODO: Create a 2D slice to store the world.
	// TODO: For all initially alive cells send a CellFlipped Event.
	// TODO: Execute all turns of the Game of Life.
	c.ioCommand <- ioInput
	c.ioFilename <- fmt.Sprintf("%vx%v", p.ImageWidth, p.ImageHeight)

	// make a 2D grid for the previous state of the world
	prevWorld := make([][]uint8, p.ImageHeight)
	for row := 0; row < p.ImageHeight; row++ {
		prevWorld[row] = make([]uint8, p.ImageWidth)
		for cell := 0; cell < p.ImageWidth; cell++ {
			prevWorld[row][cell] = <-c.ioInput
			c.events <- CellFlipped{0, util.Celld{X: cell, Y: row}}
		}
	}
	immPrevWorld := makeImmtutableMatrix(prevWorld)

	// make a 2D grid for the next state of the world
	nextWorld := make([][]uint8, p.ImageHeight)
	for row := 0; row < p.ImageHeight; row++ {
		nextWorld[row] = make([]uint8, p.ImageWidth)
	}

	// determine how to allocate rows to workers
	rowsPerSlice := p.ImageHeight / p.Threads
	extra := p.ImageHeight % p.Threads

	// run the game of life
	for turn := 0; turn < p.Turns; turn++ {
		for i := 0; i < p.Threads; i++ {
			workerRows := rowsPerSlice
			if extra > 0 {
				workerRows++
				extra--
			}
			// TODO: revise workerParams
			wp := workerParams{i, c.events, 0, 0, p.ImageHeight, p.Turns, p.Threads}
			wp.imagePartHeightStartpoint = i * wp.imagePartHeight
			go workerGoroutine(wp, immPrevWorld, nextWorld)
		}
		c.events <- TurnComplete{turn}
		immPrevWorld = nextWorld
	}
	c.events <- FinalTurnComplete{turn, calculateAliveCells(ap, prevWorld)}

	// TODO: Send correct Events when required, e.g. CellFlipped, TurnComplete and FinalTurnComplete.
	//		 See event.go for a list of all events.

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
