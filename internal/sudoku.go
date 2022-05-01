package internal

import (
	"fmt"
	"math"

	"github.com/osharko/sudoku/internal/config"
	"github.com/osharko/sudoku/internal/pogo"
)

var (
	primaryColor        = "\033[32m"
	secondaryColor      = "\033[33m"
	accentColor         = "\033[31m"
	superHighlightColor = "\033[35m"

	printEachStep bool
)

type sudokuElement struct {
	value        uint8
	isStartValue bool
}

// Well, it's simple and doesn't need explanations.
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
	currentCol        uint8          //Holds the current column
	currentRow        uint8          //Holds the current row
	iteration         uint8          //How many iteration have been done
	newFound          *sudokuElement //Holds the last found value
	startMissingValue uint8          //Holds the number of missing values at the beginning of the sudoku
	//Data
	grid [][]sudokuElement // Slice which hold all sudoku table.
	//Configuration
	size            uint8     // Size of a side of the grid. Ex: 9 for a 9x9 grid.
	shapes          [][]uint8 // Represent all the shape of the grid. Each shape is made by the element's position into the grid.
	requiredNumbers []uint8   //All the value that must be present in shape/column/row.
}

// Since sudoku is a private struct, the only way to create a new sudoku is to use the SudokuFactory function.
// This function is used to create a new sudoku, with the related configuration. It's a workaround due to golang's lack of constructor.
func SudokuFactory(grid [][]uint8, printStep bool) (s sudoku) {
	configuration := config.GetSudokuConfig()

	var elApp sudokuElement
	printEachStep = printStep

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

//Unlease the beast and solve the sudoku.
func (s *sudoku) Solve() {
	shouldGoAhead := true
	// Continue solving until the sudoku is complete or there's no more found value.
	for s.iteration = 1; shouldGoAhead; s.iteration++ {
		if printEachStep {
			s.printGrid(&s.iteration)
		}
		s.findMissingValue()

		shouldGoAhead = !s.isComplete() && s.newFound != nil
	}
	//Just a last winning print.
	s.printGrid(nil)
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
	possibleValues := s.getMissingValues()

	if nakedSingle := s.nakedSingle(possibleValues); nakedSingle != 0 {
		fmt.Println("ioZ")
		return nakedSingle
	} else if hiddenSingle := s.hiddenSingle(possibleValues); len(hiddenSingle) == 1 {
		return hiddenSingle[0]
	} else {
		return 0
	}
}

// Return the only present value, if more then one exists, then return noone (0).
func (s *sudoku) nakedSingle(possibilities []uint8) uint8 {
	if len(possibilities) == 1 {
		return possibilities[0]
	}
	return 0
}

//Iterate over all the missing values for the row/column/shape,
// and check if some of these value has only one suitable cell inside the shape.
func (s *sudoku) hiddenSingle(possibilities []uint8) []uint8 {
	return pogo.FilterArray(possibilities, s.isCellTheOnlySuitable)
}

// Check if a value has only one suitable cell inside the shape.
func (s *sudoku) isCellTheOnlySuitable(possibilities uint8) bool {
	// Get the the position of the cell, in the current shape, which has no value, excluding the current cell.
	positionToCheck := pogo.FilterArray(s.getCurrentShapeElementPosition(), func(pos uint8) bool {
		row, col := s.getCordinatesFromPosition(pos)
		//Don't check the current cell.
		return s.grid[row][col].value == 0 && !(row == s.currentRow && col == s.currentCol)
	})
	//Check if the current "valueCandidate" is applicable only to the current cell, no other cell must be a valid candidate to that value.
	return pogo.EveryInArray(positionToCheck, func(cell uint8) bool {
		row, col := s.getCordinatesFromPosition(cell)
		return !s.isValid(row, col, possibilities)
	})
}

// Get all the values not present in row, col, or shape.
func (s *sudoku) getMissingValues() []uint8 {
	presentValues := make(map[uint8]bool)

	sendToMap := func(values []uint8) {
		for _, v := range values {
			presentValues[v] = true
		}
	}

	sendToMap(s.getRowElements(s.currentRow))
	sendToMap(s.getColElements(s.currentCol))
	sendToMap(s.getShapeElements())

	return pogo.FilterArray(s.requiredNumbers, func(v uint8) bool {
		return !presentValues[v]
	})
}

//Check if the given value is valid for the given cell.
func (s *sudoku) isValid(row, col, num uint8) bool {
	return !pogo.ContainsArray(s.getRowElements(row), num) && !pogo.ContainsArray(s.getColElements(col), num)
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

//Iterate one time through the grid, and try to find as much new value as it can.
func (s *sudoku) findMissingValue() {
	//Reset the new Found value.
	s.newFound = nil

	//Iterate over the grid.
	for i, row := range s.grid {
		//Iterate over the row
		for j, value := range row {
			//Only 0 must be considered, otherwise the value has already been found.
			if value.value != 0 {
				continue
			}
			//Set row and col of the current cell.
			s.currentRow = uint8(i)
			s.currentCol = uint8(j)
			//Get the missing value. If found, set to the grid, otherwise continue scanning and finding.
			if v := s.findValue(); v != 0 {
				s.grid[i][j].value = v
				s.newFound = &s.grid[i][j]
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

//Print the grid with some stats, and different colors for different cell's status.
func (s *sudoku) printGrid(iteration *uint8) {
	head := ""
	if iteration != nil {
		head = fmt.Sprintf("Iteration %d\tMissing Values: %d\tStarting Missing Values: %d", *iteration, s.countMissingValues(), s.startMissingValue)
	} else if iteration == nil && s.isComplete() {
		head = "\t\t\t\tSolved"
	} else {
		head = fmt.Sprintf("\t\tFailed after %d iterations, found %d values.", s.iteration, s.startMissingValue-s.countMissingValues())
	}

	fmt.Printf("\n%s\n\n", head)
	fmt.Printf("\n\t\t\t")

	root := int(math.Sqrt(float64(s.size)))
	var color string

	for i, row := range s.grid {
		for j, value := range row {
			//Print the current value
			if s.newFound == &s.grid[i][j] {
				color = superHighlightColor
			} else if value.isStartValue {
				color = primaryColor
			} else if s.grid[i][j].value == 0 {
				color = secondaryColor
			} else {
				color = accentColor
			}
			fmt.Printf("%s%d%s ", color, value.value, primaryColor)
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
	fmt.Printf("\n-------------------------------------------------------------------\n\n")
}
