package practise

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func FileReadTest() {
	f,err := os.Open("D:\\workspace\\IDEA\\qtl-ms-portal\\src\\main\\resources\\application.properties")
	if err != nil {
		panic(err.Error())
	}
	defer func(){
		fmt.Printf("file[%v] close",f.Name())
		f.Close()
	}()

	reader := bufio.NewReader(f)
	fmt.Printf("file[%v] start reading...\n",f.Name())
	for {
		line, err := reader.ReadBytes('\n')
		fmt.Printf(string(line))
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				fmt.Printf("file[%v] end read\n",f.Name())
				break
			}
			fmt.Printf("file[%v] read error!\n",err)
			break
		}
	}
}
