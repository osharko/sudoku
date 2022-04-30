package internal

import (
	"fmt"
	"math"

	"github.com/osharko/sudoku/internal/config"
	"github.com/osharko/sudoku/internal/pogo"
)

type sudoku struct {
	//Stats
	currentCol        uint8 //Holds the current column
	currentRow        uint8 //Holds the current row
	iteration         uint8 //How many iteration have been done
	startMissingValue uint8 //Holds the number of missing values at the beginning of the sudoku
	//Data
	grid [][]uint8 // Slice which hold all sudoku table.
	//Configuration
	size            uint8     // Size of a side of the grid. Ex: 9 for a 9x9 grid.
	shapes          [][]uint8 // Represent all the shape of the grid. Each shape is made by the element's position into the grid.
	requiredNumbers []uint8   //All the value that must be present in shape/column/row.
}

// Since sudoku is a private struct, the only way to create a new sudoku is to use the SudokuFactory function.
// This function is used to create a new sudoku, with the related configuration. It's a workaround due to golang's lack of constructor.
func SudokuFactory(grid [][]uint8) (s sudoku) {
	configuration := config.GetSudokuConfig()

	s = sudoku{
		currentCol:      0,
		currentRow:      0,
		iteration:       0,
		grid:            grid,
		size:            configuration.SquareSize,
		shapes:          configuration.Shapes,
		requiredNumbers: configuration.RequiredNumbers,
	}

	s.startMissingValue = s.countMissingValues()

	return
}

func (s *sudoku) Solve() {
	for s.iteration = 1; s.iteration <= s.size*s.size && !s.isComplete(); s.iteration++ {
		s.PrintGrid()
		s.findMissingValue()
	}

	s.PrintGrid()
}

// Return all the element in the current row.
func (s *sudoku) getRowElements() (ret []uint8) {
	for i := uint8(0); i < s.size; i++ {
		ret = append(ret, s.grid[s.currentRow][i])
	}

	return
}

// Return all the element in the current column.
func (s *sudoku) getColElements() (ret []uint8) {
	for i := uint8(0); i < s.size; i++ {
		ret = append(ret, s.grid[i][s.currentCol])
	}

	return
}

// Accoring to the current row and column, returns all the element of the current shape.
func (s *sudoku) getShapeElements() []uint8 {
	sizeRoot := math.Sqrt(float64(s.size))

	row := uint8((math.Floor(float64(s.currentRow)/sizeRoot) * sizeRoot))
	col := uint8(math.Floor(float64(s.currentCol) / sizeRoot))

	return s.getShape(row + col)
}

// Returns all element in a given shape.
func (s *sudoku) getShape(shapePos uint8) (ret []uint8) {
	for _, element := range s.shapes[shapePos] {
		row := (element) / s.size
		col := (element) % s.size
		ret = append(ret, s.grid[row][col])
	}

	return
}

// Find all the missing number in the current row, column and shape,
// Thene check if those 3 has 1 common missing number,
// If so, fill the current cell with that number.
func (s *sudoku) findValue() uint8 {
	removeFromArray := func(origin []uint8, comparison []uint8) []uint8 {
		rem := make([]uint8, 0)

		for _, v := range origin {
			for _, value := range comparison {
				if v == value {
					rem = append(rem, v)
				}
			}
		}

		return pogo.FilterArray(origin, func(o uint8) bool {
			return !pogo.ContainsArray(rem, o)
		})
	}

	values := make([]uint8, len(s.requiredNumbers))
	copy(values, s.requiredNumbers)
	r, c, h := s.getRowElements(), s.getColElements(), s.getShapeElements()
	values = removeFromArray(values, r)
	values = removeFromArray(values, c)
	values = removeFromArray(values, h)

	if len(values) == 1 {
		return values[0]
	}

	return 0
}

// If there's no 0 value into the grid, then the sudoku is complete.
func (s *sudoku) isComplete() bool {
	return pogo.EveryInArray(s.grid, func(row []uint8) bool {
		return !pogo.ContainsArray(row, 0)
	})
}

func (s *sudoku) findMissingValue() {
	for i, row := range s.grid {
		for j, value := range row {
			if value != 0 {
				continue
			}

			s.currentRow = uint8(i)
			s.currentCol = uint8(j)

			if v := s.findValue(); v != 0 {
				s.grid[i][j] = v
				return
			}
		}
	}
}

// Return the number of 0 into the grid.
func (s *sudoku) countMissingValues() uint8 {
	count := uint8(0)
	for _, row := range s.grid {
		for _, value := range row {
			if value == 0 {
				count++
			}
		}
	}

	return count
}

func (s *sudoku) PrintGrid() {
	fmt.Printf("\nCurrent Iteration: %d\tMissing Values: %d\tStarting Missing Values: %d\tCompleted: %t\n\n", s.iteration, s.countMissingValues(), s.startMissingValue, s.isComplete())
	fmt.Printf("\n\t\t\t")

	root := int(math.Sqrt(float64(s.size)))

	for i, row := range s.grid {
		for j, value := range row {
			//Print the current value
			fmt.Printf("%d ", value)
			//Reaching the end of the shape, print a tab.
			if (j+1)%root == 0 {
				fmt.Printf("\t")
			}
		}
		// Print an additional line, at the end of the shape.
		if (i+1)%root == 0 {
			fmt.Println()
		}
		// Print a new line every time the end of a row have been reached.
		fmt.Printf("\n\t\t\t")
	}
	fmt.Printf("\n----------------------------------------------------------------------------\n\n")
}
