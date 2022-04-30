package config

import "github.com/osharko/sudoku/internal/pogo"

var (
	sudokuDataEntrySingleton pogo.Singleton
	sudokuDataEntry          *SudokuDataEntry = new(SudokuDataEntry)
)

type SudokuDataEntry struct {
	Grid [][]uint8 `yaml:"grid"`
}

func (s *SudokuDataEntry) GetLinearSlice() (ret []uint8) {
	for _, row := range s.Grid {
		ret = append(ret, row...)
	}

	return ret
}

func GetSudokuDataEntry() *SudokuDataEntry {
	sudokuDataEntrySingleton.Once(func() {
		InitConfiguration("dataEntry", sudokuDataEntry)
	})
	return sudokuDataEntry
}
