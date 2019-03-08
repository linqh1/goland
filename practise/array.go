package practise

import "fmt"

func ArrayTest() {
	var array = [...]float64{1.1,2.2,3.3}
	fmt.Printf("进函数SumFor3Float64(a *[3]float64) (sum float64)前参数的内存地址是：%p\n",&array)
	fmt.Printf("进函数SumFor3Float64(a *[3]float64) (sum float64)前的值时：%v\n",array)
	SumFor3Float64(&array)
	fmt.Printf("出函数SumFor3Float64(a *[3]float64) (sum float64)后的值时：%v\n",array)
	fmt.Printf("进函数Double3ArryValue(a [3]float64)前的值时：%v\n",array)
	Double3ArryValue(array)
	fmt.Printf("出函数Double3ArryValue(a [3]float64)后的值时：%v\n",array)
	//var array2 = make([] float64,3)
	//array2[0] = 100
	//array2[1] = 300
	//array2[2] = 300
	var array2 = [] float64 {100,300,300}
	fmt.Printf("进函数DoubleArryValue(a []float64)前参数的内存地址是：%p\n",&array2)
	fmt.Printf("进函数Double3ArryValue(a [3]float64)前的值时：%v\n",array2)
	DoubleArryValue(array2)
	fmt.Printf("出函数Double3ArryValue(a [3]float64)后的值时：%v\n",array2)
}

// The size of an array is part of its type.
// The types [10]int and [20]int are distinct.
func SumFor3Float64(a *[3]float64) (sum float64) {
	fmt.Printf("SumFor3Float64(a *[3]float64) (sum float64)函数里接收到参数的内存地址是：%p\n",&a)
	for _, v := range *a {
		sum += v
	}
	v := *a
	v[0] = 9.9 //这里改变值不会影响到返回后的值，为什么？
	v[1] = 8.8
	v[2] = 7.7
	return
}

// [i] int 值传递
// [i] int 这是数组,数组时值传递
func Double3ArryValue(a [3]float64) {
	fmt.Printf("Double3ArryValue(a [3]float64)函数里接收到参数的内存地址是：%p\n",&a)
	for index,value := range a {
		a[index] = value * 2
	}
}

// [] int 引用传递
// [] int 这是切片,切片是引用传递
func DoubleArryValue(a []float64) {
	fmt.Printf("DoubleArryValue(a []float64)函数里接收到参数的内存地址是：%p\n",&a)
	for index,value := range a {
		a[index] = value * 2
	}
}

func init() {
	fmt.Println("init in array.go")
}
