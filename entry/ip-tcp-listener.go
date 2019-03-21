package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	//conn, e := net.DialTimeout("tcp", "10.8.156.137:8090", time.Second*500)
	//if e != nil {
	//	log.Fatal("dial error",e)
	//}
	//defer conn.Close()
	//log.Print("dial success!\n")
	netaddr, err := net.ResolveIPAddr("ip4", "192.168.51.128")
	if err != nil {
		log.Fatal("ResolveIPAddr error", err)
	}
	conn, err := net.ListenIP("ip4:tcp", netaddr)
	if err != nil {
		log.Fatal("ListenIP error", err)
	}
	for {
		buf := make([]byte, 1480)
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("ReadFrom error", err)
		} else {
			tcpHeaderMain := NewTCPHeaderMain(buf[:n])
			fmt.Println(hex.Dump(buf[:n]))
			fmt.Println(string(buf[:n][tcpHeaderMain.DataOffset*4:]))
		}
	}
}

type TCPHeaderMain struct {
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
	Options     []TCPOptionMain
}

type TCPOptionMain struct {
	Kind   uint8
	Length uint8
	Data   []byte
}

// Parse packet into TCPHeader structure
func NewTCPHeaderMain(data []byte) *TCPHeaderMain {
	var tcp TCPHeaderMain
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
