package config

import "github.com/osharko/sudoku/internal/pogo"

var (
	sudokuConfigSingleton pogo.Singleton
	sudokuConfig          *SudokuConfig = new(SudokuConfig)
)

type SudokuConfig struct {
	SquareSize      uint8     `yaml:"square_size"`
	Shapes          [][]uint8 `yaml:"shapes"`
	RequiredNumbers []uint8   `yaml:"required_numbers"`
}

func GetSudokuConfig() *SudokuConfig {
	sudokuConfigSingleton.Once(func() {
		InitConfiguration("config", sudokuConfig)
	})
	return sudokuConfig
}
