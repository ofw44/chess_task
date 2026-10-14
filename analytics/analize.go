package analytics

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"sync/atomic"

	"chess2/engine"
)

type HeatMap struct {
	Cells [][]atomic.Uint64
}

func NewHeatMap(size int) *HeatMap {
	cells := make([][]atomic.Uint64, size)
	for i := range cells {
		cells[i] = make([]atomic.Uint64, size)
	}
	return &HeatMap{Cells: cells}
}

func (h *HeatMap) AddBoard(board *engine.Board) {
	for y, row := range board.Cells {
		for x, value := range row {
			if value > 0 {
				h.Cells[y][x].Add(1)
			}
		}
	}
}

func (h *HeatMap) AddFirstQueen(board *engine.Board) {
	for y, row := range board.Cells {
		for x, value := range row {
			if value == 1 {
				h.Cells[y][x].Add(1)
				return
			}
		}
	}
}

type Simulation struct {
	Name              string
	HeatMap           *HeatMap // все ферзи всех запусков
	FirstQueenSuccess *HeatMap // первый ферзь только успешных запусков
	SuccessfulHeatMap *HeatMap
	Unique            sync.Map

	BoardCount   atomic.Uint64
	SuccessCount atomic.Uint64

	outcomesMu sync.Mutex
	outcomes   map[int8]uint64
}

func NewSimulation(name string, size int) *Simulation {
	return &Simulation{
		Name:              name,
		HeatMap:           NewHeatMap(size),
		SuccessfulHeatMap: NewHeatMap(size),
		FirstQueenSuccess: NewHeatMap(size),
		outcomes:          make(map[int8]uint64),
	}
}

// Run выполняет count независимых запусков.
// makeBoard должен каждый раз создавать новую доску.
func (s *Simulation) Run(
	count int,
	workers int,
	makeBoard func() *engine.Board,
	simulate func(*engine.Board),
) {
	jobs := make(chan struct{}, workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				board := makeBoard()
				simulate(board)

				s.BoardCount.Add(1)

				s.outcomesMu.Lock()
				s.outcomes[board.Step]++
				s.outcomesMu.Unlock()

				if int(board.Step) == len(board.Cells) {
					s.SuccessCount.Add(1)
					s.SuccessfulHeatMap.AddBoard(board)
					s.FirstQueenSuccess.AddFirstQueen(board)
				}

				s.HeatMap.AddBoard(board)

				key := canonicalKey(board)
				s.Unique.LoadOrStore(key, struct{}{})

			}
		}()
	}

	for i := 0; i < count; i++ {
		jobs <- struct{}{}
	}
	close(jobs)
	wg.Wait()
}

func canonicalKey(board *engine.Board) string {
	n := len(board.Cells)
	var best string

	for transform := 0; transform < 8; transform++ {
		cells := make([]byte, n*n)

		for y, row := range board.Cells {
			for x, value := range row {
				if value <= 0 {
					continue
				}
				tx, ty := transformCoord(x, y, n, transform)
				cells[ty*n+tx] = 1
			}
		}
		key := string(cells)
		if transform == 0 || key < best {
			best = key
		}
	}
	return best
}

func transformCoord(x, y, n, transform int) (int, int) {
	switch transform {
	case 0: // исходная
		return x, y
	case 1: // поворот 90°
		return n - 1 - y, x
	case 2: // поворот 180°
		return n - 1 - x, n - 1 - y
	case 3: // поворот 270°
		return y, n - 1 - x
	case 4: // отражение по вертикали
		return n - 1 - x, y
	case 5:
		return n - 1 - y, n - 1 - x
	case 6: // отражение по горизонтали
		return x, n - 1 - y
	default:
		return y, x
	}
}

func (s *Simulation) Print() {
	size := len(s.HeatMap.Cells)

	fmt.Printf("Симуляция: %s\n", s.Name)
	fmt.Printf("Размер доски: %d\n", size)
	fmt.Printf("Всего симуляций: %d\n", s.BoardCount.Load())
	fmt.Printf("Успешных симуляций: %d\n", s.SuccessCount.Load())
	fmt.Printf("Уникальных расстановок с учётом симметрий: %d\n", s.uniqueCount())

	fmt.Println("\nHeatmap — число попаданий ферзя в клетку:")

	printHeatMap("Все ферзи во всех запусках:", s.HeatMap)
	printHeatMap("Все ферзи в успешных запусках:", s.SuccessfulHeatMap)
	printHeatMap("Первый ферзь в успешных запусках:", s.FirstQueenSuccess)

	s.outcomesMu.Lock()
	outcomes := make(map[int8]uint64, len(s.outcomes))
	for step, count := range s.outcomes {
		outcomes[step] = count
	}
	s.outcomesMu.Unlock()

	steps := make([]int, 0, len(outcomes))
	for step := range outcomes {
		steps = append(steps, int(step))
	}
	sort.Ints(steps)

	fmt.Println("\nРезультаты по количеству ферзей:")
	for _, step := range steps {
		count := outcomes[int8(step)]
		fmt.Printf("Ферзей: %d\tЗапусков: %d\n", step, count)
	}
}

func (s *Simulation) uniqueCount() uint64 {
	var count uint64
	s.Unique.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

func heatColor(t float64) (r, g, b int) {
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	stops := [][3]int{
		{40, 80, 220},
		{30, 200, 220},
		{255, 220, 40},
		{230, 40, 30},
	}

	position := t * float64(len(stops)-1)
	index := int(position)

	if index >= len(stops)-1 {
		last := stops[len(stops)-1]
		return last[0], last[1], last[2]
	}

	fraction := position - float64(index)
	left := stops[index]
	right := stops[index+1]

	r = int(float64(left[0]) + fraction*float64(right[0]-left[0]))
	g = int(float64(left[1]) + fraction*float64(right[1]-left[1]))
	b = int(float64(left[2]) + fraction*float64(right[2]-left[2]))

	return r, g, b
}

func printHeatMap(title string, h *HeatMap) {
	fmt.Println("\n" + title)

	maxCount := uint64(0)
	cellWidth := 1

	for _, row := range h.Cells {
		for _, cell := range row {
			value := cell.Load()
			if value > maxCount {
				maxCount = value
			}
			if width := len(fmt.Sprint(value)); width > cellWidth {
				cellWidth = width
			}
		}
	}

	useColor := os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"

	for _, row := range h.Cells {
		for x, cell := range row {
			if x > 0 {
				fmt.Print(" ")
			}

			value := cell.Load()
			text := fmt.Sprintf("%*d", cellWidth, value)

			if useColor && maxCount > 0 {
				r, g, b := heatColor(float64(value) / float64(maxCount))
				fmt.Printf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
			} else {
				fmt.Print(text)
			}
		}
		fmt.Println()
	}
}
