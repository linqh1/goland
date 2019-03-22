// 只在linux下能编译成功
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
	packet := TCPHeader{
		Source:      0xaa47, // Random ephemeral port
		Destination: 8090,   //这里的端口可能需要修改一下
		SeqNum:      rand.Uint32(),
		AckNum:      0,
		DataOffset:  5,      // 4 bits
		Reserved:    0,      // 3 bits
		ECN:         0,      // 3 bits
		Ctrl:        2,      // 6 bits (000010, SYN bit set)
		Window:      0xaaaa, // size of your receive window
		Checksum:    0,      // Kernel will set this if it's 0
		Urgent:      0,
		Options:     []TCPOption{},
	}

	conn, err := net.Dial("ip4:tcp", "10.8.156.137")
	if err != nil {
		log.Fatalf("Dial: %s\n", err)
	}
	defer conn.Close()

	data := packet.Marshal()
	packet.Checksum = Csum(data, parseToIp4(conn.LocalAddr().String()), [4]byte{10, 8, 156, 137})
	data = packet.Marshal()
	_, err = conn.Write(data) //在这一步过程中，如果受到目标机器返回的SYN-ACK包之后，会自动发送一个RST包给目标机器
	if err != nil {
		log.Fatal("write error", err)
	}
	bytes := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	readnum, err := conn.Read(bytes) //这里似乎不会收到10.8.156.137这个的地址发来的数据?
	if err != nil {
		fmt.Printf("read error:%v\n", err)
		return
	}
	fmt.Print(hex.Dump(bytes[:readnum]))
	ipheaderLenth := bytes[:readnum][0] & 0x0f                //第1个字节低4位表示ip头长度，单位是32bit，即4字节
	header := NewTCPHeader(bytes[:readnum][ipheaderLenth*4:]) //收到的TCP包,包中的SYN和ACK标志都应该为1,并且ACKNum应该为发送包seq+1
	fmt.Printf("send seq:%v, receive ack:%v\n", packet.SeqNum, header.AckNum)
	fmt.Printf("receive tcp packet,syn:%v,ack:%v,rst:%v \n", (header.Ctrl&0x02)>>1, (header.Ctrl&0x10)>>4, (header.Ctrl&0x04)>>2)

}

const (
	FIN = 1  // 00 0001
	SYN = 2  // 00 0010
	RST = 4  // 00 0100
	PSH = 8  // 00 1000
	ACK = 16 // 01 0000
	URG = 32 // 10 0000
)

type TCPHeader struct {
	Source      uint16
	Destination uint16
	SeqNum      uint32
	AckNum      uint32
	DataOffset  uint8 // 4 bits
	Reserved    uint8 // 3 bits
	ECN         uint8 // 3 bits
	Ctrl        uint8 // 6 bits
	Window      uint16
	Checksum    uint16 // Kernel will set this if it's 0
	Urgent      uint16
	Options     []TCPOption
}

type TCPOption struct {
	Kind   uint8
	Length uint8
	Data   []byte
}

// Parse packet into TCPHeader structure
func NewTCPHeader(data []byte) *TCPHeader {
	var tcp TCPHeader
	r := bytes.NewReader(data)
	binary.Read(r, binary.BigEndian, &tcp.Source)
	binary.Read(r, binary.BigEndian, &tcp.Destination)
	binary.Read(r, binary.BigEndian, &tcp.SeqNum)
	binary.Read(r, binary.BigEndian, &tcp.AckNum)

	var mix uint16
	binary.Read(r, binary.BigEndian, &mix)
	tcp.DataOffset = byte(mix >> 12)  // top 4 bits
	tcp.Reserved = byte(mix >> 9 & 7) // 3 bits
	tcp.ECN = byte(mix >> 6 & 7)      // 3 bits
	tcp.Ctrl = byte(mix & 0x3f)       // bottom 6 bits

	binary.Read(r, binary.BigEndian, &tcp.Window)
	binary.Read(r, binary.BigEndian, &tcp.Checksum)
	binary.Read(r, binary.BigEndian, &tcp.Urgent)

	return &tcp
}

func (tcp *TCPHeader) HasFlag(flagBit byte) bool {
	return tcp.Ctrl&flagBit != 0
}

func (tcp *TCPHeader) Marshal() []byte {

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, tcp.Source)
	binary.Write(buf, binary.BigEndian, tcp.Destination)
	binary.Write(buf, binary.BigEndian, tcp.SeqNum)
	binary.Write(buf, binary.BigEndian, tcp.AckNum)

	var mix uint16
	mix = uint16(tcp.DataOffset)<<12 | // top 4 bits
		uint16(tcp.Reserved)<<9 | // 3 bits
		uint16(tcp.ECN)<<6 | // 3 bits
		uint16(tcp.Ctrl) // bottom 6 bits
	binary.Write(buf, binary.BigEndian, mix)

	binary.Write(buf, binary.BigEndian, tcp.Window)
	binary.Write(buf, binary.BigEndian, tcp.Checksum)
	binary.Write(buf, binary.BigEndian, tcp.Urgent)

	for _, option := range tcp.Options {
		binary.Write(buf, binary.BigEndian, option.Kind)
		if option.Length > 1 {
			binary.Write(buf, binary.BigEndian, option.Length)
			binary.Write(buf, binary.BigEndian, option.Data)
		}
	}

	out := buf.Bytes()

	// Pad to min tcp header size, which is 20 bytes (5 32-bit words)
	pad := 20 - len(out)
	for i := 0; i < pad; i++ {
		out = append(out, 0)
	}

	return out
}

// TCP Checksum
func Csum(data []byte, srcip, dstip [4]byte) uint16 {

	//伪头部,计算校验和的时候使用的,结构如下
	//源IP(4个字节) 目标ip(4个字节) 0(1个字节) 协议类型(1个字节) TCP数据长度(2个字节)
	pseudoHeader := []byte{
		srcip[0], srcip[1], srcip[2], srcip[3],
		dstip[0], dstip[1], dstip[2], dstip[3],
		0,                  // zero
		6,                  // protocol number (6 == TCP)
		0, byte(len(data)), // TCP length (16 bits), not inc pseudo header 这里的TCP数据长度前8位为什么设置为0???
	}

	sumThis := make([]byte, 0, len(pseudoHeader)+len(data))
	sumThis = append(sumThis, pseudoHeader...)
	sumThis = append(sumThis, data...)
	//fmt.Printf("% x\n", sumThis)

	lenSumThis := len(sumThis)
	var nextWord uint16
	var sum uint32
	for i := 0; i+1 < lenSumThis; i += 2 {
		nextWord = uint16(sumThis[i])<<8 | uint16(sumThis[i+1])
		sum += uint32(nextWord)
	}
	if lenSumThis%2 != 0 {
		fmt.Println("Odd byte")
		sum += uint32(sumThis[len(sumThis)-1])
	}

	// Add back any carry, and any carry from adding the carry
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	fmt.Printf("checksum:%v\n", uint16(^sum))
	// Bitwise complement
	return uint16(^sum)
}

func parseToIp4(ip string) [4]byte {
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
