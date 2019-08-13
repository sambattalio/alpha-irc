package client

import (
	"fmt"
	"strings"
	"github.com/awesome-gocui/gocui"
)


type command func(*Client, *Message) (error)
// Incomming commands

var systemCommandList = map[string]command{
	"PING": handlePing,
	"NICK": handleNickChange,
	"353": handleConnectedNames,
	"372": writeToView,
	"PRIVMSG": writeToView,
	"NOTICE": writeToView,
}

func handlePing(c *Client, msg *Message) error {
	fmt.Fprintf(c.conn, "PONG %v\r\n", msg.Parameters[0])
	return nil
}

func handleNickChange(c *Client, msg *Message) error {
	// TODO implement
	//	c.user.Nick = msg.Parameters[0]
	return nil
}

func handleConnectedNames(c *Client, msg *Message) error {
	// hacky need to fix
	c.users = nil
	c.users = make([]string, 0)
	c.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("users")
		if err != nil {
			return err
		}

		for _, user := range strings.Fields(msg.Parameters[3]) {
			fmt.Fprintf(v, "%v\r\n", user)
			c.users = append(c.users, user)
		}
		return nil
	})
	return nil
}

func handleMOTD(c *Client, msg *Message) error {
	writeToView(c, msg)
	return nil
}

func handlePRIVMSG(c *Client, msg *Message) error {
	writeToView(c, msg)
	return nil
}

func handleNOTICE(c *Client, msg *Message) error {
	writeToView(c, msg)
	return nil
}

// Outgoing commands

var slashCommandList = map[string]command{
	"/join": handleJoin,
	"/nick": handleSendNick,
	"/msg": handleMSG,
}

func handleJoin(c *Client, msg *Message) error {
	fmt.Fprintf(c.conn, "JOIN %v\r\n", msg.Parameters[0])
	c.setChannel(msg.Parameters[0])
	return nil
}

func handleSendNick(c *Client, msg *Message) error {
	fmt.Fprintf(c.conn, "NICK %v\r\n", msg.Parameters[0])
	return nil
}

func handleMSG(c *Client, msg *Message) error {
	if (len(msg.Parameters) < 2) {
		return nil
	}
	c.setChannel(msg.Parameters[0])
	fmt.Fprintf(c.conn, "PRIVMSG %v :%v\r\n", msg.Parameters[0], strings.Join(msg.Parameters[1:], " "))
	c.writeInputToScreen(strings.Join(msg.Parameters[1:], " "))
	return nil
}
