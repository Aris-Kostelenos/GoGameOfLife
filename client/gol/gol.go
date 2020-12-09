package gol

import (
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

func makeWorld(IoInput chan uint8, height int, width int) [][]uint8 {
	world := make([][]uint8, height)
	for row := 0; row < height; row++ {
		world[row] = make([]uint8, width)
		for cell := 0; cell < width; cell++ {
			world[row][cell] = <-IoInput
		}
	}
	return world
}

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {

	IoCommand := make(chan ioCommand)
	IoIdle := make(chan bool)
	IoFilename := make(chan string)
	IoInput := make(chan uint8)
	IoOutput := make(chan uint8)

	ioChannels := ioChannels{
		command:  IoCommand,
		idle:     IoIdle,
		filename: IoFilename,
		output:   IoOutput,
		input:    IoInput,
	}
	go startIo(p, ioChannels)

	// read the data from the file to construct a 2D grid
	IoCommand <- ioInput
	IoFilename <- fmt.Sprintf("%vx%v", p.ImageWidth, p.ImageHeight)

	//world := makeWorld16(IoInput)
	world := makeWorld(IoInput, p.ImageWidth, p.ImageHeight)

	// parse the command-line flags
	serverAddress := "localhost:8030"
	// flag.StringVar(&serverAddress, "server", "localhost:8030", "IP:Port string of the server")

	// dial the server
	server, err := rpc.Dial("tcp", serverAddress)
	if err != nil {
		panic(err)
	}
	// defer server.Close()

	// start the game of life simulation on the server

	//stringWorld := encoder(p.ImageHeight, p.ImageWidth, world)

	args := stubs.StartArgs{
		Turns:   p.Turns,
		Threads: p.Threads,
		Height:  p.ImageHeight,
		Width:   p.ImageWidth,
		World:   world,
	}

	def := new(stubs.Default)
	status := new(stubs.Status)
	server.Call(stubs.Connect, def, status)
	if status.Running {
		killReply := new(stubs.Turn)
		server.Call(stubs.Kill, def, killReply)
	}

	reply := new(stubs.ID)
	err = server.Call(stubs.StartGoL, args, reply)
	if err != nil {
		panic(err)
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

}
