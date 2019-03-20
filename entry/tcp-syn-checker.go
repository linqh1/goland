// 依赖github.com/tevino/tcp-shaker
// 扫描端口时，不进行三次握手，而是只发送一次SYN包，如果收到服务端相应的SYN-ACK包就默认端口可用
package main

import (
	"context"
	"fmt"
	"github.com/tevino/tcp-shaker"
	"time"
)

func main() {
	c := tcp.NewChecker()

	ctx, stopChecker := context.WithCancel(context.Background())
	defer stopChecker()
	go func() {
		if err := c.CheckingLoop(ctx); err != nil {
			fmt.Println("checking loop stopped due to fatal error: ", err)
		}
	}()

	<-c.WaitReady()

	timeout := time.Second * 1
	err := c.CheckAddr("10.8.156.137:8090", timeout)
	switch err {
	case tcp.ErrTimeout:
		fmt.Println("Connect to Google timed out")
	case nil:
		fmt.Println("Connect to Google succeeded")
	default:
		fmt.Println("Error occurred while connecting: ", err)
	}
}
