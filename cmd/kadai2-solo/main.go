package main

import (
	"contentssecurity/pkg/table"
	"time"

	conn "github.com/uecconsecexp/secexp2022/se_go/connector"
)

func main() {
	yc := make(chan bool)
	go func() {
		yobikouSide()
		yc <- true
	}()

	time.Sleep(500 * time.Millisecond)

	cc := make(chan bool)
	go func() {
		chugakuSide("0.0.0.0")
		cc <- true
	}()

	<-yc
	<-cc
}

func chugakuSide(addr string) {
	client, err := conn.NewChugakuClient(addr)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Import seiseki table
	A, _ := table.Import("data/seiseki.txt")

	// Create random matrix M
	MRows, MCols := len(A.Cols), len(A.Cols)
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
	BPrime, err := table.ChugakuReceive(&client, nil, nil)
	if err != nil {
		panic(err)
	}

	// Calculate and send A''
	APrimePrime := A.Mul(MRight).Mul(BPrime)
	if err = APrimePrime.ChugakuSend(&client); err != nil {
		panic(err)
	}

	// Wait and receive result
	result, err := table.ChugakuReceive(&client, nil, nil)
	if err != nil {
		panic(err)
	}

	// Export result
	result.ExportResult("data/out/kekka.txt", "合否")
}

func yobikouSide() {
	server, err := conn.NewYobikouServer()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	B, _ := table.Import("data/omomi.txt")
	C, _ := table.Import("data/saiteiten.txt")

	// Wait and receive M
	M, err := table.YobikouReceive(&server, nil, nil)
	if err != nil {
		panic(err)
	}

	// Calculate M' and slice
	MInv := M.Inv()
	MRows, MCols := len(M.Rows), len(M.Cols)
	MInvTop, MInvBottom := MInv.Slice(0, MRows/2, 0, MCols), MInv.Slice(MRows/2, MRows, 0, MCols)

	// Wait and receive A'
	APrime, err := table.YobikouReceive(&server, nil, nil)
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
	APrimePrime, err := table.YobikouReceive(&server, nil, nil)
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
