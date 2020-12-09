package gol

import (
	"fmt"
	"log"
	"net/rpc"

	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

type clientChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioInput    <-chan uint8
	ioOutput   chan<- uint8
	keyPresses <-chan rune
}

// Client performs server interaction
type Client struct {
	t Ticker
}

func saveWorld(c clientChannels, p Params, world [][]uint8, turn int) {
	c.ioCommand <- ioOutput
	outputFilename := fmt.Sprintf("%vx%vx%v", p.ImageWidth, p.ImageHeight, turn)
	c.ioFilename <- outputFilename
	for row := 0; row < p.ImageHeight; row++ {
		for cell := 0; cell < p.ImageWidth; cell++ {
			c.ioOutput <- world[row][cell]
		}
	}
}

func array16ToSlice(world [16][16]uint8) [][]uint8 {
	sliceWorld := make([][]uint8, 16)
	for row := 0; row < 16; row++ {
		sliceWorld[row] = make([]uint8, 16)
		for col := 0; col < 16; col++ {
			sliceWorld[row][col] = world[row][col]
		}
	}
	return sliceWorld
}

func extractAlive(world [][]uint8) []util.Cell {
	alive := make([]util.Cell, 0)
	for row := range world {
		for col := range world[row] {
			if world[row][col] == 255 {
				alive = append(alive, util.Cell{X: col, Y: row})
			}
		}
	}
	return alive
}

func (client *Client) getWorldGeneric(server *rpc.Client) (world [][]uint8, turn int) {
	args := new(stubs.Default)
	reply := new(stubs.WorldGeneric)
	fmt.Println("getting worlf Genreric")
	err := server.Call(stubs.GetWorldGeneric, args, reply)
	fmt.Println("err", err)
	fmt.Println(reply.World)
	return decoder(reply.Height, reply.Width, reply.World), reply.Turn
}

func (client *Client) getWorld16(server *rpc.Client) (world [][]uint8, turn int) {
	args := new(stubs.Default)
	reply := new(stubs.World16)
	fmt.Println("getting world16")
	err := server.Call(stubs.GetWorld16, args, reply)
	fmt.Println("err:", err)
	return array16ToSlice(reply.World), reply.Turn
}

func (client *Client) pauseServer(server *rpc.Client) (turn int) {
	args := new(stubs.Default)
	reply := new(stubs.Turn)
	server.Call(stubs.Pause, args, reply)
	return reply.Turn
}

func (client *Client) killServer(server *rpc.Client) (turn int) {
	args := new(stubs.Default)
	reply := new(stubs.Turn)
	server.Call(stubs.Kill, args, reply)
	return reply.Turn
}

func (client *Client) getAlive(p Params, server *rpc.Client, events chan<- Event) (done bool) {
	args := new(stubs.Default)
	reply := new(stubs.Alive)
	server.Call(stubs.GetNumAlive, args, reply)
	events <- AliveCellsCount{reply.Turn, reply.Num}
	if reply.Turn >= p.Turns {
		return true
	}
	return false
}

func (client *Client) handleEvents(c clientChannels, p Params, server *rpc.Client) (turn int) {
	turn = 0
	for quit := false; !quit; {
		select {
		case <-client.t.tick:
			done := client.getAlive(p, server, c.events)
			if done {
				world, turn := client.getWorldGeneric(server)
				saveWorld(c, p, world, turn)
				alive := extractAlive(world)
				c.events <- FinalTurnComplete{turn, alive}
				turn = client.killServer(server)
				quit = true
			}
		case key := <-c.keyPresses:
			switch key {
			case 's':
				world, turn := client.getWorldGeneric(server)
				saveWorld(c, p, world, turn)
			case 'q':
				quit = true
			case 'p':
				// tell the server to pause
				turn = client.pauseServer(server)
				fmt.Printf("Paused. %v turns complete\n", turn)
				// wait for resume keypress
				var key rune
				for key != 'p' {
					key = <-c.keyPresses
				}
				// tell the server to resume
				client.pauseServer(server)
				fmt.Println("Continuing...")
			case 'k':
				world, turn := client.getWorldGeneric(server)
				saveWorld(c, p, world, turn)
				turn = client.killServer(server)
				quit = true
			default:
				log.Fatalf("Unexpected keypress: %v", key)
			}
		}
	}
	return turn
}

func (client *Client) run(p Params, c clientChannels, server *rpc.Client) {

	// create a ticker
	client.t = Ticker{}
	client.t.stop = make(chan bool)
	client.t.tick = make(chan bool)
	go client.t.startTicker(c.events)

	// main loop
	turn := client.handleEvents(c, p, server)

	// end the ticker
	client.t.stop <- true
	server.Close()

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
