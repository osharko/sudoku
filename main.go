package main

import (
	"github.com/osharko/sudoku/internal"
	"github.com/osharko/sudoku/internal/config"
)

func main() {
	grid := config.GetSudokuDataEntry().Grid
	sudoku := internal.SudokuFactory(grid, false)

	sudoku.Solve()
}
