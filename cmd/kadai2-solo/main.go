package main

import (
	"contentssecurity/pkg/table"

	conn "github.com/uecconsecexp/secexp2022/se_go/connector"
)

func main() {
	// Chugaku
	A, _ := table.Import("data/seiseki.txt")

	// Yobikou
	B, _ := table.Import("data/omomi.txt")
	C, _ := table.Import("data/saiteiten.txt")

	// Create random matrix M
	MRows, MCols := len(A.Cols), len(B.Rows)
	M := table.Random(MRows, MCols)

	// Send M and calculate M inverse
	MInv := M.Inv()

	// Slice M
	MLeft, MRight := M.Slice(0, MRows, 0, MCols/2), M.Slice(0, MRows, MCols/2, MCols)
	MInvTop, MInvBottom := MInv.Slice(0, MRows/2, 0, MCols), MInv.Slice(MRows/2, MRows, 0, MCols)

	// Calculate and send A' and B'
	APrime := A.Mul(MLeft)
	BPrime := MInvBottom.Mul(B)

	// Calculate and send A'' and B''
	APrimePrime := A.Mul(MRight).Mul(BPrime)
	BPrimePrime := APrime.Mul(MInvTop).Mul(B)

	// Calculate aptitude and evaulate result
	aptitudeTable := APrimePrime.Add(BPrimePrime)
	resultTable := aptitudeTable.Evaluate(C)

	// Export
	aptitudeTable.Export("bin/tekisei.txt", "適性")
	resultTable.ExportResult("bin/kekka.txt", "合否")
}

func yobikouSide() {
	yobikou, err := conn.NewYobikouServer()
	if err != nil {
		panic(err)
	}
	defer yobikou.Close()
}

func chugakuSide(addr string) {
	chugaku, err := conn.NewChugakuClient(addr)
	if err != nil {
		panic(err)
	}
	defer chugaku.Close()

}
