package main

import (
	"contentssecurity/pkg/table"
)

func main() {
	gradeTable, _ := table.ImportTable("data/seiseki.txt")
	weightTable, _ := table.ImportTable("data/omomi.txt")
	minGradeTable, _ := table.ImportTable("data/saiteiten.txt")

	aptitudeTable := gradeTable.Mul(weightTable)
	resultTable := aptitudeTable.Evaluate(minGradeTable)

	aptitudeTable.ExportTable("bin/tekisei.txt", "適性")
	resultTable.ExportResult("bin/kekka.txt", "合否")
}
