package practise

import (
	"fmt"
	"reflect"
	"strings"
)

type MyType1 string
type MyType2 string

func ReflectTest() {

	var str = "abc"
	fmt.Println("send a string")
	accept(str)
	fmt.Println("send a MyType1")
	accept(MyType1(str))
	fmt.Println("send a MyType2")
	accept(MyType2(str))
}

func accept(a interface{}) {
	//fmt.Printf("type:%v\n",reflect.TypeOf(a))
	switch a.(type) {
	case string:
		str := a.(string)
		fmt.Println("i'm a string", str)
	case MyType1:
		type1 := a.(MyType1)
		type1.type1Method()
	default:
		fmt.Println("can not recongnize type", reflect.TypeOf(a))
	}
	fmt.Println()
}

func (t1 MyType1) cal() string {
	return strings.ToUpper(string(t1))
}

func (t1 MyType1) type1Method() {
	fmt.Println("i'm a method of MyType1")
}

func (t2 MyType2) cal() string {
	return strings.ToLower(string(t2))
}

func (t2 MyType2) type2Method() {
	fmt.Println("i'm a method of MyType2")
}
