package main

import (
	"fmt"
	"github.com/EdmundMartin/greenstalk/protocol"
	"github.com/EdmundMartin/greenstalk/protocol/postgres"
	"net"
)

func main() {
	l, err := net.Listen("tcp4", ":11300")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	dbConn := postgres.NewPGConn(`localhost`, `postgres`, `edmund`, `beanstalk`, 5432)
	postgres.CreateTable(dbConn)
	fmt.Println("Serving on port 11300")
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		conn := protocol.NewClientConn(c)
		conn.Db = dbConn
		//go conn.handleConnection()
		go conn.HandleConnection()
	}
}
