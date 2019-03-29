// 解析外部字符串数组到相应flag上
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	os.Args = []string{"curl.exe", "-X", "POST", "-d", "123", "www.baidu.com"}
	var method, data string
	flag.StringVar(&method, "X", "GET", "http method")
	flag.StringVar(&data, "d", "", "request body")
	flag.Parse()
	fmt.Printf("method:%v data:%v flag.args:%v\n", method, data, flag.Args())

	set := flag.NewFlagSet("new", flag.ExitOnError)
	set.StringVar(&method, "X", "GET", "http method")
	set.StringVar(&data, "d", "", "request body")
	set.Parse([]string{"-X", "PUT", "-d", "zeihahahaha", "www.google.com"})
	fmt.Printf("method:%v data:%v flag.args:%v\n", method, data, set.Args())
}
