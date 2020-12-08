package server

import (
	"sync"

	"github.com/ChrisGora/semaphore"
	"uk.ac.bris.cs/gameoflife/util"
)

type distributor struct {
	currentTurn int
	numOfTurns  int
	threads     int
	imageWidth  int
	imageHeight int
	prevWorld   [][]uint8
	nextWorld   [][]uint8
	workers     []worker
	mutex       sync.Mutex
	paused      chan bool
	quit        bool
}

// creates a grid to represent the next state of the world
func (d *distributor) makeNextWorld() {
	d.nextWorld = make([][]uint8, d.imageHeight)
	for row := 0; row < d.imageHeight; row++ {
		d.nextWorld[row] = make([]uint8, d.imageWidth)
	}
}

// creates workers and starts their goroutines
func (d *distributor) makeWorkers() {
	rowsPerSlice := d.imageHeight / d.threads
	extra := d.imageHeight % d.threads
	startRow := 0
	d.workers = make([]worker, d.threads)
	for i := 0; i < d.threads; i++ {
		// determine the number of rows for this worker
		workerRows := rowsPerSlice
		if extra > 0 {
			workerRows++
			extra--
		}
		// create the worker
		w := worker{}
		w.prevWorld = &d.prevWorld
		w.nextWorld = &d.nextWorld
		w.startRow = startRow
		w.endRow = startRow + workerRows - 1
		w.width = d.imageWidth
		w.work = semaphore.Init(1, 1)
		w.space = semaphore.Init(1, 0)
		go w.processStrip()
		d.workers[i] = w
		// prep for the next iteration
		startRow = w.endRow + 1
	}
}

// returns a slice of the alive cells in prevWorld
func (d *distributor) getAliveCells() []util.Cell {
	alive := make([]util.Cell, 0)
	for row := range d.prevWorld {
		for col := range d.prevWorld[row] {
			if d.prevWorld[row][col] == LIVE {
				alive = append(alive, util.Cell{X: col, Y: row})
			}
		}
	}
	return alive
}

// distributor divides the work between workers and interacts with other goroutines.
func (d *distributor) run() {

	// initialise the state
	d.makeNextWorld()
	d.makeWorkers()

	// run the game of life
	for d.currentTurn = 0; d.currentTurn < d.numOfTurns && !d.quit; d.currentTurn++ {

		// wait for all workers to complete this turn
		for _, w := range d.workers {
			w.space.Wait()
		}

		// swap the previous and next grids
		d.mutex.Lock()
		temp := d.prevWorld
		d.prevWorld = d.nextWorld
		d.nextWorld = temp
		d.mutex.Unlock()

		if <-d.paused {
			// pause the workers
			<-d.paused
			// resume the workers
		}

		// order the workers to start the next turn and notify the ticker
		for i := 0; i < d.threads && d.quit == false; i++ {
			d.workers[i].work.Post()
		}
	}
	// notify that the end is complete
}
