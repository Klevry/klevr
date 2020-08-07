package main

import (
"fmt"

"github.com/mackerelio/go-osstat/memory"
)

func main() {
memory, err := memory.Get()
if err != nil{
fmt.Println(err)
}

fmt.Println(int(memory.Total/1024/1024))
}