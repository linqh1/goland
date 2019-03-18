// 依赖 miekg.dns库
package main

import (
	"context"
	"fmt"
	"github.com/bogdanovich/dns_resolver"
	"log"
	"net"
)

var customResolver = net.Resolver{
	PreferGo: true,
	Dial:     GoogleDNSDialer,
}

func GoogleDNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
	d := net.Dialer{}
	fmt.Println("trigger custom resolver...")
	return d.DialContext(ctx, "udp", "8.8.8.8:53")
}

func main() {
	//www.myhostname.linqh是一个不存在的域名,在hosts文件中被指向127.0.0.1
	//blog.csdn.net在hosts文件中被指向127.0.0.1
	localDnsResolverTest("www.myhostname.linqh")
	localDnsResolverTest("blog.csdn.net")

	resolver := dns_resolver.New([]string{"8.8.8.8"})
	resolver.RetryTimes = 5
	googleDnsResolverTest(resolver, "www.myhostname.linqh")
	googleDnsResolverTest(resolver, "blog.csdn.net")
}

func localDnsResolverTest(addr string) {
	hosts, err := net.LookupIP(addr)
	if err != nil {
		log.Printf("local dns lookup host error:" + err.Error())
	}
	log.Printf("local dns lookup host:%v result:%v\n", addr, hosts)
}

func googleDnsResolverTest(resovler *dns_resolver.DnsResolver, addr string) {
	hosts, err := resovler.LookupHost(addr)
	if err != nil {
		log.Printf("google dns lookup host %v error:%v\n", addr, err.Error())
	}
	log.Printf("google dns lookup host:%v result:%v\n", addr, hosts)
}
