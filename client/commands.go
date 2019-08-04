package client

import (
	"fmt"
	"strings"
	"github.com/awesome-gocui/gocui"
)

type command func(*Client, *Message)

// Incomming commands

var systemCommandList = map[string]command{
	"PING": handlePing,
	"NICK": handleNickChange,
	"353": handleConnectedNames,
}

func handlePing(c *Client, msg *Message) {
	fmt.Fprintf(c.conn, "PONG %v\r\n", msg.Parameters[0])
}

func handleNickChange(c *Client, msg *Message) {
	// TODO implement
	//	c.user.Nick = msg.Parameters[0]
}

func handleConnectedNames(c *Client, msg *Message) {
	// hacky need to fix
	c.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("users")
		if err != nil {
			return err
		}

		for _, user := range strings.Fields(msg.Parameters[3]) {
			fmt.Fprintf(v, "%v\r\n", user)
		}
		return nil
	})
}

// Outgoing commands

var slashCommandList = map[string]command{
	"/join": handleJoin,
	"/nick": handleSendNick,
}

func handleJoin(c *Client, msg *Message) {
	fmt.Fprintf(c.conn, "JOIN %v\r\n", msg.Parameters[0])
	c.setChannel(msg.Parameters[0])
}

func handleSendNick(c *Client, msg *Message) {
	fmt.Fprintf(c.conn, "NICK %v\r\n", msg.Parameters[0])
}
