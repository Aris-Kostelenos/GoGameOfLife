package gol

// Server is the interface for the server-side GoL engine
type Server struct{}

// StartArgs contains params for starting the Game of Life
type StartArgs struct {
    Turns int
    Width int
    Height int
	ClientAddress string
}

// Default args/reply for all methods
type Default struct{}

// ImageReply contains an address to send the world to
type ImageReply struct {
	ClientAddress string
}

// DisconnectArgs contains the address of the client disconnecting
type DisconnectArgs struct {
	ClientAddress string
}

// ConnectArgs contains the address of the client connecting
type ConnectArgs struct {
	ClientAddress string
}

// PauseArgs contains the current turn
type PauseArgs struct {
	Turn int
}

// KillReply contains the address to send the world to
type KillReply struct {
	ClientAddress string
}

// StartGoL starts processing the Game of Life
func (s *Server) StartGoL(args StartArgs, reply *Default) error {
	return nil
}

// GetImage sends prevWorld to the client
func (s *Server) GetImage(args Default, reply *ImageReply) error {
	return nil
}

// Disconnect tells the server to stop sending messages to that client
func (s *Server) Disconnect(args DisconnectArgs, reply *Default) error {
	return nil
}

// Connect tells the server to start communicating with a new client
func (s *Server) Connect(args ConnectArgs, reply *Default) error {
	return nil
}

// Pause starts/stops the server until further notice
func (s *Server) Pause(args Default, *reply PauseArgs) error {
	return nil
}

// Kill shuts down the server
func (s *Server) Kill(args Default, *reply KillReply) error {
	return nil
}

// ReportAlive returns the number of alive cells
func (s *Server) ReportAlive(args, *reply Default) error {
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
        listener, err := net.Listen("tcp", ":"+*port)
        if err != nil {
            panic(err)
        }
        defer listener.Close()
        // accept a listener
        go rpc.Accept(listener)
    }
}
