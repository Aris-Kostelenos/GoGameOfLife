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
/*
func (s *Server) StartGoL16(args stubs.Start16, reply *stubs.Default) error {
	fmt.Println("starting new GoL16")
	if s.inProgress {
		return errors.New("Simulation already in progress")
	}
	// convert the world from an array to a slice
	worldSlice := make([][]uint8, 16)
	for row := 0; row < 16; row++ {
		worldSlice[row] = make([]uint8, 16)
		for col := 0; col < 16; col++ {
			worldSlice[row][col] = args.World[row][col]
		}
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

*/

//StartGoLGeneric starts processing -ur mom- a world of any size given to it as a string
func (s *Server) StartGoLGeneric(args stubs.StartGeneric, reply *stubs.Default) error {
	fmt.Println("starting new GoLGeneric")
	if s.inProgress {
		return errors.New("Simulation already in progress")
	}
	WorldSlice := decoder(args.Height, args.Width, args.World)
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

// StartGoL16 starts processing a 16x16 Game of Life simulation
func (s *Server) StartGoL16(args stubs.Start16, reply *stubs.Default) error {
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

// GetWorldGeneric returns the latest state of a world
func (s *Server) GetWorldGeneric(args stubs.Default, reply *stubs.WorldGeneric) error {
	s.distributor.mutex.Lock()
	reply.Turn = s.distributor.currentTurn
	reply.World = encoder(s.distributor.imageHeight, s.distributor.imageWidth, s.distributor.prevWorld)
	reply.Height = s.distributor.imageHeight
	reply.Width = s.distributor.imageWidth
	s.distributor.mutex.Unlock()
	return nil
}

/*
func (s *Server) GetWorld16(args stubs.Default, reply *stubs.World16) error {
	if s.distributor.imageHeight != 16 || s.distributor.imageWidth != 16 {
		message := fmt.Sprintf("Simulated world is %vx%v, not 16x16", s.distributor.imageWidth, s.distributor.imageHeight)
		return errors.New(message)
	}
	s.distributor.mutex.Lock()
	reply.Turn = s.distributor.currentTurn
	reply.World = [16][16]uint8{}
	for row := 0; row < 16; row++ {
		reply.World[row] = [16]uint8{}
		for col := 0; col < 16; col++ {
			reply.World[row][col] = s.distributor.prevWorld[row][col]
		}
	}
	s.distributor.mutex.Unlock()
	return nil
}
*/

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
	fmt.Println("pausing")
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
