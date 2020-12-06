package gol

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

type ClientChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioInput    <-chan uint8
	ioOutput   chan<- uint8
	keyPresses <-chan rune
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

func client(p Params, c ClientChannels) {

	// start GOL on the server
	// handle keypresses
	// receive and deal with stuff from the server

	netConn, error1 := net.Dial("tcp", "127.0.0.1:8030")
	if error1 != nil {
		handleError(error1)
	}

	rpcConn, error2 := rpc.Dial("tcp", "127.0.0.1:8030")
	
	if error1 != nil {
		handleError(error1)
	}

	//TODO Start asynchronously reading and displaying messages
	//TODO Start getting and sending user messages.
	go read(&netConn)
	for {

	}
}

func read(conn *net.Conn) {
	//TODO In a continuous loop, read a message from the server and display it.
	reader := bufio.NewReader(*conn)
	for {
		msg, error3 := reader.ReadString('\n')
		if error3 != nil {
			handleError(error3)
		}
		fmt.Print(msg)
	}
}

/*
SERVER TO CLIENT (using net)

ReportAlive() {
	every 2 seconds:
		send numberOfAliveCells to client
}

Finished() {
	after all turns complete:
		send imDone message to client
		send dataFromPrevWorld to client
}
*/

func handleError(err error) {
	fmt.Println(err)
	os.Exit(1)
	// TODO: all
	// Deal with an error event.
}

/*
func client(p, c) {

	conn, error1 := net.Dial("tcp", "127.0.0.1:8030")
	if error1 != nil {
		handleError(error1)
	}
	//TODO Start asynchronously reading and displaying messages
	//TODO Start getting and sending user messages.
	for {

	}
}
*/
