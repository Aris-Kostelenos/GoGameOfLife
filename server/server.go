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

// StartGoL starts processing a world of any size given to it as a string
func (s *Server) StartGoL(args stubs.StartArgs, reply *stubs.Default) error {
	if s.inProgress {
		return errors.New("Simulation already in progress")
	}
	WorldSlice := args.World
	s.distributor = Distributor{
		currentTurn: 0,
		numOfTurns:  args.Turns,
		threads:     args.Threads,
		imageWidth:  args.Width,
		imageHeight: args.Height,
		prevWorld:   WorldSlice,
		paused:      make(chan bool),
	}
	go s.distributor.run()
	s.inProgress = true
	return nil
}

// GetWorld returns the latest state of a world
func (s *Server) GetWorld(args stubs.Default, reply *stubs.World) error {
	s.distributor.mutex.Lock()
	reply.Turn = s.distributor.currentTurn
	reply.World = s.distributor.prevWorld
	reply.Height = s.distributor.imageHeight
	reply.Width = s.distributor.imageWidth
	s.distributor.mutex.Unlock()
	return nil
}

// Connect returns the necessary information for a client to start communicating with the server
func (s *Server) Connect(args stubs.Default, reply *stubs.Status) error {
	reply.Running = s.inProgress
	if s.inProgress {
		s.distributor.mutex.Lock()
		reply.CurrentTurn = s.distributor.currentTurn
		reply.NumOfTurns = s.distributor.numOfTurns
		reply.Width = s.distributor.imageWidth
		reply.Height = s.distributor.imageHeight
		s.distributor.mutex.Unlock()
	}
	return nil
}

// Pause starts/stops the server until further notice
func (s *Server) Pause(args stubs.Default, reply *stubs.Turn) error {
	reply.Turn = s.distributor.currentTurn
	s.distributor.paused <- true
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
	s.inProgress = false
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

// CheckDone returns true if the server is finished processing the current simulation
func (s *Server) CheckDone(args stubs.Default, reply *stubs.Done) error {
	if !s.inProgress {
		return errors.New("No simulation is in progress")
	}
	s.distributor.mutex.Lock()
	if s.distributor.numOfTurns == s.distributor.currentTurn {
		reply.Done = true
	} else {
		reply.Done = false
	}
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
		rpc.Accept(listener)
	}
}
