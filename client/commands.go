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
	"PRIVMSG": handlePRIVMSG,
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

func handlePRIVMSG(c *Client, msg *Message) {
	writeToView(c, msg)
}

// Outgoing commands

var slashCommandList = map[string]command{
	"/join": handleJoin,
	"/nick": handleSendNick,
	"/msg": handleMSG,
}

func handleJoin(c *Client, msg *Message) {
	fmt.Fprintf(c.conn, "JOIN %v\r\n", msg.Parameters[0])
	c.setChannel(msg.Parameters[0])
}

func handleSendNick(c *Client, msg *Message) {
	fmt.Fprintf(c.conn, "NICK %v\r\n", msg.Parameters[0])
}

func handleMSG(c *Client, msg *Message) {
	if (len(msg.Parameters) < 2) {
		return
	}
	c.setChannel(msg.Parameters[0])
	fmt.Fprintf(c.conn, "PRIVMSG %v :%v\r\n", msg.Parameters[0], strings.Join(msg.Parameters[1:], " "))
	c.writeInputToScreen(strings.Join(msg.Parameters[1:], " "))
}
