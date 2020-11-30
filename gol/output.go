package gol

import (
	"fmt"
	"io/ioutil"
	"log"
)

func writeFile(grid [][]uint8, width, height, turns int) {
	filename := "out/" + fmt.Sprintf("%vx%vx%v.pgm", width, height, turns)
	header := "P5\n" + fmt.Sprintf("%d %d\n%d\n", width, height, 255)
	content := []byte(header)
	for row := 0; row < height; row++ {
		for cell := 0; cell < width; cell++ {
			content = append(content, byte(grid[row][cell]))
		}
	}
	err := ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
