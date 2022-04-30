package internal

import (
	"fmt"
	"math"

	"github.com/osharko/sudoku/internal/config"
	"github.com/osharko/sudoku/internal/pogo"
)

var (
	baseColor      = "\033[32m"
	highlightColor = "\033[33m"
)

type sudokuSize uint8

func (*sudokuSize) ToArray(arr []uint8) []sudokuSize {
	ret := make([]sudokuSize, len(arr))
	for i, v := range arr {
		ret[i] = sudokuSize(v)
	}
	return ret
}

func (*sudokuSize) ToMatrix(matr [][]uint8) [][]sudokuSize {
	ret := make([][]sudokuSize, len(matr))
	for i, arr := range matr {
		ret[i] = make([]sudokuSize, len(arr))
		for j, v := range arr {
			ret[i][j] = sudokuSize(v)
		}
	}
	return ret
}

type sudokuElement struct {
	value        sudokuSize
	isStartValue bool
}

func (e *sudokuElement) FromSudokusize(matr [][]sudokuSize) [][]sudokuElement {
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
	currentCol        sudokuSize //Holds the current column
	currentRow        sudokuSize //Holds the current row
	iteration         sudokuSize //How many iteration have been done
	startMissingValue sudokuSize //Holds the number of missing values at the beginning of the sudoku
	//Data
	grid [][]sudokuElement // Slice which hold all sudoku table.
	//Configuration
	size            sudokuSize     // Size of a side of the grid. Ex: 9 for a 9x9 grid.
	shapes          [][]sudokuSize // Represent all the shape of the grid. Each shape is made by the element's position into the grid.
	requiredNumbers []sudokuSize   //All the value that must be present in shape/column/row.
}

// Since sudoku is a private struct, the only way to create a new sudoku is to use the SudokuFactory function.
// This function is used to create a new sudoku, with the related configuration. It's a workaround due to golang's lack of constructor.
func SudokuFactory(grid [][]uint8) (s sudoku) {
	configuration := config.GetSudokuConfig()

	var typeApp sudokuSize
	var elApp sudokuElement

	s = sudoku{
		currentCol:      sudokuSize(0),
		currentRow:      sudokuSize(0),
		iteration:       sudokuSize(0),
		grid:            elApp.FromSudokusize(typeApp.ToMatrix(grid)),
		size:            sudokuSize(configuration.SquareSize),
		shapes:          typeApp.ToMatrix(configuration.Shapes),
		requiredNumbers: typeApp.ToArray(configuration.RequiredNumbers),
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
func (s *sudoku) getRowElements() (ret []sudokuSize) {
	for i := sudokuSize(0); i < s.size; i++ {
		ret = append(ret, s.grid[s.currentRow][i].value)
	}

	return
}

// Return all the element in the current column.
func (s *sudoku) getColElements() (ret []sudokuSize) {
	for i := sudokuSize(0); i < s.size; i++ {
		ret = append(ret, s.grid[i][s.currentCol].value)
	}

	return
}

// Accoring to the current row and column, returns all the element of the current shape.
func (s *sudoku) getShapeElements() []sudokuSize {
	sizeRoot := math.Sqrt(float64(s.size))

	row := sudokuSize((math.Floor(float64(s.currentRow)/sizeRoot) * sizeRoot))
	col := sudokuSize(math.Floor(float64(s.currentCol) / sizeRoot))

	return s.getShape(row + col)
}

// Returns all element in a given shape.
func (s *sudoku) getShape(shapePos sudokuSize) (ret []sudokuSize) {
	for _, element := range s.shapes[shapePos] {
		row := (element) / s.size
		col := (element) % s.size
		ret = append(ret, s.grid[row][col].value)
	}

	return
}

// Find all the missing number in the current row, column and shape,
// Thene check if those 3 has 1 common missing number,
// If so, fill the current cell with that number.
func (s *sudoku) findValue() sudokuSize {
	removeFromArray := func(origin []sudokuSize, comparison []sudokuSize) []sudokuSize {
		rem := make([]sudokuSize, 0)

		for _, v := range origin {
			for _, value := range comparison {
				if v == value {
					rem = append(rem, v)
				}
			}
		}

		return pogo.FilterArray(origin, func(o sudokuSize) bool {
			return !pogo.ContainsArray(rem, o)
		})
	}

	values := make([]sudokuSize, len(s.requiredNumbers))
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
	return pogo.EveryInArray(s.grid, func(row []sudokuElement) bool {
		flat := pogo.MapArray(row, func(e sudokuElement) sudokuSize {
			return e.value
		})
		return !pogo.ContainsArray(flat, 0)
	})
}

func (s *sudoku) findMissingValue() {
	for i, row := range s.grid {
		for j, value := range row {
			if value.value != 0 {
				continue
			}

			s.currentRow = sudokuSize(i)
			s.currentCol = sudokuSize(j)

			if v := s.findValue(); v != 0 {
				s.grid[i][j].value = v
				return
			}
		}
	}
}

// Return the number of 0 into the grid.
func (s *sudoku) countMissingValues() sudokuSize {
	count := sudokuSize(0)
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
			if value.isStartValue {
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
