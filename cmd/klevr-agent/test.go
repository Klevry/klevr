package main

import (
	"fmt"
)



func main() {
	go Ping()

	fmt.Scanln() // main 함수가 종료되지 않도록 대기
}