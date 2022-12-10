package table

import "gonum.org/v1/gonum/mat"

func (a *Table) Add(b *Table) *Table {
	var m mat.Dense
	m.Add(a.Values, b.Values)

	return &Table{
		Values: &m,
		Rows:   a.Rows,
		Cols:   b.Cols,
	}
}

func (a *Table) Mul(b *Table) *Table {
	var m mat.Dense
	m.Mul(a.Values, b.Values)

	return &Table{
		Values: &m,
		Rows:   a.Rows,
		Cols:   b.Cols,
	}
}

func (a *Table) Inv() *Table {
	var m mat.Dense
	m.Inverse(a.Values)

	return &Table{
		Values: &m,
		Rows:   a.Rows,
		Cols:   a.Cols,
	}
}

func (a *Table) Slice(rowStart, rowEnd, colStart, colEnd int) *Table {
	m := a.Values.Slice(rowStart, rowEnd, colStart, colEnd)

	return &Table{
		Values: mat.DenseCopyOf(m),
		Rows:   a.Rows[rowStart:rowEnd],
		Cols:   a.Cols[colStart:colEnd],
	}
}

func (a *Table) Evaluate(b *Table) *Table {
	rowSize, colSize := a.Values.Caps()
	result := Table{
		Rows:   a.Rows,
		Cols:   a.Cols,
		Values: mat.NewDense(rowSize, colSize, nil),
	}

	for rowIndex := 0; rowIndex < rowSize; rowIndex++ {
		for colIndex := 0; colIndex < colSize; colIndex++ {
			if a.Values.At(rowIndex, colIndex) >= b.Values.At(0, colIndex) {
				result.Values.Set(rowIndex, colIndex, 1)
			} else {
				result.Values.Set(rowIndex, colIndex, 0)
			}
		}
	}

	return &result
}
