package main

import (
	"testing"

	"uk.ac.bris.cs/gameoflife/gol"
)

func BenchmarkGoLBigOrGoLHome(b *testing.B) {

	params := gol.Params{
		Turns:       10000,
		ImageWidth:  5120,
		ImageHeight: 5120,
		Threads:     16,
	}

	//the function is run several times. N is increased automatically by the benchamrk runner
	//until the stability of the benchmark is confirmed
	for n := 0; n < b.N; n++ {
		eventsPlaceholder := make(chan gol.Event, 10)
		gol.Run(params, eventsPlaceholder, nil)
		for range eventsPlaceholder {

		}
	}

}

/*
func benchmarkFib(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
			Fib(i)
	}
}

func BenchmarkFib1(b *testing.B)  { benchmarkFib(1, b) }
func BenchmarkFib2(b *testing.B)  { benchmarkFib(2, b) }
func BenchmarkFib3(b *testing.B)  { benchmarkFib(3, b) }
func BenchmarkFib10(b *testing.B) { benchmarkFib(10, b) }
func BenchmarkFib20(b *testing.B) { benchmarkFib(20, b) }
func BenchmarkFib40(b *testing.B) { benchmarkFib(40, b) }
*/
