package internal

import (
	"fmt"
	"math"

	"github.com/osharko/sudoku/internal/config"
)

type sudoku struct {
	//Stats
	currentCol          uint8 //Holds the current column
	currentRow          uint8 //Holds the current row
	iteration           uint8 //How many iteration have been done
	startMissingValue   uint8 //Holds the number of missing values at the beginning of the sudoku
	currentMissingValue uint8 //Holds the number of missing values at the current iteration
	//Data
	grid []uint8 // Slice which hold all the rows as single slice. That should simplify the operation when searching by col.
	//Configuration
	size            uint8     // Size of a side of the grid. Ex: 9 for a 9x9 grid.
	shapes          [][]uint8 // Represent all the shape of the grid. Each shape is made by the element's position into the grid.
	requiredNumbers []uint8   //All the value that must be present in shape/column/row.
}

// Since sudoku is a private struct, the only way to create a new sudoku is to use the SudokuFactory function.
// This function is used to create a new sudoku, with the related configuration. It's a workaround due to golang's lack of constructor.
func SudokuFactory() (s sudoku) {
	configuration := config.GetSudokuConfig()
	grid := config.GetSudokuDataEntry().GetLinearSlice()

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
	s.currentMissingValue = s.startMissingValue

	return
}

// Return all the element in the current row.
func (s *sudoku) GetRowElements() []uint8 {
	return s.grid[(s.size * s.currentRow):((s.size * s.currentRow) + s.size)]
}

// Return all the element in the current column.
func (s *sudoku) GetColElements() []uint8 {
	ret := make([]uint8, s.size)

	for i := uint8(0); i < s.size; i++ {
		ret[i] = s.grid[(i*s.size)+s.currentCol]
	}

	return ret
}

// Accoring to the current row and column, returns all the element of the current shape.
func (s *sudoku) GetShapeElements() []uint8 {
	sizeRoot := math.Sqrt(float64(s.size))

	row := uint8((math.Floor(float64(s.currentRow)/sizeRoot) * sizeRoot))
	col := uint8(math.Round(float64(s.currentCol) / sizeRoot))

	return s.getShape(row + col)
}

// Returns all element in a given shape.
func (s *sudoku) getShape(shapePos uint8) []uint8 {
	ret := make([]uint8, s.size)

	for i, element := range s.shapes[shapePos] {
		ret[i] = s.grid[element]
	}

	return ret
}

// Find all the missing number in the current row, column and shape,
// Thene check if those 3 has 1 common missing numbe,
// If so, then we can fill the current cell with that number.
func (s *sudoku) FindValue() *uint8 {
	missingRow := s.getMissingNumber(s.GetRowElements())
	missingCol := s.getMissingNumber(s.GetColElements())
	missingShape := s.getMissingNumber(s.GetShapeElements())

	if len(missingCol) > 0 && len(missingRow) > 0 && len(missingShape) > 0 {
		values := make(map[uint8]bool)
		duplicates := make(map[uint8]bool)

		findDuplicates := func(slice []uint8) {
			for _, value := range slice {
				if values[value] {
					duplicates[value] = true
				}
				values[value] = true
			}
		}

		findDuplicates(missingRow)
		findDuplicates(missingCol)
		findDuplicates(missingShape)

		if len(duplicates) == 1 {
			for key := range duplicates {
				return &key
			}
		}

		//fmt.Println("Duplicates: ", duplicates)
	}

	//fmt.Println("Missing number in row: ", missingRow)
	//fmt.Println("Missing number in col: ", missingCol)
	//fmt.Println("Missing number in shape: ", missingShape)

	return nil
}

// Returns all the missing number, from sudoku.requiredNumbers, in the given slice.
func (s *sudoku) getMissingNumber(slice []uint8) []uint8 {
	missing := make([]uint8, 0)

	for _, number := range s.requiredNumbers {
		if !contains(slice, number) {
			missing = append(missing, number)
		}
	}

	return missing
}

// If there's no 0 value into the grid, then the sudoku is complete.
func (s *sudoku) IsComplete() bool {
	return !contains(s.grid, 0)
}

func (s *sudoku) FindMissingValue() {
	s.iteration++
	s.PrintGrid()

	for _, value := range s.grid {
		if value == 0 {
			continue
		}

		if v := s.FindValue(); v != nil {
			fmt.Println(v)
		}
	}
}

// Return the number of 0 into the grid.
func (s *sudoku) countMissingValues() uint8 {
	count := uint8(0)
	for _, value := range s.grid {
		if value == 0 {
			count++
		}
	}

	return count
}

func contains(slice []uint8, element uint8) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}

func (s *sudoku) PrintGrid() {
	fmt.Printf("\nCurrent Iteration: %d\tMissing Values: %d\tStarting Missing Values: %d\n\n", s.iteration, s.currentMissingValue, s.startMissingValue)
	fmt.Printf("\n\t\t\t")

	root := uint8(math.Sqrt(float64(s.size)))

	for i, value := range s.grid {
		// +1 Because the grid is 0 based
		index := uint8(i) + 1

		//Print the current value
		fmt.Printf("%d ", value)

		//If we reach the end of a column, we need to go to the next row.
		if index%root == 0 && i != 0 {
			fmt.Printf("\t")
		}

		// Print a new line every time we reach the end of a row.
		if index%s.size == 0 {
			// Print an additional line, at the end of the shape
			if index%(s.size*root) == 0 {
				fmt.Println()
			}
			fmt.Printf("\n\t\t\t")
		}
	}
	fmt.Printf("\n\n")
}
