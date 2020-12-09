package main

import (
	"fmt"
	"strings"
)

func encoder(height int, width int, world [][]uint8) string {
	var s string
	if height%2 == 0 {
		for i := 0; i < height; i += 2 {
			for j := 0; j < width; j++ {
				//s += strconv.Itoa(int(world[i][j])) + ","
				if world[i][j] == 0 && world[i+1][j] == 0 {
					s += "0"
				} else if world[i][j] == 0 && world[i+1][j] == 1 {
					s += "1"
				} else if world[i][j] == 1 && world[i+1][j] == 0 {
					s += "2"
				} else {
					s += "3"
				}

			}
		}

	} else {
		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {
				//s += strconv.Itoa(int(world[i][j])) + ","
				if world[i][j] == 0 {
					s += "0"
				} else {
					s += "1"
				}

			}
		}
		return s
	}

}

func decoder(height int, width int, s string) [][]uint8 {
	x := 0
	world := make([][]uint8, height)
	strings.Split(s, ",")
	for i := 0; i < height; i++ {
		world[i] = make([]uint8, width)
		for j := 0; j < width; j++ {
			if s[x] == '0' {
				world[i][j] = 0
			} else {
				world[i][j] = 255
			}
			x++
		}
	}
	return world
}

func main() {
	// make cock array
	cock := make([][]uint8, 3)
	cock[0] = make([]uint8, 3)
	cock[1] = make([]uint8, 3)
	cock[2] = make([]uint8, 3)
	cock[0][0] = 0
	cock[0][1] = 255
	cock[0][2] = 0
	cock[1][0] = 0
	cock[1][1] = 255
	cock[1][2] = 0
	cock[2][0] = 255
	cock[2][1] = 255
	cock[2][2] = 255
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Println(cock[i][j])
		}
	}
	penis := encoder(3, 3, cock)
	fmt.Println(penis)
	dicc := decoder(3, 3, penis)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Println(dicc[i][j])
		}
	}

	//255

}
