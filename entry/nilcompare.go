package main

import (
	"fmt"
)

type Animal interface {
	m1()
}

type Duck struct {
}

func (*Duck) m1() {
}

func main() {
	var animal Animal = nil
	var duck1 *Duck = nil
	//animal = duck1 // OK
	fmt.Println(animal == nil)   //true
	fmt.Println(duck1 == nil)    //true
	fmt.Println(animal == duck1) //false
	animal = returnNilDuck()
	// 注意,指针是拥有类型的.同类型的指针的nil才是一样的
	fmt.Println(animal == nil)         //false.此时比较式左边的animal指向*Duck的nil值,而右边的nil指向左边值的类型(Animal)的nil值,两者是不一样的
	fmt.Println(animal.(*Duck) == nil) //true.此时比较式左边的值类型被转为*Duck,而右边的nil仍然指向左边值的类型(*Duck)的nil值.
	fmt.Println(animal == duck1)       //true
}

func returnNilDuck() *Duck {
	return nil
}

func init() {
	fmt.Println("====init start====")
	defer fmt.Println("====init end====")
}
