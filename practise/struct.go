package practise

import "fmt"

type Child struct {
	name string
	age int
}

func StructTest() {
	child := Child{"Shindou",18}
	fmt.Printf("%#v\n",child)
	fmt.Printf("poniter in main is : %p\n",&child)
	changeStruct(child)
	fmt.Printf("%#v\n",child)
	changeStructByPointer(&child)
	fmt.Printf("%#v\n",child)
}

func changeStruct(child Child) {
	defer untrace(trace("changeStruct"))
	fmt.Printf("poniter in changeStruct is : %p\n",&child)
	child.name = "newChild"
	fmt.Printf("%#v\n",child)
}

func changeStructByPointer(child *Child) {
	defer untrace(trace("changeStructByPointer"))
	fmt.Printf("poniter in changeStruct is : %p\n",&child)
	child.name = "pointer"
	fmt.Printf("%#v\n",child)
}

func trace(info string) string {
	fmt.Println("======enter",info,"======")
	return info
}

func untrace(info string) {
	fmt.Println("======leaving",info,"======")
}
