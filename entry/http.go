package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PollResult struct {
	totalMillisecond int64
	dnsMillisecond int64
	connectMillisecond int64
	sendMillisecond int64
	responseMillisecond int64
	redirectMillisecond int64
	downloadMillisecond int64
	fileSize int64
	downloadSpeed float64 // kb/s
}

func main()  {
	pollURL("http://localhost:8090")
}

func pollURL(urlstring string) (result PollResult) {
	startTime := time.Now()
	_,result.dnsMillisecond = dnsparse(urlstring)
	endTime := time.Now()
	result.totalMillisecond = endTime.Sub(startTime).Nanoseconds()/1000000
	httpgettest(urlstring)
	return
}

// dnsparse解析指定的url,并返回解析后的地址和解析消耗的时间
func dnsparse(u string) (addr []string,costTime int64) {
	parseResult, e := url.Parse(u)
	HandlerErr(e)
	host := parseResult.Host
	if i := strings.Index(parseResult.Host,":");i > 0 {
		host = host[:i]
	}
	dnsStartTime := time.Now()
	addr, e = net.LookupHost(host)
	dnsEndTime := time.Now()
	HandlerErr(e)
	costTime = dnsEndTime.Sub(dnsStartTime).Nanoseconds() / 1000000
	return
}

func httpgettest(url string) {
	resp, err := http.Get(url)
	HandlerErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	HandlerErr(err)
	fmt.Println(string(body))
	fmt.Println("StatusCode:", resp.StatusCode)
	fmt.Println(resp.Request.URL)//这个URL可能是重定向后的URL
}

func noRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func httpGetWithNoRedirect(url string) {
	client := &http.Client{
		CheckRedirect: noRedirect,
	}
	resp, err := client.Get(url)
	defer resp.Body.Close()
	HandlerErr(err)
	fmt.Println("StatusCode:", resp.StatusCode)
	fmt.Println(resp.Request.URL)
	fmt.Println("Redirect Location:",resp.Header.Get("Location"))
}

func HandlerErr(err error){
	if err != nil {
		panic(err.Error())
	}
}