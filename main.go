package main

import (
	"chess2/engine"
	"fmt"
	"sort"
	"time"
)

const (
	Counts = 100000
	Size   = 12
)

func main() {

	var results []*engine.Board
	for i := 0; i < Counts-1; i++ {
		results = append(results, engine.InitBoard(Size))
	}
	var best *engine.Board
	var b int8
	score := make(map[int8]int)
	start := time.Now()
	for _, res := range results {
		res.SimRandom()
		if res.Step > b {
			best = res
			b = res.Step
		}
		score[res.Step]++
	}
	fmt.Println("Time elapsed: ", time.Since(start))
	fmt.Println("Time for 1 simulations: ", time.Since(start)/Counts)
	fmt.Println("Всего ферзей:", best.Step)
	fmt.Println("Порядок расстановки", best.Ranking)
	best.Print()
	keys := make([]int, 0, len(score))
	for k := range score {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Println("Ферзей:", k, "Результатов:", score[int8(k)], "Доля:", (float32(score[int8(k)])/Counts)*100, "%")
	}

}
