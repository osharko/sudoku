package main

import (
	"github.com/osharko/sudoku/internal"
)

func main() {
	sudoku := internal.SudokuFactory()

	/* fmt.Println(sudoku.GetRowElements())
	fmt.Println(sudoku.GetColElements())
	fmt.Println(sudoku.GetShapeElements()) */

	sudoku.FindMissingValue()
}
