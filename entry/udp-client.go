package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/net/icmp"
	"net"
	"os"
	"strings"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	serverAddr := "10.8.227.27:10006"
	conn, err := net.Dial("udp", serverAddr)
	checkError(err)
	defer conn.Close()

	ipport := conn.LocalAddr().String()
	// 经测试，该方法在Linux下能够监听到所有的icmp包，但是在windows下似乎监听不到外部请求的icmp响应包（如:通过外部ping得到的icmp响应包就监听不到）
	packetConn, err := icmp.ListenPacket("ip4:icmp", ipport[:strings.LastIndex(ipport, ":")])
	if err != nil {
		panic("ListenPacket Error:" + err.Error())
	}
	defer packetConn.Close()

	bytes := []byte("i'm from vmware")
	n, err := conn.Write(bytes)
	checkError(err)
	fmt.Println("Write:", string(bytes[:n]))
	msg := make([]byte, 100)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	n, err = conn.Read(msg)
	if err != nil {
		icmpmsg := make([]byte, 1024)
		fmt.Println("read udp error" + err.Error())
		packetConn.SetReadDeadline(time.Now().Add(time.Second))
		// 注意：ICMP响应报文通常情况下会被拦截
		i, addr, err := packetConn.ReadFrom(icmpmsg)
		if err != nil {
			fmt.Printf("read icmp message error:%v\n", err.Error())
		} else {
			fmt.Println("received icmp message")
			fmt.Println(addr.String())
			fmt.Print(hex.Dump(icmpmsg[:i]))
			//注意：这里似乎会接受到所有的ICMP报文，并不只是刚才发送的UDP请求对应的ICMP报文
			//TODO 校验接收到的icmp是否匹配上次请求
			hanlderICMPPacked(icmpmsg[:i])
		}
	} else {
		fmt.Println("Response:", string(msg[:n]))
	}
}

// 处理icmp报文
func hanlderICMPPacked(bytes []byte) {
	itype := bytes[0]
	icode := bytes[1]
	ichecksum := int(bytes[2])*256 + int(bytes[3])
	fmt.Printf("icmp type:%v\n", itype)
	fmt.Printf("icmp code:%v\n", icode)
	fmt.Printf("icmp checksum:%v\n", ichecksum)
}
