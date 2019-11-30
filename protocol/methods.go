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
	}
}

func handlePut(cmd string, c *ClientConn) {
	fmt.Println(cmd)
	j := &Job{Tube: c.Using}
	_, err := fmt.Sscanf(cmd,"put %d %d %d %d\r\n", &j.Priority, &j.Delay, &j.TimeToRun, &j.TotalBytes)
	if err != nil {
		return
	}
	res, _, _ := c.reader.ReadLine()
	if err != nil {
		return
	}
	j.Body = string(res)
	_, err = c.SendAll([]byte("INSERTED 1 \r\n"))
	if err != nil {
		return
	}
}
