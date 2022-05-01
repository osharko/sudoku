package internal

import (
	"fmt"
	"math"

	"github.com/osharko/sudoku/internal/config"
	"github.com/osharko/sudoku/internal/pogo"
)

var (
	baseColor           = "\033[32m"
	highlightColor      = "\033[33m"
	superHighlightColor = "\033[35m"

	newFound *sudokuElement
)

type sudokuElement struct {
	value        uint8
	isStartValue bool
}

func (e *sudokuElement) Fromuint8(matr [][]uint8) [][]sudokuElement {
	ret := make([][]sudokuElement, len(matr))
	for i, arr := range matr {
		ret[i] = make([]sudokuElement, len(arr))
		for j, v := range arr {
			ret[i][j] = sudokuElement{
				value:        v,
				isStartValue: v != 0,
			}
		}
	}
	return ret
}

type sudoku struct {
	//Stats
	currentCol        uint8 //Holds the current column
	currentRow        uint8 //Holds the current row
	iteration         uint8 //How many iteration have been done
	startMissingValue uint8 //Holds the number of missing values at the beginning of the sudoku
	//Data
	grid [][]sudokuElement // Slice which hold all sudoku table.
	//Configuration
	size            uint8     // Size of a side of the grid. Ex: 9 for a 9x9 grid.
	shapes          [][]uint8 // Represent all the shape of the grid. Each shape is made by the element's position into the grid.
	requiredNumbers []uint8   //All the value that must be present in shape/column/row.
}

// Since sudoku is a private struct, the only way to create a new sudoku is to use the SudokuFactory function.
// This function is used to create a new sudoku, with the related configuration. It's a workaround due to golang's lack of constructor.
func SudokuFactory(grid [][]uint8) (s sudoku) {
	configuration := config.GetSudokuConfig()

	var elApp sudokuElement

	s = sudoku{
		currentCol:      0,
		currentRow:      0,
		iteration:       0,
		grid:            elApp.Fromuint8(grid),
		size:            uint8(configuration.SquareSize),
		shapes:          configuration.Shapes,
		requiredNumbers: configuration.RequiredNumbers,
	}

	s.startMissingValue = s.countMissingValues()

	return
}

func (s *sudoku) Solve() {
	shouldContinue := true
	for s.iteration = 1; s.iteration <= s.size*s.size && shouldContinue; s.iteration++ {
		s.PrintGrid()
		s.findMissingValue()

		shouldContinue = !s.isComplete() && newFound != nil
	}

	s.PrintGrid()
}

// Return all the element in the current row.
func (s *sudoku) getRowElements(row uint8) []uint8 {
	ret := make([]uint8, s.size)

	for i := uint8(0); i < s.size; i++ {
		ret[i] = s.grid[row][i].value
	}

	return ret
}

// Return all the element in the current column.
func (s *sudoku) getColElements(col uint8) []uint8 {
	ret := make([]uint8, s.size)

	for i := uint8(0); i < s.size; i++ {
		ret[i] = s.grid[i][col].value
	}

	return ret
}

// Returns all the position of a given shape.
func (s *sudoku) getCurrentShapeElementPosition() []uint8 {
	sizeRoot := math.Sqrt(float64(s.size))

	row := uint8((math.Floor(float64(s.currentRow)/sizeRoot) * sizeRoot))
	col := uint8(math.Floor(float64(s.currentCol) / sizeRoot))

	return s.shapes[row+col]
}

// Return the row and column of the element at the given position.
func (s *sudoku) getCordinatesFromPosition(pos uint8) (row, col uint8) {
	return pos / s.size, pos % s.size
}

// Accoring to the current row and column, returns all the element of the current shape.
func (s *sudoku) getShapeElements() []uint8 {
	ret := make([]uint8, s.size)

	for i, element := range s.getCurrentShapeElementPosition() {
		row, col := s.getCordinatesFromPosition(element)
		ret[i] = s.grid[row][col].value
	}

	return ret
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

	//Remove all the value that doesn't fit to that cell.
	values := make([]uint8, len(s.requiredNumbers))
	copy(values, s.requiredNumbers)
	r, c, b := s.getRowElements(s.currentRow), s.getColElements(s.currentCol), s.getShapeElements()
	values = removeFromArray(values, r)
	values = removeFromArray(values, c)
	values = removeFromArray(values, b)

	//If only one possibility is left, than use that value.
	if len(values) == 1 {
		return values[0]
	}

	/* //Iterate over all the values and check if there is only one missing value.
	v := pogo.FilterArray(values, func(valueCandidates uint8) bool {
		cellToCheck := pogo.FilterArray(s.getCurrentShapeElementPosition(), func(pos uint8) bool {
			row, col := s.getCordinatesFromPosition(pos)
			return s.grid[row][col].value == 0
		})

		isRight := pogo.EveryInArray(cellToCheck, func(cell uint8) bool {
			row, col := s.getCordinatesFromPosition(cell)
			return s.isValid(row, col, valueCandidates)
		})

		return isRight
	})

	if len(v) == 1 {
		return values[0]
	}  */

	return 0
}

func (s *sudoku) isValid(row, col, num uint8) bool {
	return !pogo.ContainsArray(s.getRowElements(row), num) && !pogo.ContainsArray(s.getColElements(col), num) && !pogo.ContainsArray(s.getShapeElements(), num)
}

// If there's no 0 value into the grid, then the sudoku is complete.
func (s *sudoku) isComplete() bool {
	return pogo.EveryInArray(s.grid, func(row []sudokuElement) bool {
		flat := pogo.MapArray(row, func(e sudokuElement) uint8 {
			return e.value
		})
		return !pogo.ContainsArray(flat, 0)
	})
}

func (s *sudoku) findMissingValue() {
	newFound = nil

	for i, row := range s.grid {
		for j, value := range row {
			if value.value != 0 {
				continue
			}

			s.currentRow = uint8(i)
			s.currentCol = uint8(j)

			if v := s.findValue(); v != 0 {
				s.grid[i][j].value = v
				newFound = &s.grid[i][j]
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
			if value.value == 0 {
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
	var color string

	for i, row := range s.grid {
		for j, value := range row {
			//Print the current value
			if newFound == &s.grid[i][j] {
				color = superHighlightColor
			} else if value.isStartValue {
				color = baseColor
			} else {
				color = highlightColor
			}
			fmt.Printf("%s%d ", color, value.value)
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
