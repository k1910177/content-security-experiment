package table

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func (table *Table) Print(title string) {
	fc := mat.Formatted(table.Values, mat.Prefix(""), mat.Squeeze())
	fmt.Println(title)
	fmt.Printf("%v\n\n", fc)
}
