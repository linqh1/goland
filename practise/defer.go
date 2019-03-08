package practise

import (
	"fmt"
)

// defer 貌似和java的finally一样
// 在defer所在的func返回之前,会执行defer 语句
// 最后声明的defer将最先被执行（延迟函数以LIFO顺序执行）
func Deftest() {
	fmt.Println("====deftest====")
	fmt.Println("first line")
	defer fmt.Println("defer line")
	defer fmt.Println("defer line2")
	fmt.Println("last line")
	defer fmt.Println("defer line3")

	//defer 语句的参数在defer语句声明时就会计算保留下来
	//而不是等到实际执行时才计算
	for i := 0; i < 5; i++ {
		defer fmt.Println("defer in for %d ", i)
	}
	//如果是对象引用的话,比如defer print(a)
	//defer 时,a.age=10
	//return前,a.age=20,那么实际执行refer的时候会打印出什么？？？
}

// 1
func F() (result int) {
	defer func() {
		result++
	}()
	return 0
}

// 5
func F1() (r int) {
	t := 5
	defer func() {
		t = t + 5
	}()
	return t
}

// 1
func F2() (r int) {
	defer func(r int) {
		r = r + 5 //这里的r并不是返回值的r了,所以对它的操作不会体现在返回值中
	}(r)
	return 1
}

// 5
func F3() (r int) {
	defer func(t int) {
		r = t + 5
	}(r)
	return 1
}

// 6
func F4() (r int) {
	defer func(t int) {
		r += t + 5
	}(r)
	return 1
}

func DeferOperateAObject() map[string] string {
	var result map[string] string
	result = make(map[string] string)
	result["key1"] = "value1"
	defer fmt.Printf("the value of key1 is [%s] when defer is called\n",result["key1"])
	defer fmt.Printf("the value is [%s] when defer is called\n",result)
	result["key1"] = "value2"
	return result
}
