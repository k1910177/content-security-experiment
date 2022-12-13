package main

import (
	"contentssecurity/pkg/table"
	"fmt"
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

func yobikouSide() {
	yobikou, err := conn.NewYobikouServer()
	if err != nil {
		panic(err)
	}
	defer yobikou.Close()

	B, _ := table.ImportTable("data/omomi.txt")
	C, _ := table.ImportTable("data/saiteiten.txt")

	// Wait and receive for byte array
	m, err := yobikou.Receive()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received: %s\n", m)

	// Send table
	matrix := [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	err = yobikou.SendTable(matrix)
	if err != nil {
		panic(err)
	}
}

func chugakuSide(addr string) {
	chugaku, err := conn.NewChugakuClient(addr)
	if err != nil {
		panic(err)
	}
	defer chugaku.Close()

	// Send byte array
	err = chugaku.Send([]byte("ping"))
	if err != nil {
		panic(err)
	}

	// Wait and receive table
	matrix, err := chugaku.ReceiveTable()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received: %v\n", matrix)
}
