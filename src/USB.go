package main

import (
	"fmt"
	"runtime"
)

// Mac,Windows,Linuxそれぞれ別の形式でUSBを検出する模様なので、それぞれのOSに対応する関数を作成する
// モジュール系統での実装を行いたいが、コンパイル時に恐らくエラーが発生してしまうため
// OSに実装されているコマンドで実装する
func DetectionUSB() {
	if runtime.GOOS == "windows" {
		//
	} else if runtime.GOOS == "darwin" {
		//diskutil list | grep external
		//上記のコードを実行するとUSBデバイスを検出できる
		fmt.Println("USBデバイスを検出中...")
	} else if runtime.GOOS == "linux" {
		//lsblk | grep disk
		//上記のコードを実行するとUSBデバイスを検出できる
		fmt.Println("USBデバイスを検出中...")
	}
}
