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
	id              int
	offset          int
	imagePartWidth  int
	imagePartHeight int
	imageWidth      int
	imageHeight     int
	turns           int
	threads         int
	prevWorld       *[][]uint8
	nextWorld       *[][]uint8
}

type workerChannels struct {
	events   chan<- Event
	syncChan []chan int
	confChan []chan bool
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

func makeNextWorld(height int, width int) [][]uint8 {
	nextWorld := make([][]uint8, height)
	for row := 0; row < height; row++ {
		nextWorld[row] = make([]uint8, width)
	}
	return nextWorld
}

// func makeWorkers() {
// 	for i := 0; i < p.Threads; i++ {

// 		//since we iterate over p.Threads we may as well initialise the channels.
// 		wc.syncChan[i] = make(chan int)
// 		wc.confChan[i] = make(chan bool)

// 		workerRows := rowsPerSlice
// 		if extra > 0 {
// 			workerRows++
// 			extra--
// 		}
// 		// TODO: revise workerParams

// 		//id is literally the number of the channel counting from 0.
// 		wp.id = i
// 		wp.imagePartHeight = workerRows
// 		//TODO: make the workers and distributor communicate via channels instead of reading and writing to common matrices.
// 		go workerGoroutine(wp, wc)

// 		//the offset for the next worker is defined as the previous offset plus the number of rows of the previous worker
// 		wp.offset += workerRows

// 	}
// }

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// start reading data from the file
	c.ioCommand <- ioInput
	c.ioFilename <- fmt.Sprintf("%vx%v", p.ImageWidth, p.ImageHeight)

	// make a 2D grid for the previous and next state of the world
	prevWorld := makePrevWorld(p.ImageHeight, p.ImageWidth, c)
	nextWorld := makeNextWorld(p.ImageHeight, p.ImageWidth)

	// make a struct for worker channels
	wc := workerChannels{}
	wc.events = c.events
	wc.syncChan = make([]chan int, p.Threads)
	wc.confChan = make([]chan bool, p.Threads)

	// wp is defined outside the loop it is passed by value to each worker.
	wp := workerParams{0, 0, p.ImageWidth, 0, p.ImageWidth, p.ImageHeight, p.Turns, p.Threads, &prevWorld, &nextWorld}

	// determine the number of rows to be allocated to each worker
	rowsPerSlice := p.ImageHeight / p.Threads
	extra := p.ImageHeight % p.Threads

	for i := 0; i < p.Threads; i++ {

		//since we iterate over p.Threads we may as well initialise the channels.
		wc.syncChan[i] = make(chan int)
		wc.confChan[i] = make(chan bool)

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
		go workerGoroutine(wp, wc)

		//the offset for the next worker is defined as the previous offset plus the number of rows of the previous worker
		wp.offset += workerRows

	}

	ds := state{}
	ds.turns = make(chan int)
	ds.stop = make(chan bool)
	ds.previousWorld = make(chan [][]uint8)
	ds.mutex = make(chan bool)
	go startTicker(c.events, ds)

	var turn int
	// run the game of life
	for turn = 0; turn < p.Turns; turn++ {

		//receive a message from every thread sayng they are done with the turn.
		for i := 0; i < p.Threads; i++ {
			x := <-wc.syncChan[i]
			if x != turn {
				//TODO: send an error
			}
		}

		c.events <- TurnComplete{turn}

		prevWorld = nextWorld
		nextWorld = makeNextWorld(p.ImageHeight, p.ImageWidth)

		for i := 0; i < p.Threads; i++ {
			wc.confChan[i] <- true
		}

		select {
		case x := <-ds.mutex:
			if x == true {
				//fmt.Println("yin")
				<-ds.mutex
				//fmt.Println("yan")
			}
		default:
			break
		}

		// update the ticker
		ds.turns <- turn
		ds.previousWorld <- prevWorld
	}

	ds.stop <- true
	c.events <- FinalTurnComplete{turn, calculateAliveCells(prevWorld)}

	// TODO: Send correct Events when required, e.g. CellFlipped, TurnComplete and FinalTurnComplete.
	//		 See event.go for a list of all events.

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
