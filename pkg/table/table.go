package table

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

type Table struct {
	values *mat.Dense
}

func New(values [][]float64) *Table {
	rowSize, colSize := len(values), len(values[0])
	records := make([]float64, rowSize*colSize)
	for rowIndex := 0; rowIndex < rowSize; rowIndex++ {
		for colIndex := 0; colIndex < colSize; colIndex++ {
			records[rowIndex*colSize+colIndex] = values[rowIndex][colIndex]
		}
	}

	return &Table{
		values: mat.NewDense(rowSize, colSize, records),
	}
}

func Random(rowSize, colSize int) *Table {
	table := Table{
		values: mat.NewDense(rowSize, colSize, nil),
	}

	for rowIndex := 0; rowIndex < rowSize; rowIndex++ {
		for colIndex := 0; colIndex < colSize; colIndex++ {
			table.values.Set(rowIndex, colIndex, rand.Float64())
		}
	}

	return &table
}

func Import(path string) (*Table, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	rowSize, colSize := len(records)-1, len(records[0])-1
	table := Table{
		values: mat.NewDense(rowSize, colSize, nil),
	}

	for rowIndex, rowItems := range records {
		for colIndex, item := range rowItems {
			if colIndex == 0 || rowIndex == 0 {
				continue
			}

			value, err := strconv.ParseFloat(item, 64)
			if err != nil {
				return nil, err
			}
			table.values.Set(rowIndex-1, colIndex-1, value)
		}
	}

	return &table, nil
}

func (table *Table) Export(path string, title string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	rowSize, colSize := table.values.Caps()
	for rowIndex := 0; rowIndex < rowSize; rowIndex++ {
		record := make([]string, colSize)
		for colIndex := range record {
			value := table.values.At(rowIndex, colIndex)
			record[colIndex] = fmt.Sprintf("%.2f", value)
		}

		err := writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (table *Table) ExportResult(path string, title string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	rowSize, colSize := table.values.Caps()
	for rowIndex := 0; rowIndex < rowSize; rowIndex++ {
		record := make([]string, colSize)
		for colIndex := range record {
			if table.values.At(rowIndex, colIndex) == 0 {
				record[colIndex] = "否"
			} else {
				record[colIndex] = "合"
			}
		}

		err := writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (table *Table) ToSlice() [][]float64 {
	rowSize, colSize := table.values.Caps()
	records := make([][]float64, rowSize)
	for rowIndex := 0; rowIndex < rowSize; rowIndex++ {
		records[rowIndex] = make([]float64, colSize)
		for colIndex := 0; colIndex < colSize; colIndex++ {
			records[rowIndex][colIndex] = table.values.At(rowIndex, colIndex)
		}
	}
	return records
}

func (table *Table) RowSize() int {
	r, _ := table.values.Caps()
	return r
}

func (table *Table) ColSize() int {
	_, c := table.values.Caps()
	return c
}
