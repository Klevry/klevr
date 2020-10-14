package main

import "fmt"

type Test struct {
	name string
	a    []Number
}

type Number struct {
	num int
}

func main() {
	a := &Test{}
	n := &Number{}
	n1 := &Number{}

	n.num = 0
	n1.num = 1

	b := make([]Number, 0)
	b = append(b, *n)
	b = append(b, *n1)

	a.name = "cjh"
	a.a = b

	fmt.Println(a)
}
