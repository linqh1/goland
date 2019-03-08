package practise

import (
	"fmt"
	"time"
)

func ChannelTest() {
	produceAndConsumer()
}

// 无缓冲chan
func blockTest() {
	c := make(chan int) //这时候没有指定容量,即为无缓冲chan,每次发送/读取都会阻塞goroutine
	go func() {
		fmt.Printf("%v: getting c...run at go func.\n",time.Now())
		get := <- c
		fmt.Printf("%v: geted c[%v]! run at go func.\n",time.Now(),get)
	}()
	fmt.Printf("%v: set c after 5 seconds.run at main func.\n",time.Now())
	time.Sleep(time.Second * 2)
	c <- 5
	fmt.Printf("%v: set c! run at main func.\n",time.Now())
	//如果这里不阻塞main程的话，这里一执行完，go程序就退出了。
	// 此时go协程可能走完了，也可能还没走完
	time.Sleep(time.Second * 2)
}

// 无缓冲chan
func blockTest2() {
	c := make(chan int)
	gocomplete := make(chan int)
	go func() {
		fmt.Printf("%v: getting c...run at go func.\n",time.Now())
		get := <- c
		fmt.Printf("%v: geted c[%v]! run at go func.\n",time.Now(),get)
		fmt.Printf("%v: sleep 5 second! run at go func.\n",time.Now())
		time.Sleep(time.Second * 5)
		gocomplete <- 1
	}()
	fmt.Printf("%v: set c after 5 seconds.run at main func.\n",time.Now())
	time.Sleep(time.Second * 2)
	c <- 5
	fmt.Printf("%v: set c! run at main func.\n",time.Now())
	<-gocomplete //可以通过一个chan来判断go协程是否完成
	fmt.Printf("%v: goroutine run complete and exit at main func.\n",time.Now())
}

func produceAndConsumer() {
	ch := make(chan int,10)
	go produce(ch)
	go consumer(ch)
	time.Sleep(1 * time.Minute)
}

func produce(p chan<- int) {
	i := 0
	for {
		i++
		p <- i
		fmt.Printf("%v : send:%v\n", time.Now(),i)
		time.Sleep(time.Millisecond * 500) //防止发送太快，每0.5秒发送一次
		if i == 50 {
			time.Sleep(time.Second * 10) //观察接受者是否阻塞
		}
	}
}
func consumer(c <-chan int) {
	fmt.Println("consumer sleep 5 seconds")
	time.Sleep(5 * time.Second)
	i := 0
	for {
		i++
		v := <-c
		fmt.Printf("%v : receive:%v\n", time.Now(),v)
	}
}
