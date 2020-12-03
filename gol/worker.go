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

func calculateNextState(p workerParams, oldWorld [][]uint8, newWorld [][]uint8, topRow []uint8, bottomRow []uint8, events chan<- Event) {

	//fmt.Println("shit : ", len(oldWorld))
	//fmt.Println("cum : ", len(oldWorld[0]))
	//fmt.Println("penis : ", len(newWorld))
	//fmt.Println("cock even : ", len(newWorld[0]))

	for row := 0; row < p.imagePartHeight; row++ {
		for cell := 0; cell < p.imageWidth; cell++ {
			var left int
			var right int
			top := row - 1
			bottom := row + 1
			var x uint8

			if cell == 0 {
				left = p.imagePartWidth - 1
			} else {
				left = cell - 1
			}

			if cell == p.imagePartWidth-1 {
				right = 0
			} else {
				right = cell + 1
			}

			x += oldWorld[row][left] / 255
			x += oldWorld[row][right] / 255

			if row > 0 && row <= p.imagePartHeight-1 {

				x += oldWorld[top][left] / 255
				x += oldWorld[top][cell] / 255
				x += oldWorld[top][right] / 255

			} else if row == 0 {

				x += topRow[left] / 255
				x += topRow[cell] / 255
				x += topRow[right] / 255

			} else {
				fmt.Println("error1!")
			}

			if row >= 0 && row < p.imagePartHeight-1 {

				x += oldWorld[bottom][left] / 255
				x += oldWorld[bottom][cell] / 255
				x += oldWorld[bottom][right] / 255

			} else if row == p.imagePartHeight-1 {

				x += bottomRow[left] / 255
				x += bottomRow[cell] / 255
				x += bottomRow[right] / 255

			} else {
				fmt.Println("error2!")
			}

			if oldWorld[row][cell] == 0 {
				if x == 3 {
					newWorld[row][cell] = 255
					events <- CellFlipped{255, util.Cell{X: row + p.offset, Y: cell}}
				} else {
					newWorld[row][cell] = 0
				}
			}

			if oldWorld[row][cell] == 255 {
				if x == 3 || x == 2 {
					newWorld[row][cell] = 255
				} else {
					newWorld[row][cell] = 0
					events <- CellFlipped{0, util.Cell{X: row + p.offset, Y: cell}}
				}
			}
		}
	}
}

func workerGoroutine(p workerParams, c workerChannels) {

	//makes a new grid
	oldGridPart := make([][]uint8, p.imagePartHeight)
	for row := 0; row < p.imagePartHeight; row++ {
		oldGridPart[row] = make([]uint8, p.imagePartWidth)
		for cell := 0; cell < p.imageWidth; cell++ {
			oldGridPart[row][cell] = (*p.prevWorld)[p.offset+row][cell]
		}
	}

	newGridPart := make([][]uint8, p.imagePartHeight)
	for row := 0; row < p.imagePartHeight; row++ {
		newGridPart[row] = make([]uint8, p.imagePartWidth)

	}

	topRow := make([]uint8, p.imageWidth)
	bottomRow := make([]uint8, p.imageWidth)

	var placeholderGridPart [][]uint8

	// main loop that runs for all the turns
	for turns := 0; turns < p.turns; turns++ {

		//every turn it reads and copies the part of the prevWorld that is relevant to it.
		//the +2 here again is because we need the info of the rows exactly above and below the rows the worker processes.

		row := -1
		for cell := 0; cell < p.imagePartWidth; cell++ {
			actualRow := p.offset + row
			if actualRow >= 0 && actualRow < p.imageHeight {
				topRow[cell] = (*p.prevWorld)[actualRow][cell]
			} else if actualRow == -1 {
				topRow[cell] = (*p.prevWorld)[p.imageHeight-1][cell]
			} else {
				fmt.Println("error!")
			}
		}

		row = p.imagePartHeight

		for cell := 0; cell < p.imagePartWidth; cell++ {
			actualRow := p.offset + row
			if actualRow >= 0 && actualRow < p.imageHeight {
				bottomRow[cell] = (*p.prevWorld)[actualRow][cell]
			} else if actualRow == p.imageHeight {
				bottomRow[cell] = (*p.prevWorld)[0][cell]
			} else {

				fmt.Println("error!")
			}
		}

		//it then sends its board to the above function to do the calculations for one turn

		calculateNextState(p, oldGridPart, newGridPart, topRow, bottomRow, c.events)

		//even though the local grid is 2 rows bigger, the top and bottom row are ommited when writing back to nextWorld

		for row := 0; row < p.imagePartHeight; row++ {
			for cell := 0; cell < p.imagePartWidth; cell++ {
				(*p.nextWorld)[row+p.offset][cell] = newGridPart[row][cell]
			}
		}

		placeholderGridPart = oldGridPart
		oldGridPart = newGridPart
		newGridPart = placeholderGridPart

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
