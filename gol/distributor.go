package gol

import (
	"fmt"
	"log"

	"github.com/ChrisGora/semaphore"
	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioInput    <-chan uint8
	ioOutput   chan<- uint8
	keyPresses <-chan rune
}

// saves the given world as a pgm file
func saveWorld(c distributorChannels, p Params, turn int, world [][]uint8) {
	c.ioCommand <- ioOutput
	outputFilename := fmt.Sprintf("%vx%vx%v", p.ImageWidth, p.ImageHeight, turn)
	c.ioFilename <- outputFilename
	for row := 0; row < p.ImageHeight; row++ {
		for cell := 0; cell < p.ImageWidth; cell++ {
			c.ioOutput <- world[row][cell]
		}
	}
}

// creates a grid to represent the current state of the world
func makePrevWorld(height int, width int, c distributorChannels) [][]uint8 {
	prevWorld := make([][]uint8, height)
	for row := 0; row < height; row++ {
		prevWorld[row] = make([]uint8, width)
		for cell := 0; cell < width; cell++ {
			prevWorld[row][cell] = <-c.ioInput
			c.events <- CellFlipped{0, util.Cell{X: cell, Y: row}}
		}
	}
	return prevWorld
}

// creates a grid to represent the next state of the world
func makeNextWorld(height int, width int) [][]uint8 {
	nextWorld := make([][]uint8, height)
	for row := 0; row < height; row++ {
		nextWorld[row] = make([]uint8, width)
	}
	return nextWorld
}

// creates workers and starts their goroutines
func makeWorkers(p Params, c distributorChannels, prevWorld, nextWorld *[][]uint8) []worker {
	rowsPerSlice := p.ImageHeight / p.Threads
	extra := p.ImageHeight % p.Threads
	startRow := 0
	workers := make([]worker, p.Threads)
	for i := 0; i < p.Threads; i++ {
		// determine the number of rows for this worker
		workerRows := rowsPerSlice
		if extra > 0 {
			workerRows++
			extra--
		}
		// create the worker
		w := worker{}
		w.prevWorld = prevWorld
		w.nextWorld = nextWorld
		w.events = c.events
		w.startRow = startRow
		w.endRow = startRow + workerRows - 1
		w.width = p.ImageWidth
		w.work = semaphore.Init(1, 1)
		w.space = semaphore.Init(1, 0)
		go w.processStrip()
		workers[i] = w
		// prep for the next iteration
		startRow = w.endRow + 1
	}
	return workers
}

// handles keypresses received from SDL
func handleKeyPresses(c distributorChannels, p Params, turn int, prevWorld [][]uint8) bool {
	quit := false
	select {
	case x := <-c.keyPresses:
		switch x {
		case 's':
			saveWorld(c, p, turn, prevWorld)
		case 'q':
			quit = true
		case 'p':
			<-c.keyPresses
		case 'k':
			break
		default:
			log.Fatalf("Unexpected keypress: %v", x)
		}
	default:
		break
	}
	return quit
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// start reading data from the file
	c.ioCommand <- ioInput
	c.ioFilename <- fmt.Sprintf("%vx%v", p.ImageWidth, p.ImageHeight)

	// make a 2D grid for the previous and next state of the world
	prevWorld := makePrevWorld(p.ImageHeight, p.ImageWidth, c)
	nextWorld := makeNextWorld(p.ImageHeight, p.ImageWidth)

	// create the workers and start them off
	workers := makeWorkers(p, c, &prevWorld, &nextWorld)

	// create a ticker
	t := Ticker{}
	t.turns = make(chan int)
	t.stop = make(chan bool)
	t.prevWorld = &prevWorld
	go t.startTicker(c.events)

	// run the game of life
	var turn int
	quit := false
	for turn = 0; turn < p.Turns && quit == false; turn++ {

		// wait for all workers to complete this turn
		for _, w := range workers {
			w.space.Wait()
		}
		c.events <- TurnComplete{turn}

		// swap the previous and next grids
		t.mutex.Lock()
		temp := prevWorld
		prevWorld = nextWorld
		nextWorld = temp
		t.mutex.Unlock()

		// handle key presses
		quit = handleKeyPresses(c, p, turn, prevWorld)

		// order the workers to start the next turn and notify the ticker
		for i := 0; i < p.Threads && quit == false; i++ {
			workers[i].work.Post()
		}
		t.turns <- turn
	}

	// end the game of life
	t.stop <- true
	saveWorld(c, p, turn, prevWorld)
	c.events <- FinalTurnComplete{turn, getAliveCells(prevWorld)}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
