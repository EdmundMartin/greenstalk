package protocol

import (
	"fmt"
	"strings"
)

func handleCmd(cmd string, c *ClientConn) {
	verb := strings.Split(cmd, " ")[0]
	switch  verb {
	case "put":
		handlePut(cmd, c)
	case "reserve":
		handleReserve(cmd, c)
	}
}

func handleReserve(cmd string, c *ClientConn) {
	fmt.Println(cmd)
	j, err := c.Db.Reserve(c.Watching)
	fmt.Println(j)
	if err != nil {
		fmt.Println(err)
		c.writer.Write([]byte("TIMED_OUT\r\n"))
	}
	c.SendAll([]byte(fmt.Sprintf("RESERVED %d %d\r\n", j.ID, j.TotalBytes)))
	c.SendAll([]byte(fmt.Sprintf("%s\r\n", j.Body)))
}

func handlePut(cmd string, c *ClientConn) {
	fmt.Println(cmd)
	j := &Job{Tube: c.Using}
	_, err := fmt.Sscanf(cmd,"put %d %d %d %d\r\n", &j.Priority, &j.Delay, &j.TimeToRun, &j.TotalBytes)
	if err != nil {
		return
	}
	for c.scanner.Scan() {
		j.Body = c.scanner.Text()
		res, err := c.Db.Save(j)
		if err != nil {
			c.SendAll([]byte("DRAINING\r\n"))
			return
		}
		_, err = c.SendAll([]byte(fmt.Sprintf("INSERTED %d \r\n", res)))
		if err != nil {
			return
		}
		return
	}
}
