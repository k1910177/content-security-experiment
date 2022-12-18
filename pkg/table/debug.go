package table

import (
	"fmt"
)

func (table *Table) Print(title string) {
	fmt.Println(title)
	fmt.Printf("Rows: %v\n", table.Rows)
	fmt.Printf("Cols: %v\n", table.Cols)

	rowSize, colSize := table.values.Caps()
	for rowIndex := 0; rowIndex < rowSize; rowIndex++ {
		for colIndex := 0; colIndex < colSize; colIndex++ {
			fmt.Printf("%-8.2f", table.values.At(rowIndex, colIndex))
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n")
}
