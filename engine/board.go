package engine

import (
	"fmt"
	"strconv"
)

type Board struct {
	Cells   [][]int8
	Step    int8
	Ranking []Coordinate
}

func InitBoard(size int) *Board {
	b := &Board{}
	b.Cells = make([][]int8, size)
	for i := range b.Cells {
		b.Cells[i] = make([]int8, size)
	}
	return b
}

func (b *Board) Print() {
	for _, row := range b.Cells {
		fmt.Println("")
		fmt.Print("|")
		for _, col := range row {
			if col == 0 {
				fmt.Print(" |")
			} else if col == -1 {
				fmt.Print("x|")
			} else {
				if col > 9 {
					fmt.Print(strconv.Itoa(int(col)) + "|")
				} else {
					fmt.Print(strconv.Itoa(int(col)) + "|")
				}
			}
		}
	}
	fmt.Println()
}

func (b *Board) Set(x, y, value int8) {
	b.Cells[y][x] = value
	if value > 0 {
		b.Ranking = append(b.Ranking, Coordinate{x, y})
	}
}

func (b *Board) DrawAttackedLine(x, y int) {
	for i := 0; i < len(b.Cells); i++ {
		for j := 0; j < len(b.Cells); j++ {
			if (i == y || j == x) && !(i == y && j == x) {
				b.Cells[i][j] = -1
			}
		}
	}
	i := y + 1
	j := x + 1
	for i < len(b.Cells) && j < len(b.Cells) {
		b.Cells[i][j] = -1
		j++
		i++
	}
	i = y - 1
	j = x + 1
	for i >= 0 && j < len(b.Cells) {
		b.Cells[i][j] = -1
		j++
		i--
	}
	i = y - 1
	j = x - 1
	for i >= 0 && j >= 0 {
		b.Cells[i][j] = -1
		j--
		i--
	}
	i = y + 1
	j = x - 1
	for i < len(b.Cells) && j >= 0 {
		b.Cells[i][j] = -1
		j--
		i++
	}
}
