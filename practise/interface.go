package practise

import "fmt"

type Living interface {
	eat() string
}

type Human struct {
	name string
}

// If struct implements a interface with pointer receiver
// This means the eat() method is in the method set of the *Human type, but not in that of Human
func (human *Human) eat() string {
	return human.name + " is eating"
}

type Bird struct{}

func (bird Bird) eat() string {
	return "bird is eating"
}

func InterfaceTest() {
	// 如果human实现eat方法是带有指针的话,那么这里需要加上取指针符号,否则不用
	// 因为方法的归属对象有两个：值类型、引用类型
	// 所以值对象human其实是不能调用eat方法的，之所以Human{"shindou"}.eat()不报错是因为go编译时自动转化成&Human{"shindou"}.eat()
	livingEat(&Human{"shindou"})
	livingEat(Bird{})

	var man1, man2 = Human{"name1"}, &Human{"name2"}
	fmt.Println("before enter func", man1, man2)
	receiverWithoutPointer(man1)
	//receiverWithPointer(man1)// 编译错误
	//receiverWithoutPointer(man2)// 编译错误
	receiverWithPointer(man2)
	fmt.Println("after enter func", man1, man2)
}

func livingEat(living Living) {
	fmt.Println(living.eat())
}

func receiverWithPointer(human *Human) {
	human.name = "receiverWithPointer"
}

func receiverWithoutPointer(human Human) {
	human.name = "receiverWithPointer"
}
