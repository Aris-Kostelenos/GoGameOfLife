package gol

import (
	"uk.ac.bris.cs/gameoflife/util"
)

func calculateNextState(p workerParams, world [][]uint8) [][]uint8 {
	world1 := make([][]uint8, p.imagePartHeight+2)
	for z := 0; z < p.imagePartHeight+2; z++ {
		world1[z] = make([]uint8, p.imagePartWidth)
	}
	for row := range world {
		for cell := range world[row] {
			var a int
			var b int
			var c int
			var d int

			if cell == 0 {
				a = p.imagePartWidth - 1
			} else {
				a = cell - 1
			}
			if row == 0 {
				c = p.imagePartHeight + 1
			} else {
				c = row - 1
			}
			if cell == p.imagePartWidth-1 {
				b = 0
			} else {
				b = cell + 1
			}
			if row == p.imagePartHeight+1 {
				d = 0
			} else {
				d = row + 1
			}
			var x uint8
			x += world[c][a] / 255
			x += world[c][cell] / 255
			x += world[row][a] / 255
			x += world[c][b] / 255
			x += world[row][b] / 255
			x += world[d][a] / 255
			x += world[d][cell] / 255
			x += world[d][b] / 255

			if world[row][cell] == 0 {
				if x == 3 {
					world1[row][cell] = 255
					//p.events <- CellFlipped{255, util.Cell{X: row, Y: cell}}
				} else {
					world1[row][cell] = 0
				}
			}

			if world[row][cell] == 255 {
				if x == 3 || x == 2 {
					world1[row][cell] = 255
				} else {
					world1[row][cell] = 0
					//p.events <- CellFlipped{0, util.Cell{X: row, Y: cell}}
				}
			}
		}
	}
	return world1
}

func calculateAliveCells(p workerParams, world [][]uint8) []util.Cell {
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

func workerGoroutine(p workerParams, immPrevWorld func(row, cell int) uint8, nextWorld [][]uint8, turnComplete []chan int) {

	gridPart := make([][]uint8, p.imagePartHeight+2)
	for row := 0; row < p.imagePartHeight+2; row++ {
		gridPart[row] = make([]uint8, p.imagePartWidth)
	}

	for turns := 0; turns < p.turns; turns++ {

		for row := 0; row < p.imagePartHeight+2; row++ {
			for cell := 0; cell < p.imagePartWidth; cell++ {
				gridPart[row][cell] = immPrevWorld(p.imagePartHeightStartpoint-1+row, cell)
			}
		}
		gridPart = calculateNextState(p, gridPart)
		for row := 0; row < p.imagePartHeight; row++ {
			for cell := 0; cell < p.imagePartWidth; cell++ {
				nextWorld[row+p.imagePartHeightStartpoint][cell] = gridPart[row+1][cell]
			}
		}

		turnComplete[p.id] <- turns
		<-turnComplete[p.id]

	}
}
