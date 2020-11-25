package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// Maximum request is 1024 bytes.
	request := make([]byte, 1024)
	_, err := bufio.NewReader(conn).Read(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Print the request
	fmt.Println(request)
	// Send ACK response
	_, err = bufio.NewWriter(conn).WriteString("Ack")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	// By default, the gateway listens to port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go handleConnection(conn)
	}
}
