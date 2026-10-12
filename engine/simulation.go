package engine

import (
	"math/rand"
)

type Coordinate struct {
	x int8
	y int8
}

func (b *Board) SimRandom() *Board {
	for {
		if b.Step == int8(len(b.Cells)) {
			break
		}
		empty := b.EmptyCells()
		if len(empty) == 0 {
			break
		}
		choose := empty[rand.Intn(len(empty))]
		b.Step++
		b.Set(choose.x, choose.y, b.Step)
		b.DrawAttackedLine(int(choose.x), int(choose.y))
	}
	return b
}

func (b *Board) EmptyCellsCount() int8 {
	var count int8
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell == 0 {
				count++
			}
		}
	}
	return count
}

func (b *Board) EmptyCells() []Coordinate {
	var Crd []Coordinate
	for y, row := range b.Cells {
		for x, cell := range row {
			if cell == 0 {
				Crd = append(Crd, Coordinate{int8(x), int8(y)})
			}
		}
	}
	return Crd
}
