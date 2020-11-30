package gol

import (
	"fmt"
	"os"
	"strconv"
)

func makeFile(imageHeight int, imageWidth int, world [][]byte, filename string) {
	_ = os.Mkdir("out", os.ModePerm)
	_ = os.Chdir("out")

	file, _ := os.Create(filename)
	//check(ioError)
	defer file.Close()

	_, _ = file.WriteString("P5\n")
	//_, _ = file.WriteString("# PGM file writer by pnmmodules (https://github.com/owainkenwayucl/pnmmodules).\n")
	_, _ = file.WriteString(strconv.Itoa(imageWidth))
	_, _ = file.WriteString(" ")
	_, _ = file.WriteString(strconv.Itoa(imageHeight))
	_, _ = file.WriteString("\n")
	_, _ = file.WriteString(strconv.Itoa(255))
	_, _ = file.WriteString("\n")

	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			_, _ = file.Write([]byte{world[y][x]})
			//check(ioError)
		}
	}

	//ioError = file.Sync()
	//check(ioError)

	fmt.Println("File", filename, "output done!")
	_ = os.Chdir("..")
}
