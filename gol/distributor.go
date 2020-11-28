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

//this function includes some logic that gives the bottom row (len(matrix-1)) when asked for row -1 and gives the first row (0) when asked for row (len(matrix))
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
			c.events <- CellFlipped{0, util.Cell{X: cell, Y: row}}
		}
	}
	immPrevWorld := makeImmutableMatrix(prevWorld)

	// make a 2D grid for the next state of the world
	nextWorld := make([][]uint8, p.ImageHeight)
	for row := 0; row < p.ImageHeight; row++ {
		nextWorld[row] = make([]uint8, p.ImageWidth)
	}

	// determine how to allocate rows to workers
	rowsPerSlice := p.ImageHeight / p.Threads
	fmt.Println("->", rowsPerSlice)
	extra := p.ImageHeight % p.Threads

	//syncChan and confChan are used for synchronisation of the distributor and workers.
	syncChan := make([]chan int, p.Threads)
	confChan := make([]chan bool, p.Threads)

	//wp is defined outside the loop it is passed by value to each worker.
	wp := workerParams{0, c.events, 0, p.ImageWidth, 0, p.Turns, p.Threads}

	for i := 0; i < p.Threads; i++ {

		//since we iterate over p.Threads we may as well initialise the channels.
		syncChan[i] = make(chan int)
		confChan[i] = make(chan bool)

		workerRows := rowsPerSlice
		if extra > 0 {
			workerRows++
			extra--
		}
		// TODO: revise workerParams

		//id is literally the number of the channel counting from 0.
		wp.id = i
		wp.imagePartHeight = workerRows
		//TODO: make the workers and distributor communicate via channels instead of reading and writing to common matrices.
		go workerGoroutine(wp, immPrevWorld, nextWorld, syncChan, confChan)

		//the offset for the next worker is defined as the previous offset plus the number of rows of the previous worker
		wp.imagePartHeightStartpoint += workerRows

	}

	var turn int
	// run the game of life
	for turn = 0; turn < p.Turns; turn++ {
		//receive a message from every thread sayng they are done with the turn.
		for i := 0; i < p.Threads; i++ {
			x := <-syncChan[i]
			if x != turn {
				//TODO: send an error
			} else {
				//
			}
		}

		//transfer the state of each cell of nextWorld to prevWorld. It is important that this is done because prevWorld is the only pointer that the workers have
		//in order to find out the previous turn. Doing prevWorld = nextWorld would change the pointer for distributor but leave the workers in the dark
		//a channel could give the workers the pointer to the new matrix for a possibly speedier approach

		c.events <- TurnComplete{turn}
		for row := range prevWorld {
			for cell := range prevWorld[0] {
				prevWorld[row][cell] = nextWorld[row][cell]
			}
		}

		//after prevWorld is updated the workers are notified to continue their work.
		for i := 0; i < p.Threads; i++ {
			confChan[i] <- true
		}
	}

	ap := workerParams{0, c.events, 0, p.ImageWidth, p.ImageHeight, p.Turns, p.Threads}

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
