package protocol

import (
	"fmt"
	"strings"
)

func handleCmd(cmd string, c *ClientConn) {
	verb := strings.Split(cmd, " ")[0]
	switch  verb {
	case "watch":
		handleWatch(cmd, c)
	case "ignore":
		handleIgnore(cmd, c)
	case "put":
		handlePut(cmd, c)
	case "reserve":
		handleReserve(cmd, c)
	case "delete":
		handleDelete(cmd, c)
	case "use":
		handleUse(cmd, c)
	default:
		fmt.Println("unsupported command")
	}
}

func removeIdx(watching []string, index int) []string {
	return append(watching[:index], watching[index+1:]...)
}

func handleIgnore(cmd string, c *ClientConn) {
	fmt.Println(cmd)
	var toIgnore string
	fmt.Sscanf(cmd, "watch %s\r\n", &toIgnore)
	if len(c.Watching) == 1 {
		c.SendAll([]byte("NOT_IGNORED\r\n"))
		return
	}
	foundTube := c.findTube(toIgnore)
	if foundTube != -1 {
		c.Watching = removeIdx(c.Watching, foundTube)
		c.SendAll([]byte(fmt.Sprintf("WATCHING %d\r\n", len(c.Watching))))
		return
	}
	c.SendAll([]byte(fmt.Sprintf("WATCHING %d\r\n", len(c.Watching))))
}

func handleWatch(cmd string, c *ClientConn) {
	fmt.Println(cmd)
	var toWatch string
	fmt.Sscanf(cmd, "watch %s\r\n", &toWatch)
	c.insertWatching(toWatch)
	c.SendAll([]byte(fmt.Sprintf("WATCHING %d\r\n", len(c.Watching))))
}

func handleReserve(cmd string, c *ClientConn) {
	fmt.Println(cmd)
	j, err := c.Db.Reserve(c.Watching)
	if err != nil {
		c.SendAll([]byte("TIMED_OUT\r\n"))
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

func handleDelete(cmd string, c *ClientConn) {
	fmt.Println(cmd)
	var jobID int
	_, err := fmt.Sscanf(cmd, "delete %d\r\n", &jobID)
	if err != nil {
		return
	}
	found := c.Db.Delete(&Job{ID: jobID})
	if found {
		c.SendAll([]byte("DELETED\r\n"))
		return
	}
	c.SendAll([]byte("NOT_FOUND\r\n"))
}

func handleUse(cmd string, c *ClientConn) {
	fmt.Println(cmd)
	var toUse string
	_, err := fmt.Sscanf(cmd, "use %s\r\n", &toUse)
	if err != nil {
		return
	}
	c.Using = toUse
	msg := fmt.Sprintf("USING %s\r\n", toUse)
	c.SendAll([]byte(msg))
}