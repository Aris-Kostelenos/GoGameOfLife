package main

import (
	"github.com/ChrisGora/semaphore"
)

type worker struct {
	prevWorld *[][]uint8
	nextWorld *[][]uint8
	startRow  int
	endRow    int
	width     int
	work      semaphore.Semaphore
	space     semaphore.Semaphore
}

// LIVE constant
const LIVE = byte(255)

// DEAD constant
const DEAD = byte(0)

// completes a single turn for its strip
func (w *worker) calculateNextState(turn int) {
	for row := w.startRow; row <= w.endRow; row++ {
		for col := 0; col < w.width; col++ {
			// initialise variables
			left := col - 1
			right := col + 1
			top := row - 1
			bottom := row + 1
			live := 0

			// wrap indices if necessary
			if row == 0 {
				top = len(*w.prevWorld) - 1
			} else if row == len(*w.prevWorld)-1 {
				bottom = 0
			}

			if col == 0 {
				left = w.width - 1
			} else if col == w.width-1 {
				right = 0
			}

			// total up the number of live neighbours
			live += int((*w.prevWorld)[top][col]) / 255
			live += int((*w.prevWorld)[bottom][col]) / 255
			live += int((*w.prevWorld)[row][left]) / 255
			live += int((*w.prevWorld)[row][right]) / 255
			live += int((*w.prevWorld)[top][left]) / 255
			live += int((*w.prevWorld)[top][right]) / 255
			live += int((*w.prevWorld)[bottom][left]) / 255
			live += int((*w.prevWorld)[bottom][right]) / 255

			// determine if the cell is worthy
			if (*w.prevWorld)[row][col] == LIVE {
				if live == 2 || live == 3 {
					(*w.nextWorld)[row][col] = LIVE
				} else {
					(*w.nextWorld)[row][col] = DEAD
				}
			} else {
				if live == 3 {
					(*w.nextWorld)[row][col] = LIVE
				} else {
					(*w.nextWorld)[row][col] = DEAD
				}
			}
		}
	}
}

// completes turns for its strip
func (w *worker) processStrip() {
	for turn := 0; ; turn++ {
		w.work.Wait()
		w.calculateNextState(turn)
		w.space.Post()
	}
}

// returns a slice of the alive cells in prevWorld
// func getAliveCells(prevWorld [][]uint8) []util.Cell {
// 	alive := make([]util.Cell, 0)
// 	for row := range prevWorld {
// 		for col := range prevWorld[row] {
// 			if prevWorld[row][col] == LIVE {
// 				alive = append(alive, util.Cell{X: col, Y: row})
// 			}
// 		}
// 	}
// 	return alive
// }
