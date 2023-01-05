package main

import (
	"contentssecurity/pkg/table"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
	"math/rand"
	"time"

	conn "github.com/uecconsecexp/secexp2022/se_go/connector"
)

func main() {
	yc := make(chan bool)
	go func() {
		yobikouSide()
		yc <- true
	}()

	time.Sleep(100 * time.Millisecond)

	cc := make(chan bool)
	go func() {
		chugakuSide("0.0.0.0")
		cc <- true
	}()

	<-yc
	<-cc
}

func yobikouSide() {
	server, err := conn.NewYobikouServer()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	B, _ := table.Import("data/omomi.txt")
	C, _ := table.Import("data/saiteiten.txt")

	// Wait and receive key
	key, err := server.Receive()
	if err != nil {
		panic(err)
	}

	// Create M
	MRows, MCols := B.RowSize(), B.RowSize()
	M := hashTable(MRows, MCols, key)

	// Calculate M' and slice
	MInv := M.Inv()
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

func chugakuSide(addr string) {
	client, err := conn.NewChugakuClient(addr)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Import seiseki table
	A, _ := table.Import("data/seiseki.txt")

	// Create random key
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		panic(err)
	}

	// Send key to yobikou
	if err = client.Send(key); err != nil {
		panic(err)
	}

	// Create matrix M
	MRows, MCols := A.ColSize(), A.ColSize()
	M := hashTable(MRows, MCols, key)

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

func hashTable(rowSize, colSize int, seed []byte) *table.Table {
	slice := make([][]float64, rowSize)
	seedInt := int(binary.LittleEndian.Uint64(seed))
	for rowIndex := 0; rowIndex < rowSize; rowIndex++ {
		slice[rowIndex] = make([]float64, colSize)
		for colIndex := 0; colIndex < colSize; colIndex++ {
			data := seedInt + rowIndex*colSize + colIndex
			bytes := big.NewInt(int64(data)).Bytes()
			hash := sha256.Sum256(bytes)
			value := binary.LittleEndian.Uint64(hash[:])
			slice[rowIndex][colIndex] = float64(value)
		}
	}

	return table.New(slice)
}
