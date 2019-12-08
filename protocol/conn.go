package protocol

import (
	"bufio"
	"net"
)

const minBuffer = 1500

type ClientConn struct {
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	scanner  *bufio.Scanner
	Using    string
	Watching []string
	Db       Storage
}

func NewClientConn(c net.Conn) *ClientConn {
	return &ClientConn{
		conn:     c,
		reader:   bufio.NewReader(c),
		writer:   bufio.NewWriter(c),
		scanner:  bufio.NewScanner(c),
		Using:    "default",
		Watching: []string{"default"},
	}
}

func isNetTempErr(err error) bool {
	if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
		return true
	}
	return false
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

func (c *ClientConn) HandleConnection() {
	for c.scanner.Scan() {
		text := c.scanner.Text()
		handleCmd(text, c)
	}
}