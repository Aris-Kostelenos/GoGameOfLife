package gol

import (
	"fmt"

	"uk.ac.bris.cs/gameoflife/util"
)

// The calculateNextState function is a copy of the function of the first lab with minimal changes.
// The most notable is the +2 next to p.imagePartHeight
// This function receives and outputs 1 extra row at the top and one at the bottom than it would have needed.
// It can be optimised in that regard.
// Additionally it makes a new board for every turn , while it could instead edit the board the worker is given.

func calculateNextState(p workerParams, world [][]uint8, events chan<- Event) [][]uint8 {
	world1 := make([][]uint8, p.imagePartHeight+2)
	for row := 0; row < p.imagePartHeight+2; row++ {
		world1[row] = make([]uint8, p.imagePartWidth)
	}
	for row := range world {
		for cell := range world[row] {
			var left int
			var right int
			var top int
			var bottom int

			if cell == 0 {
				left = p.imagePartWidth - 1
			} else {
				left = cell - 1
			}
			if row == 0 {
				top = p.imagePartHeight + 1
			} else {
				top = row - 1
			}
			if cell == p.imagePartWidth-1 {
				right = 0
			} else {
				right = cell + 1
			}
			if row == p.imagePartHeight+1 {
				bottom = 0
			} else {
				bottom = row + 1
			}
			var x uint8

			x += world[top][left] / 255
			x += world[top][cell] / 255
			x += world[top][right] / 255
			x += world[row][left] / 255
			x += world[row][right] / 255
			x += world[bottom][left] / 255
			x += world[bottom][cell] / 255
			x += world[bottom][right] / 255

			if world[row][cell] == 0 {
				if x == 3 {
					world1[row][cell] = 255
					events <- CellFlipped{255, util.Cell{X: row, Y: cell}}
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

func workerGoroutine(p workerParams, c workerChannels) {

	//*(p.prevWorld)[row][cell]

	/*
		if row >= 0 && row < len(matrix) {
			return matrix[row][cell]
		} else if row == -1 {
			return matrix[len(matrix)-1][cell]
		} else if row == len(matrix) {
			return matrix[0][cell]
		} else {
			return 0
		}

	*/

	//makes a new grid
	gridPart := make([][]uint8, p.imagePartHeight+2)
	for row := 0; row < p.imagePartHeight+2; row++ {
		gridPart[row] = make([]uint8, p.imagePartWidth)
	}

	// main loop that runs for all the turns
	for turns := 0; turns < p.turns; turns++ {

		//every turn it reads and copies the part of the prevWorld that is relevant to it.
		//the +2 here again is because we need the info of the rows exactly above and below the rows the worker processes.

		for row := 0; row < p.imagePartHeight+2; row++ {
			for cell := 0; cell < p.imagePartWidth; cell++ {
				actualRow := p.offset - 1 + row
				//gridPart[row][cell] = (*p.prevWorld)[p.offset-1+row][cell]
				if actualRow >= 0 && actualRow < p.imageHeight-1 {
					gridPart[row][cell] = (*p.prevWorld)[actualRow][cell]
				} else if actualRow == -1 {
					gridPart[row][cell] = (*p.prevWorld)[p.imageHeight-1]
				} else if actualRow == p.imageHeight {
					gridPart[row][cell] = (*p.prevWorld)[0][cell]
				} else {
					return 0
					fmt.Println("error!")
				}
			}
		}

		//it then sends its board to the above function to do the calculations for one turn
		gridPart = calculateNextState(p, gridPart, c.events)

		//even though the local grid is 2 rows bigger, the top and bottom row are ommited when writing back to nextWorld

		for row := 0; row < p.imagePartHeight; row++ {
			for cell := 0; cell < p.imagePartWidth; cell++ {
				(*p.nextWorld)[row+p.offset][cell] = gridPart[row+1][cell]
			}
		}

		//the worker sends the turn it is in as a signal to distributor to transfer the data of newWorld to oldWorld
		c.syncChan[p.id] <- turns

		//when the distributor is done it sends a bool value to the following channel. These channels act like a mutex lock basically.
		<-c.confChan[p.id]

	}
}

//the function from the first coursework with basically no differences.
func calculateAliveCells(world [][]uint8) []util.Cell {
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
