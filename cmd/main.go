package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIDTH  = 800
	HEIGHT = 800

	PADDING = 10

	FPS = 10

	CELLSIZE            = 5
	LiveCellProbability = 0.1

	rowCount = (HEIGHT - (PADDING * 2)) / CELLSIZE
	colCount = (WIDTH - (PADDING * 2)) / CELLSIZE
)

var (
	NeighboursMatrix = [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1} /*{0, 0},*/, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	Glider = []struct{ Row, Col int }{
		{0, 1},
		{1, 2},
		{2, 0}, {2, 1}, {2, 2},
	}

	SmallExploder = []struct{ Row, Col int }{
		{0, 1},
		{1, 0}, {1, 1}, {1, 2},
		{2, 0}, {2, 2},
		{3, 1},
	}

	Toad = []struct{ Row, Col int }{
		{1, 1}, {1, 2}, {1, 3},
		{2, 0}, {2, 1}, {2, 2},
	}

	Beacon = []struct{ Row, Col int }{
		{0, 0}, {0, 1},
		{1, 0},
		{2, 3},
		{3, 2}, {3, 3},
	}

	Pulsar = []struct{ Row, Col int }{
		{2, 4}, {2, 5}, {2, 6}, {2, 10}, {2, 11}, {2, 12},
		{4, 2}, {5, 2}, {6, 2}, {4, 7}, {5, 7}, {6, 7},
		{4, 9}, {5, 9}, {6, 9}, {4, 14}, {5, 14}, {6, 14},
		{7, 4}, {7, 5}, {7, 6}, {7, 10}, {7, 11}, {7, 12},
		{9, 4}, {9, 5}, {9, 6}, {9, 10}, {9, 11}, {9, 12},
		{10, 2}, {11, 2}, {12, 2}, {10, 7}, {11, 7}, {12, 7},
		{10, 9}, {11, 9}, {12, 9}, {10, 14}, {11, 14}, {12, 14},
		{14, 4}, {14, 5}, {14, 6}, {14, 10}, {14, 11}, {14, 12},
	}

	GosperGliderGun = []struct{ Row, Col int }{
		{5, 1}, {5, 2}, {6, 1}, {6, 2},
		{3, 13}, {3, 14}, {4, 12}, {4, 16}, {5, 11}, {5, 17}, {6, 11}, {6, 15},
		{6, 17}, {6, 18}, {7, 11}, {7, 17}, {8, 12}, {8, 16}, {9, 13}, {9, 14},
		{1, 25}, {2, 23}, {2, 25}, {3, 21}, {3, 22}, {4, 21}, {4, 22}, {5, 21},
		{5, 22}, {6, 23}, {6, 25}, {7, 25},
		{3, 35}, {3, 36}, {4, 35}, {4, 36},
	}
)

type Cell struct {
	Row             int
	Col             int
	IsAlive         bool
	IsAliveInFuture bool
}

type Configration struct {
	AliveCells []struct{ Row, Col int }
	UseConfig  bool
}

func NewCell(r int, c int, a bool) *Cell {
	return &Cell{
		Row:             r,
		Col:             c,
		IsAlive:         a,
		IsAliveInFuture: a,
	}
}

func NewConfiguration(c *[]struct{ Row, Col int }) *Configration {
	return &Configration{
		AliveCells: *c,
		UseConfig:  true,
	}
}

func InitRandomGrid() [][]Cell {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)

	grid := make([][]Cell, rowCount)

	for r := range rowCount {
		grid[r] = make([]Cell, colCount)
		for c := range colCount {
			a := generator.Float64() < LiveCellProbability
			grid[r][c] = *NewCell(r, c, a)

		}
	}

	return grid
}

func InitGridFromConfig(config *Configration) [][]Cell {
	grid := make([][]Cell, rowCount)

	for r := range rowCount {
		grid[r] = make([]Cell, colCount)
		for c := range colCount {
			grid[r][c] = *NewCell(r, c, false)
		}
	}

	if config != nil && config.UseConfig {
		for _, cell := range config.AliveCells {
			grid[cell.Row][cell.Col] = *NewCell(cell.Row, cell.Col, true)
		}
	}

	return grid
}

func DrawCellGrid(g *[][]Cell) {
	gridValue := *g

	for r := range gridValue {
		for c := range gridValue[r] {
			cell := gridValue[r][c]

			if cell.IsAlive {
				cellXPos := PADDING + cell.Col*CELLSIZE
				cellYPos := PADDING + cell.Row*CELLSIZE

				rl.DrawRectangle(int32(cellXPos), int32(cellYPos), int32(CELLSIZE), int32(CELLSIZE), rl.Green)
			}
		}
	}
}

func UpdateCellGrid(g *[][]Cell) {
	for r := range rowCount {
		for c := range colCount {
			liveNeighbourCount := 0
			for _, offset := range NeighboursMatrix {
				rowIdx := r + offset[0]
				colIdx := c + offset[1]

				if rowIdx >= 0 && rowIdx < rowCount && colIdx >= 0 && colIdx < colCount && (*g)[rowIdx][colIdx].IsAlive {
					liveNeighbourCount++
				}
			}

			cell := &(*g)[r][c]

			if cell.IsAlive {
				if liveNeighbourCount < 2 || liveNeighbourCount > 3 {
					cell.IsAliveInFuture = false
				} else {
					cell.IsAliveInFuture = true
				}
			} else {
				if liveNeighbourCount == 3 {
					cell.IsAliveInFuture = true
				} else {
					cell.IsAliveInFuture = false
				}
			}
		}
	}

	for r := range rowCount {
		for c := range colCount {
			(*g)[r][c].IsAlive = (*g)[r][c].IsAliveInFuture
		}
	}
}

func main() {
	if colCount <= 0 || rowCount <= 0 {
		panic("Calculated grid dimensions are not positive. Adjust PADDING or CELLSIZE.")
	}

	rl.InitWindow(WIDTH, HEIGHT, "Conway's Game of Life")
	defer rl.CloseWindow()

	rl.SetTargetFPS(FPS)

	// grid := InitRandomGrid()

	c := NewConfiguration(&Pulsar)
	grid := InitGridFromConfig(c)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		DrawCellGrid(&grid)
		UpdateCellGrid(&grid)

		rl.EndDrawing()
	}
}
