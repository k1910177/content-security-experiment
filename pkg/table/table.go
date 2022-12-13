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
	Values *mat.Dense
	Rows   []string
	Cols   []string
}

func RandomTable(rowSize, colSize int) *Table {
	table := Table{
		Values: mat.NewDense(rowSize, colSize, nil),
		Rows:   make([]string, rowSize),
		Cols:   make([]string, colSize),
	}

	for rowIndex := 0; rowIndex < rowSize; rowIndex++ {
		for colIndex := 0; colIndex < colSize; colIndex++ {
			table.Values.Set(rowIndex, colIndex, rand.Float64())
		}
	}

	return &table
}

func ImportTable(path string) (*Table, error) {
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
		Values: mat.NewDense(rowSize, colSize, nil),
		Rows:   make([]string, rowSize),
		Cols:   make([]string, colSize),
	}

	for rowIndex, rowItems := range records {
		for colIndex, item := range rowItems {
			if colIndex == 0 && rowIndex == 0 {
				continue
			}
			if rowIndex == 0 && colIndex > 0 {
				table.Cols[colIndex-1] = item
				continue
			}
			if colIndex == 0 && rowIndex > 0 {
				table.Rows[rowIndex-1] = item
				continue
			}

			value, err := strconv.ParseFloat(item, 64)
			if err != nil {
				return nil, err
			}
			table.Values.Set(rowIndex-1, colIndex-1, value)
		}
	}

	return &table, nil
}

func (table *Table) ExportTable(path string, title string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	rowSize, colSize := table.Values.Caps()
	for rowIndex := 0; rowIndex < rowSize+1; rowIndex++ {
		record := make([]string, colSize+1)
		for colIndex := range record {
			if colIndex == 0 && rowIndex == 0 {
				record[colIndex] = title
				continue
			}
			if rowIndex == 0 && colIndex > 0 {
				record[colIndex] = table.Cols[colIndex-1]
				continue
			}
			if colIndex == 0 && rowIndex > 0 {
				record[0] = table.Rows[rowIndex-1]
				continue
			}
			value := table.Values.At(rowIndex-1, colIndex-1)
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

	rowSize, colSize := table.Values.Caps()
	for rowIndex := 0; rowIndex < rowSize+1; rowIndex++ {
		record := make([]string, colSize+1)
		for colIndex := range record {
			if colIndex == 0 && rowIndex == 0 {
				record[colIndex] = title
				continue
			}
			if rowIndex == 0 && colIndex > 0 {
				record[colIndex] = table.Cols[colIndex-1]
				continue
			}
			if colIndex == 0 && rowIndex > 0 {
				record[0] = table.Rows[rowIndex-1]
				continue
			}

			if table.Values.At(rowIndex-1, colIndex-1) == 0 {
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