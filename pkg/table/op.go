package table

import "gonum.org/v1/gonum/mat"

func (a *Table) Add(b *Table) *Table {
	var m mat.Dense
	m.Add(a.values, b.values)

	return &Table{
		values: &m,
	}
}

func (a *Table) Mul(b *Table) *Table {
	var m mat.Dense
	m.Mul(a.values, b.values)

	return &Table{
		values: &m,
	}
}

func (a *Table) Inv() *Table {
	var m mat.Dense
	m.Inverse(a.values)

	return &Table{
		values: &m,
	}
}

func (a *Table) Slice(rowStart, rowEnd, colStart, colEnd int) *Table {
	m := a.values.Slice(rowStart, rowEnd, colStart, colEnd)

	return &Table{
		values: mat.DenseCopyOf(m),
	}
}

func (a *Table) Evaluate(b *Table) *Table {
	rowSize, colSize := a.values.Caps()
	result := Table{
		values: mat.NewDense(rowSize, colSize, nil),
	}

	for rowIndex := 0; rowIndex < rowSize; rowIndex++ {
		for colIndex := 0; colIndex < colSize; colIndex++ {
			if a.values.At(rowIndex, colIndex) >= b.values.At(0, colIndex) {
				result.values.Set(rowIndex, colIndex, 1)
			} else {
				result.values.Set(rowIndex, colIndex, 0)
			}
		}
	}

	return &result
}
