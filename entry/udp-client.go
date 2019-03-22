package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	serverAddr := "10.8.156.137:10006"
	conn, err := net.Dial("udp", serverAddr)
	checkError(err)

	defer conn.Close()

	bytes := []byte("i'm from vmware")
	n, err := conn.Write(bytes)
	checkError(err)
	fmt.Println("Write:", string(bytes[:n]))
	msg := make([]byte, 100)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	n, err = conn.Read(msg)
	checkError(err)
	fmt.Println("Response:", string(msg[:n]))
}
