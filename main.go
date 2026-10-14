package main

import (
	"chess2/analytics"
	"chess2/engine"
	"math/rand"
	"runtime"
	"time"
)

const (
	Counts = 100000
	Size   = 12
)

func main() {
	count := 1_000_000
	workers := runtime.NumCPU()
	size := 10
	randomSimulation := analytics.NewSimulation("random", size)
	randomSimulation.Run(
		count,
		workers,
		func() *engine.Board {
			return engine.InitBoard(size)
		},
		func(board *engine.Board) {
			board.SimRandom()
		},
	)
	randomSimulation.Print()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	board := engine.InitBoard(size)
	board.SimAnnealing(rng, 10_000, 10, 0.01)
	board.Print()
}
