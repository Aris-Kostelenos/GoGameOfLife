package main

// TODO: rename package to main when separated?

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/rpc"

	"server/stubs"
)

// Server is the interface for the server-side GoL engine
type Server struct {
	inProgress  bool
	distributor Distributor
}

// StartGoL16 starts processing a 16x16 Game of Life simulation
func (s *Server) StartGoL16(args stubs.Start16, reply *stubs.Default) error {
	if s.inProgress {
		return errors.New("Simulation already in progress")
	}
	// convert the world from an array to a slice
	worldSlice := make([][]uint8, 16)
	for row := 0; row < 16; row++ {
		copy(worldSlice[row], args.World[row][:15])
	}
	// start the distributor
	s.distributor = Distributor{
		currentTurn: 0,
		numOfTurns:  args.Turns,
		threads:     args.Threads,
		imageWidth:  16,
		imageHeight: 16,
		prevWorld:   worldSlice,
		paused:      make(chan bool),
	}
	go s.distributor.run()
	s.inProgress = true
	return nil
}

// StartGoL64 starts processing a 64x64 Game of Life simulation
func (s *Server) StartGoL64(args stubs.Start64, reply *stubs.Default) error {
	return nil
}

// StartGoL512 starts processing a 512x512 Game of Life simulation
func (s *Server) StartGoL512(args stubs.Start512, reply *stubs.Default) error {
	return nil
}

// GetWorld16 returns the latest state of a 16x16 world
func (s *Server) GetWorld16(args stubs.Default, reply *stubs.World16) error {
	if s.distributor.imageHeight != 16 || s.distributor.imageWidth != 16 {
		message := fmt.Sprintf("Simulated world is %vx%v, not 16x16", s.distributor.imageWidth, s.distributor.imageHeight)
		return errors.New(message)
	}
	s.distributor.mutex.Lock()
	reply.Turn = s.distributor.currentTurn
	reply.World = [16][16]uint8{}
	for row := 0; row < 16; row++ {
		for col := 0; col < 16; col++ {
			reply.World[row][col] = s.distributor.prevWorld[row][col]
		}
	}
	s.distributor.mutex.Unlock()
	return nil
}

// Connect returns the necessary information for a client to start communicating with the server
func (s *Server) Connect(args stubs.Default, reply *stubs.Status) error {
	reply.Running = s.inProgress
	if s.inProgress {
		reply.CurrentTurn = s.distributor.currentTurn
		reply.NumOfTurns = s.distributor.numOfTurns
		reply.Width = s.distributor.imageWidth
		reply.Height = s.distributor.imageHeight
	}
	return nil
}

// Pause starts/stops the server until further notice
func (s *Server) Pause(args stubs.Default, reply *stubs.Turn) error {
	s.distributor.mutex.Lock()
	s.distributor.paused <- true
	reply.Turn = s.distributor.currentTurn
	s.distributor.mutex.Unlock()
	return nil
}

// Kill shuts down the server
func (s *Server) Kill(args stubs.Default, reply *stubs.Turn) error {
	if s.distributor.quit || s.distributor.currentTurn > s.distributor.numOfTurns {
		return errors.New("The engine has already been quit")
	}
	s.distributor.mutex.Lock()
	s.distributor.quit = true
	reply.Turn = s.distributor.currentTurn
	s.distributor.mutex.Unlock()
	return nil
}

// GetNumAlive returns the number of alive cells and current turn
func (s *Server) GetNumAlive(args stubs.Default, reply *stubs.Alive) error {
	s.distributor.mutex.Lock()
	reply.Num = len(s.distributor.getAliveCells())
	reply.Turn = s.distributor.currentTurn
	s.distributor.mutex.Unlock()
	return nil
}

func main() {
	// parse compiler flags
	port := flag.String("this", "8030", "Port for this service to listen on")
	flag.Parse()
	// register the interface
	rpc.Register(new(Server))
	// listen for calls
	active := true
	for active {
		fmt.Println("listening...")
		listener, err := net.Listen("tcp", ":"+*port)
		if err != nil {
			panic(err)
		}
		defer listener.Close()
		// accept a listener
		fmt.Println("...listener received")
		rpc.Accept(listener)
	}
}
