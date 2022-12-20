package main

import (
	"contentssecurity/pkg/table"

	conn "github.com/uecconsecexp/secexp2022/se_go/connector"
)

func main() {
	server, err := conn.NewYobikouServer()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	B, _ := table.Import("data/omomi.txt")
	C, _ := table.Import("data/saiteiten.txt")

	// Wait and receive M
	M, err := table.YobikouReceive(&server)
	if err != nil {
		panic(err)
	}

	// Calculate M' and slice
	MInv := M.Inv()
	MRows, MCols := MInv.RowSize(), MInv.ColSize()
	MInvTop, MInvBottom := MInv.Slice(0, MRows/2, 0, MCols), MInv.Slice(MRows/2, MRows, 0, MCols)

	// Wait and receive A'
	APrime, err := table.YobikouReceive(&server)
	if err != nil {
		panic(err)
	}

	// Calculate and send B'
	BPrime := MInvBottom.Mul(B)
	if err = BPrime.YobikouSend(&server); err != nil {
		panic(err)
	}

	// Calculate and send B''
	BPrimePrime := APrime.Mul(MInvTop).Mul(B)

	// Wait and receive A''
	APrimePrime, err := table.YobikouReceive(&server)
	if err != nil {
		panic(err)
	}

	// Calculate and export aptitude
	aptitude := APrimePrime.Add(BPrimePrime)
	aptitude.Export("data/out/tekisei.txt", "適性")

	// Evaluate and send result
	result := aptitude.Evaluate(C)
	if err = result.YobikouSend(&server); err != nil {
		panic(err)
	}
}
