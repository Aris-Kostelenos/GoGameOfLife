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
	t    Ticker
	quit bool
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

func (client *Client) getWorld(server *rpc.Client) (world [][]uint8, turn int) {
	args := new(stubs.Default)
	reply := new(stubs.World)
	err := server.Call(stubs.GetWorld, args, reply)
	if err != nil {
		fmt.Println("err", err)
	}
	return reply.World, reply.Turn
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
	err := server.Call(stubs.Kill, args, reply)
	if err != nil {
		panic(err)
	}
	return reply.Turn
}

func (client *Client) getAlive(p Params, server *rpc.Client, events chan<- Event) (turn int) {
	args := new(stubs.Default)
	reply := new(stubs.Alive)
	server.Call(stubs.GetNumAlive, args, reply)
	events <- AliveCellsCount{reply.Turn, reply.Num}
	return reply.Turn
}

func (client *Client) checkDone(server *rpc.Client) (done bool) {
	args := new(stubs.Default)
	reply := new(stubs.Done)
	err := server.Call(stubs.CheckDone, args, reply)
	if err != nil {
		panic(err)
	}
	return reply.Done
}

func (client *Client) handleEvents(c clientChannels, p Params, server *rpc.Client) (turn int) {
	turn = 0
	client.quit = false
	for running := true; running; {
		select {
		case tick := <-client.t.tick:
			if tick == true {
				turn = client.getAlive(p, server, c.events)
			} else if tick == false {
				if client.checkDone(server) {
					client.t.stop <- true
					world, currentTurn := client.getWorld(server)
					turn = currentTurn
					saveWorld(c, p, world, turn)
					alive := extractAlive(world)
					c.events <- FinalTurnComplete{turn, alive}
					running = false
				}
			}
		case key := <-c.keyPresses:
			switch key {
			case 's':
				world, currentTurn := client.getWorld(server)
				turn = currentTurn
				saveWorld(c, p, world, turn)
			case 'q':
				client.quit = true
				running = false
			case 'p':
				// tell the server to pause
				turn = client.pauseServer(server)
				fmt.Printf("Paused. %v turns complete\n", turn)
				// wait for resume keypress
				var nextKey rune
				for nextKey != 'p' {
					nextKey = <-c.keyPresses
				}
				// tell the server to resume
				client.pauseServer(server)
				fmt.Println("Continuing...")
			case 'k':
				world, currentTurn := client.getWorld(server)
				turn = currentTurn
				saveWorld(c, p, world, turn)
				running = false
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

	if !client.quit {
		client.killServer(server)
	}
	server.Close()
	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle
	c.events <- StateChange{turn, Quitting}
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
