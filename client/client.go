package client

import (
	"fmt"
	"strings"
	"errors"
	"net"
	"bufio"
	"github.com/awesome-gocui/gocui"
)

type User struct {
	Server, Nick, User, Name string
}

type Client struct {
	conn net.Conn
	user User
	channel string
	Gui  *gocui.Gui
}

type Message struct {
	Tags       []string
	Source     string
	Command    string
	Parameters []string
}

func (c *Client) Connect(u *User) error {
	fmt.Printf("Creating connection to %v\n", u.Server)

	conn, err := net.Dial("tcp", u.Server)
	if err != nil {
		fmt.Println("Error dialing connection")
		return err
	}
	c.conn = conn

	go c.readLoop()

	fmt.Fprintf(c.conn, "NICK %v\r\n", u.Nick)
	fmt.Fprintf(c.conn, "USER %v - * :%v\r\n", u.User, u.Name)

	return nil
}

func (c *Client) GetInput(_ *gocui.Gui, v *gocui.View) error {
	input := v.ViewBuffer()

	parsed, err := parseMessage(input);
	if err != nil {
		fmt.Println(err)
		return err
	}

	if handler, ok := slashCommandList[parsed.Command]; ok {
		handler(c, parsed)
	} else {
		fmt.Fprintf(c.conn, "%v\r\n", input)
	}

	v.Clear()
	err = v.SetCursor(0,0)
	return err
}

func (c *Client) setChannel(name string) {
	writeToScreen(c, name, "channels")
	c.channel = name
}

func (c *Client) readLoop() {
	reader := bufio.NewReader(c.conn)

	// loop while connected
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		parsed, err := parseMessage(msg);
		if err != nil {
			fmt.Println(err)
			break
		}

		writeToScreen(c, strings.Join(parsed.Parameters, " "), "stream")

		if handler, ok := systemCommandList[parsed.Command]; ok {
			handler(c, parsed)
		}
	}
}

func parseMessage(msg string) (*Message, error) {
	if (len(msg) == 0) {
		return &Message{}, errors.New("String Length 0!")
	}

	split := strings.Fields(msg)

	parsed := &Message{}

	if (strings.HasPrefix(split[0], "@")) {
		tags, err := parseTags(split[0])
		if err != nil {
			return &Message{}, err
		}
		parsed.Tags = tags
		split = split[1:]
	}

	if (strings.HasPrefix(split[0], ":")) {
		parsed.Source = split[0][1:]
		split = split[1:]
	}

	parsed.Command = split[0]
	split = split[1:]

	for i, item := range split {
		if (strings.HasPrefix(item, ":")) {
			parsed.Parameters = append(parsed.Parameters,
				            strings.Join(split[i:], " ")[1:])
			break
		}
		parsed.Parameters = append(parsed.Parameters, item)
	}

	return parsed, nil
}

func parseTags(s string) ([]string, error) {
	// TODO
	return make([]string, 0), nil
}

func handleMsg(c *Client, msg string) error {
	fmt.Println(strings.Fields(msg))
	return nil
}

func writeToScreen(c *Client, msg string, view string) {
	c.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(view)
		if err != nil {
			return err
		}
		fmt.Fprintln(v, msg)
		//fmt.Fprintln(v, msg.Source + ":" +
		//	     strings.Join(msg.Parameters[1:], " "))
		return nil
	})
}
