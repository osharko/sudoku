package main

import (
	"github.com/osharko/sudoku/internal"
)

func main() {
	sudoku := internal.SudokuFactory()

	sudoku.FindMissingValue()
}
