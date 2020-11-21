package gol

import (
	// "fmt"
	"uk.ac.bris.cs/gameoflife/util"
)

type golParams struct {
	events      chan<- Event
	imageWidth  int
	imageHeight int
}

func calculateNextState(p golParams, world [][]uint8) [][]uint8 {
	world1 := make([][]uint8, p.imageWidth)
	for z := 0; z < p.imageWidth; z++ {
		world1[z] = make([]uint8, p.imageHeight)
	}
	for row := range world {
		for cell := range world[row] {
			var a int
			var b int
			var c int
			var d int

			if row == 0 {
				a = p.imageWidth - 1
			} else {
				a = row - 1
			}
			if cell == 0 {
				c = p.imageHeight - 1
			} else {
				c = cell - 1
			}
			if row == p.imageWidth-1 {
				b = 0
			} else {
				b = row + 1
			}
			if cell == p.imageHeight-1 {
				d = 0
			} else {
				d = cell + 1
			}

			x := ((world[a][c] / 255) + (world[a][cell] / 255) + (world[a][d] / 255) + (world[row][c] / 255) + (world[row][d] / 255) + (world[b][c] / 255) + (world[b][cell] / 255) + (world[b][d] / 255))

			if world[row][cell] == 0 {
				if x == 3 {
					world1[row][cell] = 255
					p.events <- CellFlipped{255, util.Cell{X: row, Y: cell}}
				} else {
					world1[row][cell] = 0
				}
			}

			if world[row][cell] == 255 {
				if x == 3 || x == 2 {
					world1[row][cell] = 255
				} else {
					world1[row][cell] = 0
					p.events <- CellFlipped{0, util.Cell{X: row, Y: cell}}
				}
			}
		}
	}
	return world1
}

func calculateAliveCells(p golParams, world [][]uint8) []util.Cell {
	a := make([]util.Cell, 0)
	k := 0
	for row := range world {
		for cell := range world[row] {
			if world[cell][row] == 255 {
				a = append(a, util.Cell{X: row, Y: cell})
				k++
			}
		}
	}
	return a
}
