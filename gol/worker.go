package gol

import (
	"uk.ac.bris.cs/gameoflife/util"
)

func calculateNextState(p workerParams, world [][]uint8, borders borders) [][]uint8 {
	oldWorld := make([][]uint8, p.imagePartWidth+2)
	for z := 0; z < p.imagePartWidth+2; z++ {
		oldWorld[z] = make([]uint8, p.imagePartHeight)
	}
	for i := 0; i < p.imagePartHeight; i++ {
		for z := 1; z <= p.imagePartWidth; z++ {
			oldWorld[z][i] = world[i][z-1]
		}
		oldWorld[1][i] = borders.top[i]
		oldWorld[p.imagePartWidth+1][i] = borders.bottom[i]
	}

	newWorld := make([][]uint8, p.imagePartWidth)
	for z := 0; z < p.imagePartWidth; z++ {
		newWorld[z] = make([]uint8, p.imagePartHeight)
	}

	for row := range oldWorld {
		for cell := range oldWorld[row] {
			var a int
			var b int
			var c int
			var d int

			if row == 0 {
				a = p.imagePartWidth + 1
			} else {
				a = row - 1
			}
			if cell == 0 {
				c = p.imagePartHeight - 1
			} else {
				c = cell - 1
			}
			if row == p.imagePartWidth+1 {
				b = 0
			} else {
				b = row + 1
			}
			if cell == p.imagePartHeight-1 {
				d = 0
			} else {
				d = cell + 1
			}

			x := (oldWorld[a][c] + oldWorld[a][cell] + oldWorld[a][d] + oldWorld[row][c] + oldWorld[row][d] + oldWorld[b][c] + oldWorld[b][cell] + oldWorld[b][d]) / 255
			if row != 0 && row != p.imagePartWidth+1 {
				if oldWorld[row][cell] == 0 {
					if x == 3 {
						newWorld[row-1][cell] = 255
						p.events <- CellFlipped{255, util.Cell{X: row, Y: cell}}
					} else {
						newWorld[row-1][cell] = 0
					}
				}

				if oldWorld[row][cell] == 255 {
					if x == 3 || x == 2 {
						newWorld[row-1][cell] = 255
					} else {
						newWorld[row-1][cell] = 0
						p.events <- CellFlipped{0, util.Cell{X: row, Y: cell}}
					}
				}

			}

		}
	}
	return newWorld
}

/*
func calculateNextStatePart(p workerParams, world [][]uint8, b borders) [][]uint8 {

	//this shit receives the borders from the neighbouring goroutines using borders in order to calculet if the shit cell are alive or fukcing ded

	NewWorld := make([][]uint8, p.imagePartWidth)

	for z := 0; z < p.imagePartWidth; z++ {
		NewWorld[z] = make([]uint8, p.imagePartHeight)
	}
	for row := range world {
		for cell := range world[row] {
			var left int
			var right int
			var x uint8

			if cell == 0 {
				left = p.imagePartHeight - 1
			} else {
				left = cell - 1
			}
			if cell == p.imagePartHeight-1 {
				right = 0
			} else {
				right = cell + 1
			}

			// if the cell it calculates is on top or bottom instead of using world it uses the fucking border arrays top and bottom ok

			if row == 0 {
				x = ((b.top[left] + b.top[cell] + b.top[right] + world[row][left] + world[row][right] + world[row+1][left] + world[row+1][cell] + world[row+1][right]) / 255)

			} else if row == p.imagePartWidth-1 {
				x = ((world[row-1][left] + world[row-1][cell] + world[row-1][right] + world[row][left] + world[row][right] + b.bottom[left] + b.bottom[cell] + b.bottom[right]) / 255)

			} else {
				x = ((world[row-1][left] + world[row-1][cell] + world[row-1][right] + world[row][left] + world[row][right] + world[row+1][left] + world[row+1][cell] + world[row+1][right]) / 255)

			}

			//x := ((world[above][left] / 255) + (world[above][cell] / 255) + (world[above][right] / 255) + (world[row][left] / 255) + (world[row][right] / 255) + (world[below][left] / 255) + (world[below][cell] / 255) + (world[below][right] / 255))

			if world[row][cell] == 0 {
				if x == 3 {
					NewWorld[row][cell] = 255
					p.events <- CellFlipped{255, util.Cell{X: row, Y: cell}}
				} else {
					NewWorld[row][cell] = 0
				}
			}

			if world[row][cell] == 255 {
				if x == 3 || x == 2 {
					NewWorld[row][cell] = 255
				} else {
					NewWorld[row][cell] = 0
					//p.events <- CellFlipped{0, util.Cell{X: row, Y: cell}}
				}
			}
		}
	}
	return NewWorld
}

*/

func calculateAliveCells(p workerParams, world [][]uint8) []util.Cell {
	a := make([]util.Cell, 0)
	k := 0
	for row := range world {
		for cell := range world[row] {
			if world[cell][row] == 255 {
				a = append(a, util.Cell{X: cell, Y: row})
				k++
			}
		}
	}
	return a
}

func workerGoroutine(p workerParams, immPrevWorld func(row, cell int) uint8, nextWorld [][]uint8) {

	gridPart := make([][]uint8, p.imagePartHeigth+2)
	for row := 0; row < p.imagePartHeight+2; row++ {
		gridPart[row] = make([]uint8, p.imagePartWidth)
		for cell := 0; cell < p.imagePartWidth; cell++ {
			gridPart[row][cell] = immPrevWorld(row+p.imagePartHeightStartpoint-1, cell)
		}
	}

	for turns := 0; turns < p.turns; turns++ {
		gridPart = calculateNextState(p, gridPart, temp)

		if p.id != 0 && p.id != p.threads-1 {

		} else if p.id == 0 && p.id == p.threads-1 {

		} else if p.id == 0 {

		} else {

		}

	}

	for row := 0; row < p.imagePartHeight; row++ {
		for cell := 0; cell < p.imagePartWidth; cell++ {
			world[cell][row+p.imagePartWidthStartpoint] = gridPart[row][cell]
		}
	}
	statusChannel[p.id] <- false

}
