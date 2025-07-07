package main

import "fmt"

func main() {
	e := Employee{
		Person: Person{
			Name: "张三",
			Age:  30,
		}, EmployeeID: 123}
	e.PrintInfo()
}

type Person struct {
	Name string
	Age  uint8
}

type Employee struct {
	Person
	EmployeeID uint64
}

func (e Employee) PrintInfo() {
	fmt.Printf("Name:%s,Age:%d,EmployeeID:%d", e.Name, e.Age, e.EmployeeID)
}
