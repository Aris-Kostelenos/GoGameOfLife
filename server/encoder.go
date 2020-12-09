package main

import (
	"strings"
)

func encoder(height int, width int, world [][]uint8) string {
	var s string
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
				world[i][j] = 1
			}
			x++
		}
	}
	return world
}
