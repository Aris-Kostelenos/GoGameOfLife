package stubs

// StartGoL16 starts a 16x16 Game of Life simulation on the server.
// args = Start16, reply = Default
var StartGoL16 = "Server.StartGoL16"

// StartGoL64 starts a 64x64 Game of Life simulation on the server.
// args = Start64, reply = Default
var StartGoL64 = "Server.StartGoL64"

// StartGoL512 starts a 512x512 Game of Life simulation on the server.
// args = Start512, reply = Default
var StartGoL512 = "Server.StartGoL512"

//StartGoLGeneric starts a generic Game of Life simulation on the server.
// args = StartGeneric, reply = Default
var StartGoLGeneric = "Server.StartGolLGeneric"

// GetWorld16 returns the latest 16x16 world from the server.
// args = Default, reply = World16
var GetWorld16 = "Server.GetWorld16"

// GetWorld64 returns the latest 64x64 world from the server.
// args = Default, reply = World64
var GetWorld64 = "Server.GetWorld64"

// GetWorld512 returns the latest 512x512 world from the server.
// args = Default, reply = World512
var GetWorld512 = "Server.GetWorld512"

//GetWorldGeneric reutnrssda latest nXn world form serverer wbere m is integer greater than zeroooooo.
//args = Default, reply = WorldGeneric
var GetWorldGeneric = "Server.GetWorldGeneric"

// Connect establishes a connection between client and server.
// args = Default, reply = Status
var Connect = "Server.Connect"

// GetNumAlive gives a report on the number of alive cells and the current turn.
// args = Default, reply = Alive
var GetNumAlive = "Server.GetNumAlive"

// Pause stops the server until further notice.
// args = Default, reply = Turn
var Pause = "Server.Pause"

// Kill shuts down the server.
// args = Default, reply = Turn
var Kill = "Server.Kill"

// Start16 contains params for starting a 16x16 Game of Life simulation
type Start16 struct {
	Turns   int
	Threads int
	World   [16][16]uint8
}

// Start64 contains params for starting a 64x64 Game of Life simulation
type Start64 struct {
	Turns   int
	Threads int
	World   [64][64]uint8
}

// Start512 contains params for starting a 512x512 Game of Life simulation
type Start512 struct {
	Turns   int
	Threads int
	World   [512][512]uint8
}

//StartGeneric Hey, aris here, this is the thing
type StartGeneric struct {
	Turns   int
	Threads int
	Height  int
	Width   int
	World   string
}

// Default args/reply for all methods
type Default struct{}

// World16 contains a 16x16 world
type World16 struct {
	World [16][16]uint8
	Turn  int
}

// World64 contains a 64x64 world
type World64 struct {
	World [64][64]uint8
	Turn  int
}

// World512 contains a 512x512 world
type World512 struct {
	World [512][512]uint8
	Turn  int
}

//WorldGeneric containis a string that is a world in string format
type WorldGeneric struct {
	World  string
	Height int
	Width  int
	Turn   int
}

// Turn contains the current turn
type Turn struct {
	Turn int
}

// Alive contains the number of alive cells
type Alive struct {
	Num  int
	Turn int
}

// ID contains the allocated client id
type ID struct {
	ClientID int
}

// Status contains details about the engine's current simulation
// if Running, it will include the state (the following variables)
type Status struct {
	Running     bool
	Width       int
	Height      int
	CurrentTurn int
	NumOfTurns  int
}
