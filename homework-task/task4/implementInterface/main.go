package main

import "fmt"

func main() {
	var _ IInterface = (*Implement)(nil)
}

type IInterface interface {
	OK1()
	OK2()
}

type Implement struct {
}

func (i *Implement) OK1() {
	fmt.Println("OK1")
}
func (i *Implement) OK2() {
	fmt.Println("OK1")
}
