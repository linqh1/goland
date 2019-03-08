package practise

import (
	"fmt"
	"log"
	"time"
)

//一个协程出错的话，会影响到所有协程
func UnRecoverableTest(){
	for i:=0;i<5;i=i+1 {
		go func(k int) {
			if k == 4 {
				fmt.Println(1/(k - 4))
			}else{
				fmt.Println(k)
			}
		}(i)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(10 * time.Second)
	fmt.Println("complete")
}

//可以采用recover()来捕捉处理error
func RecoverableTest(){
	for i:=0;i<6;i=i+1 {
		go func(k int) {
			//recover必须定义在panic之前的defer语句中
			defer func(){
				if err := recover();err != nil{
					log.Println("work failed:", err)
				}
			}()
			if k == 4 {
				fmt.Println(1/(k - 4))
			}else{
				fmt.Println(k)
			}
			// recover如果放在这里的话,会没有效果
		}(i)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(3 * time.Second)
	fmt.Println("complete")
}
