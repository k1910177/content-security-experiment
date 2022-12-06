package main

import (
	consec "contentssecurity"
	"time"
)

func main() {
	// ここからプログラム記述開始

	// 中学
	// server_addr := "0.0.0.0"
	// ↑ハードコードしたくない場合
	// 方法 1. コマンドライン引数で受け取る: flagパッケージ
	// 方法 2. 環境変数で受け取る: godotenvモジュール & os.Getenv関数
	// consec.ChugakuSide(server_addr)

	// 予備校
	// consec.YobikouSide()

	// ↓1人で検証用
	demo()
}

// パートナーとではなく、一人で検証を行う際は、以下のようにします。
func demo() {
	yc := make(chan interface{}) // goroutine終了を知らせるためのチャンネル
	// goroutineを使用し、マルチスレッド処理としています。
	go func() {
		// lib.goを読んでください。
		consec.YobikouSide()
		yc <- struct{}{}
	}()

	// 予備校がサーバーを建てるのを少し待ちます
	time.Sleep(100 * time.Millisecond)

	cc := make(chan interface{})
	go func() {
		consec.ChugakuSide("0.0.0.0")
		cc <- struct{}{}
	}()

	<-yc
	<-cc
}
