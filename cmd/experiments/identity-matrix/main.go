package main

import "contentssecurity/pkg/table"

func main() {
	// Chugaku: Import seiseki
	A, _ := table.Import("data/seiseki.txt")

	// Yobikou: Import omomi and saiteiten
	B, _ := table.Import("data/omomi.txt")
	C, _ := table.Import("data/saiteiten.txt")

	// Chugaku: create and send M (identity matrix)
	MRows, MCols := A.ColSize(), A.ColSize()
	MSlice := make([][]float64, MRows)
	for rowIndex := 0; rowIndex < MRows; rowIndex++ {
		MSlice[rowIndex] = make([]float64, MCols)
		for colIndex := 0; colIndex < MCols; colIndex++ {
			if rowIndex == colIndex {
				MSlice[rowIndex][colIndex] = 1
			}
		}
	}
	M := table.New(MSlice)

	// Yobikou: Calculate M'
	MInv := M.Inv()

	// Chugaku: Slice M
	MLeft, MRight := M.Slice(0, MRows, 0, MCols/2), M.Slice(0, MRows, MCols/2, MCols)

	// Yobikou: Slice M'
	MInvTop, MInvBottom := MInv.Slice(0, MRows/2, 0, MCols), MInv.Slice(MRows/2, MRows, 0, MCols)

	// Chugaku: Calculate and send A'
	APrime := A.Mul(MLeft)

	// Yobikou: Calculate and send B'
	BPrime := MInvBottom.Mul(B)

	// Chugaku: Calculate and send A''
	APrimePrime := A.Mul(MRight).Mul(BPrime)

	// Yobikou: Calculate B''
	BPrimePrime := APrime.Mul(MInvTop).Mul(B)

	APrime.Print("APrime")
	BPrime.Print("BPrime")
	APrimePrime.Print("APrimePrime")
	BPrimePrime.Print("BPrimePrime")

	// Yobikou: Calculate aptitude and evaulate result
	aptitude := APrimePrime.Add(BPrimePrime)
	result := aptitude.Evaluate(C)

	// Export
	aptitude.Print("tekisei")
	result.Print("kekka")
}
