package main

import (
	"bufio"
	"fmt"
	"net"
)

const minBuffer = 1500

type ClientConn struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewClientConn(conn net.Conn) *ClientConn {
	return &ClientConn{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func isNetTempErr(err error) bool {
	if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
		return true
	}
	return false
}

func (c *ClientConn) SendAll(msg []byte) (int, error) {
	written := 0
	forWrite := msg
	var n int
	var err error
	for written < len(msg) {
		forBuff := len(forWrite) >= minBuffer
		if forBuff {
			n, err = sendAllBuffer(c, forWrite)
			if err != nil && !isNetTempErr(err) {
				return written, err
			}
		} else {
			n, err = sendAllNoBuffer(c, forWrite)
			if err != nil && !isNetTempErr(err) {
				return written, err
			}
		}
		written += n
		forWrite = forWrite[n:]
	}
	return written, nil
}

func sendAllNoBuffer(c *ClientConn, msg []byte) (int, error) {
	n, err := c.conn.Write(msg)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func sendAllBuffer(c *ClientConn, msg []byte) (int, error) {
	n, err := c.writer.Write(msg)
	if err != nil {
		return n, err
	}
	err = c.writer.Flush()
	if err != nil {
		return n, err
	}
	return n, nil
}

type Job struct {
	Priority   int
	Delay      int
	TimeToRun  int
	TotalBytes int
	Body string
}

func handlePut(cmd string, cc *ClientConn) error {
	var pri, delay, ttr, tBytes int
	_, err := fmt.Sscanf(cmd,"put %d %d %d %d\r\n", &pri, &delay, &ttr, &tBytes)
	if err != nil {
		return err
	}
	j := &Job{
		Priority:   pri,
		Delay:      delay,
		TimeToRun:  ttr,
		TotalBytes: tBytes,
	}
	res, _, _ := cc.reader.ReadLine()
	j.Body = string(res)
	fmt.Println(j)
	_, err = cc.SendAll([]byte("INSERTED 1 \r\n"))
	if err != nil {
		return err
	}
	return nil
}

func (cc *ClientConn) handleConnection() {
	for {
		res, _, _ := cc.reader.ReadLine()
		if len(res) > 0 {
		handlePut(string(res), cc)
		}
	}
}


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
		conn := NewClientConn(c)
		//go conn.handleConnection()
		go conn.handleConnection()
	}
}
