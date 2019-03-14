// 依赖github.com/tatsushid/go-fastping
package main

import (
	"fmt"
	"goland/practise"
	"net"
	"os"
	"time"
)

func main() {
	p := practise.NewPinger()
	p.Debug = true
	p.MaxRTT = time.Second
	address := []string{"www.baidu.com", "www.google.com"} //[]string{"www.baidu.com","www.google.com","sina.com"}
	for _, each := range address {
		ra, err := net.ResolveIPAddr("ip4:icmp", each)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		p.AddIPAddr(ra)
	}
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	err := p.Run()
	if err != nil {
		fmt.Println(err)
	}
}
