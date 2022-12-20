package main

import (
	"contentssecurity/pkg/table"

	conn "github.com/uecconsecexp/secexp2022/se_go/connector"
)

func main() {
	addr := "0.0.0.0"

	client, err := conn.NewChugakuClient(addr)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Import seiseki table
	A, _ := table.Import("data/seiseki.txt")

	// Create random matrix M
	MRows, MCols := A.ColSize(), A.ColSize()
	M := table.Random(MRows, MCols)

	// Send M to yobikou
	if err = M.ChugakuSend(&client); err != nil {
		panic(err)
	}

	// Slice M
	MLeft, MRight := M.Slice(0, MRows, 0, MCols/2), M.Slice(0, MRows, MCols/2, MCols)

	// Calculate and send A'
	APrime := A.Mul(MLeft)
	if err = APrime.ChugakuSend(&client); err != nil {
		panic(err)
	}

	// Wait and receive B'
	BPrime, err := table.ChugakuReceive(&client)
	if err != nil {
		panic(err)
	}

	// Calculate and send A''
	APrimePrime := A.Mul(MRight).Mul(BPrime)
	if err = APrimePrime.ChugakuSend(&client); err != nil {
		panic(err)
	}

	// Wait and receive result
	result, err := table.ChugakuReceive(&client)
	if err != nil {
		panic(err)
	}

	// Export result
	result.ExportResult("data/out/kekka.txt", "合否")
}
