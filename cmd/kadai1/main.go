package main

import (
	"contentssecurity/pkg/table"
)

func main() {
	gradeTable, _ := table.Import("data/seiseki.txt")
	weightTable, _ := table.Import("data/omomi.txt")
	minGradeTable, _ := table.Import("data/saiteiten.txt")

	aptitudeTable := gradeTable.Mul(weightTable)
	resultTable := aptitudeTable.Evaluate(minGradeTable)

	aptitudeTable.Export("data/out/tekisei.txt", "適性")
	resultTable.ExportResult("data/out/kekka.txt", "合否")
}
