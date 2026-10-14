package engine

import (
	"math"
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

// SimAnnealing ищет расстановку методом имитации отжига.
func (b *Board) SimAnnealing(
	rng *rand.Rand,
	iterations int,
	startTemp float64,
	endTemp float64,
) *Board {
	n := len(b.Cells)
	if n == 0 || iterations <= 0 || rng == nil {
		return b
	}
	if startTemp <= 0 || endTemp <= 0 {
		return b
	}

	// Для каждого столбца храним строку, в которой стоит кандидат.
	rows := make([]int, n)
	for col := range rows {
		rows[col] = rng.Intn(n)
	}

	bestRows := append([]int(nil), rows...)
	bestConflicts := countConflicts(rows)

	currentConflicts := bestConflicts

	for i := 0; i < iterations; i++ {
		// Плавное охлаждение от startTemp до endTemp.
		progress := float64(i) / float64(iterations-1)
		if iterations == 1 {
			progress = 1
		}
		temp := startTemp * math.Pow(endTemp/startTemp, progress)

		col := rng.Intn(n)
		oldRow := rows[col]

		newRow := rng.Intn(n - 1)
		if newRow >= oldRow {
			newRow++
		}

		rows[col] = newRow
		newConflicts := countConflicts(rows)
		delta := newConflicts - currentConflicts

		accept := delta <= 0
		if !accept {
			probability := math.Exp(-float64(delta) / temp)
			accept = rng.Float64() < probability
		}

		if accept {
			currentConflicts = newConflicts
		} else {
			rows[col] = oldRow
		}

		if currentConflicts < bestConflicts {
			bestConflicts = currentConflicts
			copy(bestRows, rows)
		}

		if bestConflicts == 0 {
			break
		}
	}

	// Записываем в Board только взаимно безопасные ферзи.
	// Если полного решения не найдено, получится допустимая
	// частичная расстановка.
	b.Step = 0
	b.Ranking = nil
	for y := range b.Cells {
		for x := range b.Cells[y] {
			b.Cells[y][x] = 0
		}
	}

	for col, row := range bestRows {
		safe := true

		for _, queen := range b.Ranking {
			qx, qy := int(queen.x), int(queen.y)
			if qy == row || abs(col-qx) == abs(row-qy) {
				safe = false
				break
			}
		}

		if safe {
			b.Step++
			b.Set(int8(col), int8(row), b.Step)
		}
	}

	// Отмечаем атакованные клетки после установки всех ферзей.
	for _, queen := range b.Ranking {
		b.DrawAttackedLine(int(queen.x), int(queen.y))
	}

	return b
}

func countConflicts(rows []int) int {
	conflicts := 0

	for col1 := 0; col1 < len(rows); col1++ {
		for col2 := col1 + 1; col2 < len(rows); col2++ {
			sameRow := rows[col1] == rows[col2]
			sameDiagonal := abs(col1-col2) == abs(rows[col1]-rows[col2])

			if sameRow || sameDiagonal {
				conflicts++
			}
		}
	}

	return conflicts
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
