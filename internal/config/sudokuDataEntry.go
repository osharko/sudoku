package config

import "github.com/osharko/sudoku/internal/pogo"

var (
	sudokuDataEntrySingleton pogo.Singleton
	sudokuDataEntry          *SudokuDataEntry = new(SudokuDataEntry)
)

type SudokuDataEntry struct {
	Grid [][]uint8 `yaml:"grid"`
}

func GetSudokuDataEntry() *SudokuDataEntry {
	sudokuDataEntrySingleton.Once(func() {
		InitConfiguration("dataEntry", sudokuDataEntry)
	})
	return sudokuDataEntry
}
