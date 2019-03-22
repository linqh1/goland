// 该方法似乎无法接受到ICMP报文
package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) //实现真正的随机数
	listenPort := 0xaa47
	packet := UDPHeader{
		Source:      uint16(listenPort), // Random ephemeral port
		Destination: 10006,              //这里的端口可能需要修改一下
		DataOffset:  0,                  // 4 bits
		Checksum:    0,                  // Kernel will set this if it's 0
		Data:        []byte("i'm from vmware"),
	}

	conn, err := net.Dial("ip4:udp", "10.8.156.137")
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

	packet.DataOffset = uint16(8 + len(packet.Data))
	data := packet.Marshal()
	packet.Checksum = UDPChecksum(data, parseToIp4Byte(conn.LocalAddr().String()), [4]byte{10, 8, 156, 137})
	data = packet.Marshal()
	_, err = conn.Write(data)
	if err != nil {
		log.Fatal("write error", err)
	}
	bytes := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	readnum, err := conn.Read(bytes)
	if err != nil {
		fmt.Printf("read error:%v\n", err)
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
	for n := 1; n < len(sumThis)-1; n += 2 {
		sum += uint32(sumThis[n])*256 + uint32(sumThis[n+1])
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
