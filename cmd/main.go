package main

import (
	"fmt"
	"github.com/EdmundMartin/greenstalk/protocol"
	"net"
)

func main() {
	l, err := net.Listen("tcp4", ":11300")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		conn := protocol.NewClientConn(c)
		//go conn.handleConnection()
		go conn.HandleConnection()
	}
}
