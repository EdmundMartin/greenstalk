package main

import (
	"fmt"
	"github.com/EdmundMartin/greenstalk/protocol"
	"github.com/EdmundMartin/greenstalk/protocol/postgres"
	"github.com/EdmundMartin/greenstalk/stateManager"
	"net"
)

func main() {
	l, err := net.Listen("tcp4", ":11300")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	statusChanges := make(chan *stateManager.HeapValue)
	dbConn := postgres.NewPGConn(`localhost`, `postgres`, `edmund`, `beanstalk`, 5432)
	postgres.CreateTable(dbConn)
	dbConn.Updates = statusChanges
	go stateManager.ManageState(statusChanges)
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
