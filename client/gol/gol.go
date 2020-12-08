package gol

import (
	"flag"
	"fmt"
	"net/rpc"

	"uk.ac.bris.cs/gameoflife/stubs"
)

// Params provides the details of how to run the Game of Life and which image to load.
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// creates a grid to represent the current state of the world
func makeWorld16(IoInput chan uint8) [16][16]uint8 {
	world := [16][16]uint8{}
	for row := 0; row < 16; row++ {
		for col := 0; col < 16; col++ {
			world[row][col] = <-IoInput
			// c.events <- CellFlipped{0, util.Cell{X: col, Y: row}} // TODO: remove?
		}
	}
	return world
}

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {

	fmt.Println("started run")
	IoCommand := make(chan ioCommand)
	IoIdle := make(chan bool)
	IoFilename := make(chan string)
	IoInput := make(chan uint8)
	IoOutput := make(chan uint8)

	// read the data from the file to construct a 2D grid
	IoCommand <- ioInput
	IoFilename <- fmt.Sprintf("%vx%v", p.ImageWidth, p.ImageHeight)
	world := makeWorld16(IoInput)

	// parse the command-line flags
	var serverAddress string
	flag.StringVar(&serverAddress, "server", "localhost:8030", "IP:Port string of the server")

	// dial the server
	fmt.Println("dialling server...")
	server, err := rpc.Dial("tcp", serverAddress)
	if err != nil {
		panic(err)
	}
	defer server.Close()
	fmt.Println("...dialled server")

	// start the game of life simulation on the server
	args := stubs.Start16{
		Turns:   p.Turns,
		Threads: p.Threads,
		World:   world,
	}
	reply := new(stubs.ID)
	fmt.Println("calling server...")
	err = server.Call(stubs.StartGoL16, args, reply)
	fmt.Println("...called server")
	if err != nil {
		// try connecting instead
	} else {
		// GoL has begun
		// start handling keypresses, tickers and other events
	}

	clientChannels := clientChannels{
		events,
		IoCommand,
		IoIdle,
		IoFilename,
		IoInput,
		IoOutput,
		keyPresses,
	}
	client := Client{}
	go client.run(p, clientChannels, server)

	ioChannels := ioChannels{
		command:  IoCommand,
		idle:     IoIdle,
		filename: IoFilename,
		output:   IoOutput,
		input:    IoInput,
	}
	go startIo(p, ioChannels)
}