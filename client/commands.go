package client

import (
	"fmt"
)

type command func(*Client, *Message)

var systemCommandList = map[string]command{
	"PING": handlePing,
}

func handlePing(c *Client, msg *Message) {
	fmt.Fprintf(c.conn, "PONG %v\r\n", msg.Parameters[0])
}

var slashCommandList = map[string]command{
	"/join": handleJoin,
}

func handleJoin(c *Client, msg *Message) {
	fmt.Fprintf(c.conn, "JOIN %v\r\n", msg.Parameters[0])
}
