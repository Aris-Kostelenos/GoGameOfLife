package main

import (
	"fmt"
	"server/util"
	"sync"

	"github.com/ChrisGora/semaphore"
)

// Distributor struct
type Distributor struct {
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
func (d *Distributor) makeNextWorld() {
	d.nextWorld = make([][]uint8, d.imageHeight)
	for row := 0; row < d.imageHeight; row++ {
		d.nextWorld[row] = make([]uint8, d.imageWidth)
	}
}

// creates workers and starts their goroutines
func (d *Distributor) makeWorkers() {
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
func (d *Distributor) getAliveCells() []util.Cell {
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
func (d *Distributor) run() {

	fmt.Println("distributor started")
	// initialise the state
	d.makeNextWorld()
	d.makeWorkers()

	// run the game of life
	fmt.Println("total turns:", d.numOfTurns)
	//fmt.Println("height ", d.imageHeight)
	//fmt.Println("Width ", d.imageWidth)
	//fmt.Println("prevWorld height", len(d.prevWorld))
	//fmt.Println("prevWorld width", len(d.prevWorld[0]))
	for d.currentTurn = 0; d.currentTurn < d.numOfTurns && !d.quit; d.currentTurn++ {
		//fmt.Println("a")
		// wait for all workers to complete this turn
		for _, w := range d.workers {
			w.space.Wait()
		}
		//fmt.Println("b")
		// swap the previous and next grids
		d.mutex.Lock()
		temp := d.prevWorld
		d.prevWorld = d.nextWorld
		d.nextWorld = temp
		d.mutex.Unlock()

		//fmt.Println("c")
		/*
			switch {
			case <-d.paused:
				fmt.Println("paused")
				// pause the workers
				<-d.paused
				// resume the workers
			default:
				fmt.Println("not paused")
			}
		*/
		//fmt.Println("d")
		// order the workers to start the next turn and notify the ticker
		for i := 0; i < d.threads && d.quit == false; i++ {
			d.workers[i].work.Post()
		}
		//fmt.Println("e")
	}
	fmt.Println("done")
	// notify that the end is complete
}
