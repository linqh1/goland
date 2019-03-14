// 依赖github.com/tatsushid/go-fastping
package main

import (
	"fmt"
	"goland/practise"
	"golang.org/x/net/icmp"
	"net"
	"time"
)

func main() {
	p := practise.NewPinger()
	p.MaxRTT = time.Second
	address := []string{"10.8.227.27"} //[]string{"www.baidu.com","www.google.com","sina.com"}
	for _, each := range address {
		//ra, err := net.ResolveIPAddr("ip4:icmp", each)
		//if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		//}
		p.AddIPAddr(&net.IPAddr{
			IP: net.ParseIP(each),
		})
	}
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	p.OnICMPResponse = func(message icmp.Message) {
		fmt.Printf("receive icmp response:%#v\n", message)
	}
	p.OnTimeout = func(_ *icmp.PacketConn) {
		fmt.Printf("receive icmp response timeout\n")
	}
	err := p.Run()
	if err != nil {
		fmt.Println(err)
	}
}
