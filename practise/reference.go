package practise

import "fmt"

type Person struct {
	Name string
}

func ReferenceTest()  {
	var person = Person{"Shindou"}
	fmt.Printf("%#v\n",person)
	fmt.Printf("函数外Person的内存地址是：%p\n",&person)
	modifyPerson(person)
	fmt.Printf("%#v\n",person)
	modifyPersonByPointer(&person)
	fmt.Printf("%#v\n",person)
}

func modifyPerson(person Person)  {
	fmt.Printf("modifyPerson函数里接收到Person的内存地址是：%p\n",&person)
	person.Name = "modifyPerson"
}

func modifyPersonByPointer(person *Person)  {
	fmt.Printf("modifyPersonByPointer函数里接收到Person的内存地址是：%p\n",&person)
	person.Name = "modifyPersonByPointer"
}