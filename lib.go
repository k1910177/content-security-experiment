package contentssecurity

import (
	"fmt"

	conn "github.com/uecconsecexp/secexp2022/se_go/connector"
)

// ========== kadai1用 ==========

// main.goで使用したい機能はパッケージとして切り出しておくと便利です。
// lib.goに記述する必要はなく、別なファイルにしてもかまいません。
// ただしファイルの行頭はpackage contentssecurityとする必要があります。

// main.goで使用したい関数はパブリックとするため、大文字で始めます。
func Hello() string {
	return "Hello, world!"
}

// 便利なモジュールやその他必要な情報は README.md にまとめています。

// ========== kadai2用 ==========

// 実際に通信を行う機能もパッケージ側に書いてしまいましょう。
// YobikouSide、ChugakuSideの中身を書き換えてください。

func YobikouSide() {
	yobikou, err := conn.NewYobikouServer()
	// 例外処理: err != nilは必ず確認し、適切に処理してください。panicでかまいません。
	if err != nil {
		panic(err)
	}
	defer yobikou.Close() // TCP通信を終了するため、最後に必ず呼び出す必要があります。消さないでください。

	// バイト列を受け取るには Receive() を使用します。デバッグ用であり最終的なプログラムでは使用する必要はありません。
	m, err := yobikou.Receive()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Received: %s\n", m)

	// 配列を送るには ReceiveTable()を使用します。こちらを使用して最終的なプログラムを作成してください。
	matrix, err := yobikou.ReceiveTable()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Received: %v\n", matrix)

	// 受信する処理のみ書きましたが、ChugakuSideで行っているSend/SendTableも可能です。
}

// 中学校側はサーバーである予備校のアドレス(パートナーのアドレス)を知る必要があります。
func ChugakuSide(addr string) {
	chugaku, err := conn.NewChugakuClient(addr)
	if err != nil {
		panic(err)
	}
	defer chugaku.Close() // 消さないでください。

	// バイト列を送るには Send() を使用します。デバッグ用であり最終的なプログラムでは使用する必要はありません。
	err = chugaku.Send([]byte("ping"))
	if err != nil {
		panic(err)
	}

	// 配列を送るには SendTable()を使用します。こちらを使用して最終的なプログラムを作成してください。
	matrix := [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	err = chugaku.SendTable(matrix)
	if err != nil {
		panic(err)
	}

	// 送信する処理のみ書きましたが、YobikouSideで行っているReceive/ReceiveTableも可能です。
}
