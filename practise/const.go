package practise

import (
	"fmt"
	"os"
)

type ByteSize float64

// iota在const关键字出现时将被重置为0(const内部的第一行之前)
// const中每新增一行常量声明将使iota计数一次(iota可理解为const语句块中的行索引)
const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota) //iota=1
	MB //iota=2
	GB //iota=3
	TB //iota=4
	PB //iota=5
	EB //iota=6
	ZB //iota=7
	YB //iota=8
)

func ConstTest(){
	fmt.Printf("%T %v\n",KB,KB)
	fmt.Println(os.Getenv("GOPATH"))
}

func init(){
	fmt.Println("init in const.go")
}