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

	channels map[string]bool
}

type Message struct {
	Tags       []string
	Source     string
	Command    string
	Parameters []string
}

func NewClient(u User, gui *gocui.Gui) Client {
	return Client{
		user: u,
		channel: "stream",
		Gui: gui,
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.user.Server)
	if err != nil {
		return err
	}
	c.conn = conn

	go c.readLoop()

	fmt.Fprintf(c.conn, "NICK %v\r\n", c.user.Nick)
	fmt.Fprintf(c.conn, "USER %v - * :%v\r\n", c.user.User, c.user.Name)

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
		writeToScreen(c, input, c.channel)
		fmt.Fprintf(c.conn, "PRIVMSG %v :%v\r\n", c.channel, input)
	}

	v.Clear()
	err = v.SetCursor(0,0)
	return err
}

func (c *Client) setChannel(name string) error {
	// Check if PRIVMSG view exists
	if _, ok := c.channels[name]; !ok {
		maxX, maxY := c.Gui.Size()
		if v, err := c.Gui.SetView(name, maxX / 6 + 1, 0, maxX - 1, maxY - 4, 0); err != nil {
			if !gocui.IsUnknownView(err) {
				return err
			}

			v.Wrap = true
			v.Autoscroll = true
		}
	}

	c.channel = name
	return nil
}

func (c *Client) readLoop() {
	reader := bufio.NewReader(c.conn)

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

		writeToScreen(c, strings.Join(parsed.Parameters[1:], " "), c.channel)

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
		return nil
	})
}
