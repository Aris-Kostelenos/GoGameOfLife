package main

import (
	"testing"

	"uk.ac.bris.cs/gameoflife/gol"
)

func benchmarkGoL(turns int, size int, threads int, b *testing.B) {
	params := gol.Params{
		Turns:       turns,
		ImageWidth:  size,
		ImageHeight: size,
	}

	params.Threads = threads
	for n := 0; n < b.N; n++ {
		eventsPlaceholder := make(chan gol.Event, 1000)
		gol.Run(params, eventsPlaceholder, nil)
		for range eventsPlaceholder {

		}
	}
	//the function is run several times. N is increased automatically by the benchamrk runner
	//until the stability of the benchmark is confirmed
}

//go test -run=XXXXXXXX -bench=. -benchtime=30s -timeout 70m

func BenchmarkGoLMedium1worker1000turn(b *testing.B)  { benchmarkGoL(1000, 512, 1, b) }
func BenchmarkGoLMedium2worker1000turn(b *testing.B)  { benchmarkGoL(1000, 512, 2, b) }
func BenchmarkGoLMedium3worker1000turn(b *testing.B)  { benchmarkGoL(1000, 512, 3, b) }
func BenchmarkGoLMedium4worker1000turn(b *testing.B)  { benchmarkGoL(1000, 512, 4, b) }
func BenchmarkGoLMedium5worker1000turn(b *testing.B)  { benchmarkGoL(1000, 512, 5, b) }
func BenchmarkGoLMedium6worker1000turn(b *testing.B)  { benchmarkGoL(1000, 512, 6, b) }
func BenchmarkGoLMedium7worker1000turn(b *testing.B)  { benchmarkGoL(1000, 512, 7, b) }
func BenchmarkGoLMedium8worker1000turn(b *testing.B)  { benchmarkGoL(1000, 512, 8, b) }
func BenchmarkGoLMedium9worker1000turn(b *testing.B)  { benchmarkGoL(1000, 512, 9, b) }
func BenchmarkGoLMedium10worker1000turn(b *testing.B) { benchmarkGoL(1000, 512, 10, b) }
func BenchmarkGoLMedium11worker1000turn(b *testing.B) { benchmarkGoL(1000, 512, 11, b) }
func BenchmarkGoLMedium12worker1000turn(b *testing.B) { benchmarkGoL(1000, 512, 12, b) }
func BenchmarkGoLMedium13worker1000turn(b *testing.B) { benchmarkGoL(1000, 512, 13, b) }
func BenchmarkGoLMedium14worker1000turn(b *testing.B) { benchmarkGoL(1000, 512, 14, b) }
func BenchmarkGoLMedium15worker1000turn(b *testing.B) { benchmarkGoL(1000, 512, 15, b) }
func BenchmarkGoLMedium16worker1000turn(b *testing.B) { benchmarkGoL(1000, 512, 16, b) }

func BenchmarkGoLImageSize16(b *testing.B)   { benchmarkGoL(0, 16, 1, b) }
func BenchmarkGoLImageSize64(b *testing.B)   { benchmarkGoL(0, 64, 1, b) }
func BenchmarkGoLImageSize128(b *testing.B)  { benchmarkGoL(0, 128, 1, b) }
func BenchmarkGoLImageSize256(b *testing.B)  { benchmarkGoL(0, 256, 1, b) }
func BenchmarkGoLImageSize512(b *testing.B)  { benchmarkGoL(0, 512, 1, b) }
func BenchmarkGoLImageSize5120(b *testing.B) { benchmarkGoL(0, 5120, 1, b) }

func BenchmarkGoLMedium1workers0turns(b *testing.B)    { benchmarkGoL(0, 512, 1, b) }
func BenchmarkGoLMedium1workers100turns(b *testing.B)  { benchmarkGoL(100, 512, 1, b) }
func BenchmarkGoLMedium1workers200turns(b *testing.B)  { benchmarkGoL(200, 512, 1, b) }
func BenchmarkGoLMedium1workers300turns(b *testing.B)  { benchmarkGoL(300, 512, 1, b) }
func BenchmarkGoLMedium1workers400turns(b *testing.B)  { benchmarkGoL(400, 512, 1, b) }
func BenchmarkGoLMedium1workers500turns(b *testing.B)  { benchmarkGoL(500, 512, 1, b) }
func BenchmarkGoLMedium1workers600turns(b *testing.B)  { benchmarkGoL(600, 512, 1, b) }
func BenchmarkGoLMedium1workers700turns(b *testing.B)  { benchmarkGoL(700, 512, 1, b) }
func BenchmarkGoLMedium1workers800turns(b *testing.B)  { benchmarkGoL(800, 512, 1, b) }
func BenchmarkGoLMedium1workers900turns(b *testing.B)  { benchmarkGoL(900, 512, 1, b) }
func BenchmarkGoLMedium1workers1000turns(b *testing.B) { benchmarkGoL(1000, 512, 1, b) }

func BenchmarkGoLMedium16workers0turns(b *testing.B)    { benchmarkGoL(0, 512, 16, b) }
func BenchmarkGoLMedium16workers100turns(b *testing.B)  { benchmarkGoL(100, 512, 16, b) }
func BenchmarkGoLMedium16workers200turns(b *testing.B)  { benchmarkGoL(200, 512, 16, b) }
func BenchmarkGoLMedium16workers300turns(b *testing.B)  { benchmarkGoL(300, 512, 16, b) }
func BenchmarkGoLMedium16workers400turns(b *testing.B)  { benchmarkGoL(400, 512, 16, b) }
func BenchmarkGoLMedium16workers500turns(b *testing.B)  { benchmarkGoL(500, 512, 16, b) }
func BenchmarkGoLMedium16workers600turns(b *testing.B)  { benchmarkGoL(600, 512, 16, b) }
func BenchmarkGoLMedium16workers700turns(b *testing.B)  { benchmarkGoL(700, 512, 16, b) }
func BenchmarkGoLMedium16workers800turns(b *testing.B)  { benchmarkGoL(800, 512, 16, b) }
func BenchmarkGoLMedium16workers900turns(b *testing.B)  { benchmarkGoL(900, 512, 16, b) }
func BenchmarkGoLMedium16workers1000turns(b *testing.B) { benchmarkGoL(1000, 512, 16, b) }

/*
func BenchmarkGoLMedium1worker0turn(b *testing.B)  { benchmarkGoL(0, 512, 1, b) }
func BenchmarkGoLMedium2worker0turn(b *testing.B)  { benchmarkGoL(0, 512, 2, b) }
func BenchmarkGoLMedium3worker0turn(b *testing.B)  { benchmarkGoL(0, 512, 3, b) }
func BenchmarkGoLMedium4worker0turn(b *testing.B)  { benchmarkGoL(0, 512, 4, b) }
func BenchmarkGoLMedium5worker0turn(b *testing.B)  { benchmarkGoL(0, 512, 5, b) }
func BenchmarkGoLMedium6worker0turn(b *testing.B)  { benchmarkGoL(0, 512, 6, b) }
func BenchmarkGoLMedium7worker0turn(b *testing.B)  { benchmarkGoL(0, 512, 7, b) }
func BenchmarkGoLMedium8worker0turn(b *testing.B)  { benchmarkGoL(0, 512, 8, b) }
func BenchmarkGoLMedium9worker0turn(b *testing.B)  { benchmarkGoL(0, 512, 9, b) }
func BenchmarkGoLMedium10worker0turn(b *testing.B) { benchmarkGoL(0, 512, 10, b) }
func BenchmarkGoLMedium11worker0turn(b *testing.B) { benchmarkGoL(0, 512, 11, b) }
func BenchmarkGoLMedium12worker0turn(b *testing.B) { benchmarkGoL(0, 512, 12, b) }
func BenchmarkGoLMedium13worker0turn(b *testing.B) { benchmarkGoL(0, 512, 13, b) }
func BenchmarkGoLMedium14worker0turn(b *testing.B) { benchmarkGoL(0, 512, 14, b) }
func BenchmarkGoLMedium15worker0turn(b *testing.B) { benchmarkGoL(0, 512, 15, b) }
func BenchmarkGoLMedium16worker0turn(b *testing.B) { benchmarkGoL(0, 512, 16, b) }

func BenchmarkGoLSmall1worker1000turn(b *testing.B)  { benchmarkGoL(1000, 64, 64, 1, b) }
func BenchmarkGoLSmall2worker1000turn(b *testing.B)  { benchmarkGoL(1000, 64, 64, 2, b) }
func BenchmarkGoLSmall3worker1000turn(b *testing.B)  { benchmarkGoL(1000, 64, 64, 3, b) }
func BenchmarkGoLSmall4worker1000turn(b *testing.B)  { benchmarkGoL(1000, 64, 64, 4, b) }
func BenchmarkGoLSmall5worker1000turn(b *testing.B)  { benchmarkGoL(1000, 64, 64, 5, b) }
func BenchmarkGoLSmall6worker1000turn(b *testing.B)  { benchmarkGoL(1000, 64, 64, 6, b) }
func BenchmarkGoLSmall7worker1000turn(b *testing.B)  { benchmarkGoL(1000, 64, 64, 7, b) }
func BenchmarkGoLSmall8worker1000turn(b *testing.B)  { benchmarkGoL(1000, 64, 64, 8, b) }
func BenchmarkGoLSmall9worker1000turn(b *testing.B)  { benchmarkGoL(1000, 64, 64, 9, b) }
func BenchmarkGoLSmall10worker1000turn(b *testing.B) { benchmarkGoL(1000, 64, 64, 10, b) }
func BenchmarkGoLSmall11worker1000turn(b *testing.B) { benchmarkGoL(1000, 64, 64, 11, b) }
func BenchmarkGoLSmall12worker1000turn(b *testing.B) { benchmarkGoL(1000, 64, 64, 12, b) }
func BenchmarkGoLSmall13worker1000turn(b *testing.B) { benchmarkGoL(1000, 64, 64, 13, b) }
func BenchmarkGoLSmall14worker1000turn(b *testing.B) { benchmarkGoL(1000, 64, 64, 14, b) }
func BenchmarkGoLSmall15worker1000turn(b *testing.B) { benchmarkGoL(1000, 64, 64, 15, b) }
func BenchmarkGoLSmall16worker1000turn(b *testing.B) { benchmarkGoL(1000, 64, 64, 16, b) }

func BenchmarkGoLSmall1worker0turn(b *testing.B)  { benchmarkGoL(0, 64, 64, 1, b) }
func BenchmarkGoLSmall2worker0turn(b *testing.B)  { benchmarkGoL(0, 64, 64, 2, b) }
func BenchmarkGoLSmall3worker0turn(b *testing.B)  { benchmarkGoL(0, 64, 64, 3, b) }
func BenchmarkGoLSmall4worker0turn(b *testing.B)  { benchmarkGoL(0, 64, 64, 4, b) }
func BenchmarkGoLSmall5worker0turn(b *testing.B)  { benchmarkGoL(0, 64, 64, 5, b) }
func BenchmarkGoLSmall6worker0turn(b *testing.B)  { benchmarkGoL(0, 64, 64, 6, b) }
func BenchmarkGoLSmall7worker0turn(b *testing.B)  { benchmarkGoL(0, 64, 64, 7, b) }
func BenchmarkGoLSmall8worker0turn(b *testing.B)  { benchmarkGoL(0, 64, 64, 8, b) }
func BenchmarkGoLSmall9worker0turn(b *testing.B)  { benchmarkGoL(0, 64, 64, 9, b) }
func BenchmarkGoLSmall10worker0turn(b *testing.B) { benchmarkGoL(0, 64, 64, 10, b) }
func BenchmarkGoLSmall11worker0turn(b *testing.B) { benchmarkGoL(0, 64, 64, 11, b) }
func BenchmarkGoLSmall12worker0turn(b *testing.B) { benchmarkGoL(0, 64, 64, 12, b) }
func BenchmarkGoLSmall13worker0turn(b *testing.B) { benchmarkGoL(0, 64, 64, 13, b) }
func BenchmarkGoLSmall14worker0turn(b *testing.B) { benchmarkGoL(0, 64, 64, 14, b) }
func BenchmarkGoLSmall15worker0turn(b *testing.B) { benchmarkGoL(0, 64, 64, 15, b) }
func BenchmarkGoLSmall16worker0turn(b *testing.B) { benchmarkGoL(0, 64, 64, 16, b) }
*/
