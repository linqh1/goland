package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"
)

func main() {
	request, err := http.NewRequest("GET", "http://www.baidu.com", nil) //14.215.177.39
	trace := &httptrace.ClientTrace{
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo.Addrs)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("target ip:%+v\n", connInfo.Conn.RemoteAddr().String())
		},
	}
	request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == "www.baidu.com:80" {
				addr = "14.215.177.39:80" //这里可以指定ip
			}
			return (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext(ctx, network, addr)
		},
	}
	client := &http.Client{
		Transport: tr,
	}
	response, err := client.Do(request)
	if err != nil {
		panic("response err" + err.Error())
	}
	body := response.Body
	defer func() {
		body.Close()
	}()
	bytes, readerr := ioutil.ReadAll(response.Body)
	if readerr != nil {
		fmt.Println("read response error!" + readerr.Error())
	}
	fmt.Println("request success!", len(bytes) > 0)
	fmt.Println(string(bytes))
}
