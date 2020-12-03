package gol

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

func makeNextWorld(height int, width int) [][]uint8 {
	nextWorld := make([][]uint8, height)
	for row := 0; row < height; row++ {
		nextWorld[row] = make([]uint8, width)
	}
	return nextWorld
}

func makeWorkers(numOfWorkers int, rowsPerSlice int, extra int, wc workerChannels, wp workerParams) {
	for i := 0; i < numOfWorkers; i++ {
		wc.syncChan[i] = make(chan int)
		wc.confChan[i] = make(chan bool)
		workerRows := rowsPerSlice
		if extra > 0 {
			workerRows++
			extra--
		}
		wp.id = i
		wp.imagePartHeight = workerRows
		go workerGoroutine(wp, wc)
		wp.offset += workerRows
	}
}

func writePgmImage(imageHeight int, imageWidth int, world [][]byte, filename string) {
	_ = os.Mkdir("out", os.ModePerm)
	_ = os.Chdir("out")

	file, _ := os.Create(filename)
	//check(ioError)
	defer file.Close()

	_, _ = file.WriteString("P5\n")
	//_, _ = file.WriteString("# PGM file writer by pnmmodules (https://github.com/owainkenwayucl/pnmmodules).\n")
	_, _ = file.WriteString(strconv.Itoa(imageWidth))
	_, _ = file.WriteString(" ")
	_, _ = file.WriteString(strconv.Itoa(imageHeight))
	_, _ = file.WriteString("\n")
	_, _ = file.WriteString(strconv.Itoa(255))
	_, _ = file.WriteString("\n")

	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			_, _ = file.Write([]byte{world[y][x]})
			//check(ioError)
		}
	}

	//ioError = file.Sync()
	//check(ioError)

	fmt.Println("File", filename, "output done!")
}

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

func runGol(p Params, wc workerChannels, c distributorChannels, t *Ticker, prevWorld [][]uint8, nextWorld [][]uint8) int {
	var turn int
	quit := false
	for turn = 0; turn < p.Turns && quit == false; turn++ {

		// wait for all workers to complete this turn
		for i := 0; i < p.Threads; i++ {
			x := <-wc.syncChan[i]
			if x != turn {
				log.Fatal("Thread out of sync")
			}
		}
		c.events <- TurnComplete{turn}

		// swap the previous and next grids
		t.mutex.Lock()
		prevWorld = nextWorld
		nextWorld = makeNextWorld(p.ImageHeight, p.ImageWidth)
		t.mutex.Unlock()

		// handle key presses
		quit = handleKeyPresses(c, p, turn, prevWorld)

		// order the workers to start the next turn and notify the ticker
		for i := 0; i < p.Threads && quit == false; i++ {
			wc.confChan[i] <- true
		}
		t.turns <- turn
	}
	return turn
}

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

	makeWorkers(p.Threads, rowsPerSlice, extra, wc, wp)

	// create a ticker
	t := Ticker{}
	t.turns = make(chan int)
	t.stop = make(chan bool)
	t.prevWorld = &prevWorld
	go t.startTicker(c.events)

	// run the game of life
	// turn := runGol(p, wc, c, &t, prevWorld, nextWorld)
	var turn int
	quit := false
	for turn = 0; turn < p.Turns && quit == false; turn++ {

		// wait for all workers to complete this turn
		for i := 0; i < p.Threads; i++ {
			x := <-wc.syncChan[i]
			if x != turn {
				log.Fatal("Thread out of sync")
			}
		}
		c.events <- TurnComplete{turn}

		// swap the previous and next grids
		t.mutex.Lock()
		prevWorld = nextWorld
		nextWorld = makeNextWorld(p.ImageHeight, p.ImageWidth)
		t.mutex.Unlock()

		// handle key presses
		quit = handleKeyPresses(c, p, turn, prevWorld)

		// order the workers to start the next turn and notify the ticker
		for i := 0; i < p.Threads && quit == false; i++ {
			wc.confChan[i] <- true
		}
		t.turns <- turn
	}

	// end the game of life
	t.stop <- true
	saveWorld(c, p, turn, prevWorld)
	c.events <- FinalTurnComplete{turn, calculateAliveCells(prevWorld)}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
