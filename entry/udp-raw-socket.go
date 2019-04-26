// 该方法似乎无法接受到ICMP报文
package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/net/icmp"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) //实现真正的随机数
	listenPort := 61500
	destinationip := "10.8.227.27"
	packet := UDPHeader{
		Source:      uint16(listenPort), // Random ephemeral port
		Destination: 10006,              //这里的端口可能需要修改一下
		DataOffset:  0,                  // 4 bits
		Checksum:    0,                  // Kernel will set this if it's 0
		Data:        []byte("i'm from vmware"),
	}

	conn, err := net.Dial("ip4:udp", destinationip)
	if err != nil {
		log.Fatalf("Dial: %s\n", err)
	}
	defer conn.Close()
	//源端口如果没有设置监听，如果server回复UDP响应报文的话，client因为没有监听该端口，还会回复一个ICMP unreachable的报文
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   parseToIp4(conn.LocalAddr().String()),
		Port: listenPort,
	})
	if err != nil {
		panic("ListenUDP Error" + err.Error())
	}
	defer udpConn.Close()

	// 监听icmp报文,在获取udp报文超时的时候,接收icmp报文
	fmt.Printf("listen %v icmp\n", conn.LocalAddr().String())
	packetConn, err := icmp.ListenPacket("ip4:icmp", conn.LocalAddr().String())
	if err != nil {
		panic("ListenPacket Error:" + err.Error())
	}
	defer packetConn.Close()

	packet.DataOffset = uint16(8 + len(packet.Data))
	data := packet.Marshal()
	packet.Checksum = UDPChecksum(data, parseToIp4Byte(conn.LocalAddr().String()), parseToIp4Byte(destinationip))
	data = packet.Marshal()
	_, err = conn.Write(data)
	if err != nil {
		log.Fatal("write error", err)
	}
	bytes := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	//这里有时候会返回timeout，有时会立即返回connection refuse，不知道为什么。
	//upd-client的就不会 = =
	readnum, err := conn.Read(bytes)
	if err != nil {
		fmt.Printf("read udp message error:%v\n", err)
		complete := make(chan int)
		go receieveICMPPacked(packetConn, complete)
		<-complete
		return
	}
	fmt.Print(hex.Dump(bytes[:readnum]))
	ipheaderLenth := bytes[:readnum][0] & 0x0f //第1个字节低4位表示ip头长度，单位是32bit，即4字节
	udpdata := bytes[:readnum][ipheaderLenth*4:]
	fmt.Println("receieve udp")
	fmt.Println(hex.Dump(udpdata))
	fmt.Println("udp data:")
	fmt.Println(string(udpdata[8:]))
}

func receieveICMPPacked(packetConn *icmp.PacketConn, complete chan int) {
	//如果读取出错（超时）,那么就尝试获取以下icmp报文
	//可以用go程来做?
	// 注意：ICMP响应报文通常情况下会被拦截
	errorTimeLimit := 10
	errorCnt := 0
	for {
		icmpmsg := make([]byte, 1024)
		packetConn.SetReadDeadline(time.Now().Add(time.Second * 2))
		i, addr, err := packetConn.ReadFrom(icmpmsg)
		if err != nil {
			fmt.Printf("read icmp messge error:%v\n", err.Error())
			if strings.Index(err.Error(), "i/o timeout") >= 0 {
				complete <- 1
				break
			} else {
				//防止死循环
				errorCnt = errorCnt + 1
				if errorCnt > errorTimeLimit {
					complete <- 1
					break
				}
			}
		} else {
			fmt.Println("received icmp message")
			fmt.Println(addr.String())
			fmt.Print(hex.Dump(icmpmsg[:i]))
			icmpbytes := icmpmsg[:i]
			itype := icmpbytes[0]
			icode := icmpbytes[1]
			ichecksum := int(icmpbytes[2])*256 + int(icmpbytes[3])
			fmt.Printf("icmp type:%v\n", itype)
			fmt.Printf("icmp code:%v\n", icode)
			fmt.Printf("icmp checksum:%v\n", ichecksum)
			// TODO 这里要做校验:如果接收到的icmp报文就是响应我们刚才的请求的话就结束循环
		}
	}
}

type UDPHeader struct {
	Source      uint16
	Destination uint16
	DataOffset  uint16
	Checksum    uint16
	Data        []byte
}

func (tcp *UDPHeader) Marshal() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, tcp.Source)
	binary.Write(buf, binary.BigEndian, tcp.Destination)
	binary.Write(buf, binary.BigEndian, tcp.DataOffset)
	binary.Write(buf, binary.BigEndian, tcp.Checksum)
	binary.Write(buf, binary.BigEndian, tcp.Data)
	return buf.Bytes()
}

func UDPChecksum(msg []byte, srcip, dstip [4]byte) uint16 {
	pseudoHeader := []byte{
		srcip[0], srcip[1], srcip[2], srcip[3],
		dstip[0], dstip[1], dstip[2], dstip[3],
		0,  // zero
		17, // protocol number (17 == UDP)
		0, byte(len(msg)),
	}
	sumThis := make([]byte, 0, len(pseudoHeader)+len(msg))
	sumThis = append(sumThis, pseudoHeader...)
	sumThis = append(sumThis, msg...)
	var sum uint32
	var nextWord uint32
	for n := 0; n < len(sumThis)-1; n += 2 {
		nextWord = uint32(sumThis[n])<<8 | uint32(sumThis[n+1])
		sum += uint32(nextWord)
	}
	if len(sumThis)%2 != 0 {
		lastByte := uint16(sumThis[len(sumThis)-1]) << 8
		sum += uint32(lastByte)
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	return uint16(^sum)
}

func parseToIp4Byte(ip string) [4]byte {
	addr := net.ParseIP(ip)
	if addr == nil {
		panic("can not recongnize ip" + ip)
	}
	addr = addr.To4()
	if addr == nil {
		panic("can only support ipv4" + ip)
	}
	return [4]byte{addr[0], addr[1], addr[2], addr[3]}
}

func parseToIp4(ip string) net.IP {
	addr := net.ParseIP(ip)
	if addr == nil {
		panic("can not recongnize ip" + ip)
	}
	addr = addr.To4()
	if addr == nil {
		panic("can only support ipv4" + ip)
	}
	return addr
}
