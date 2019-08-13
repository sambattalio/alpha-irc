package client

import (
	"fmt"
	"strings"
	"unicode"
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
	gui  *gocui.Gui

	channel string
	channels map[string]bool

	users []string
}

type Message struct {
	Tags       []string
	Source     string
	Command    string
	Parameters []string
}

func NewClient(u User) *Client {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		fmt.Println(err)
	}
	g.Cursor = true

	c := Client{
		user: u,
		gui: g,
		channels: make(map[string]bool),
		users: make([]string, 0),
	}
	g.SetManager(&c)
	c.setKeybindings()
	return &c
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

	c.startGui()

	return nil
}

func (c *Client) GetInput(_ *gocui.Gui, v *gocui.View) error {
	input := v.ViewBuffer()

	parsed, err := parseMessage(input);
	if err != nil {
		// Garbage input don't crash just _gracefully_ handle
		return nil
	}

	parsed.Source = c.user.Nick

	if handler, ok := slashCommandList[parsed.Command]; ok {
		handler(c, parsed)
	} else {
		c.writeInputToScreen(input)
		fmt.Fprintf(c.conn, "PRIVMSG %v :%v\r\n", c.channel, input)
	}

	v.Clear()
	err = v.SetCursor(0,0)
	return err
}

func (c *Client) setChannel(name string) error {
	var err error = nil
	if !c.isChannel(name) {
		err = c.addChannel(name)
	}
	c.channel = name
	c.gui.SetViewOnTop(name)
	if err = c.resetUsersTab(); err != nil {
		return err
	}
	c.reloadUsers()
	return err
}

func (c *Client) deleteChannel(_ *gocui.Gui, v *gocui.View) error {
	channel, err := lineOfText(v);
	if err != nil {
		return err
	}

	c.channels[channel] = false
	err = c.gui.DeleteView(channel)
	if err != nil {
		return err
	}
	c.updateChannelsList()

	channelsList, err := c.gui.View("channels")
	if err != nil {
		return err
	}
	c.setChannelView(nil, channelsList)

	return nil
}

func (c *Client) setChannelView(_ *gocui.Gui, v *gocui.View) error {
	var channel string
        var err error
	if channel, err = lineOfText(v); err != nil {
		return err
	}
	c.channel = channel
	c.gui.SetCurrentView("input")
        c.gui.SetViewOnTop(channel)
	if err = c.resetUsersTab(); err != nil {
		return err
	}
	c.reloadUsers()
        return nil
}

func (c *Client) autoFill(_ *gocui.Gui, v *gocui.View) error {
	// Try to fill with users if existi
	v, err := c.gui.View("input")
	if err != nil {
		return err
	}

	_, curY := v.Cursor()
	input := v.ViewBuffer()
	words := strings.Fields(input)
	if len(words) == 0 {
		return nil
	}
	substring := words[len(words) - 1]

	for _, name := range c.users {
		name = strings.TrimLeftFunc(name, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		if strings.HasPrefix(name, substring) {
			toWrite := strings.Split(name, substring)[1]
			v.Clear()
			fmt.Fprintf(v, "%s%s", input, toWrite)
			v.SetCursor(len(input) + len(toWrite), curY)
			break
		}
	}

	return nil
}

func lineOfText(v *gocui.View) (string, error) {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	return line, err
}

func (c *Client) reloadUsers() {
	fmt.Fprintf(c.conn, "NAMES %v\r\n", c.channel)
}

func (c *Client) resetUsersTab() error {
	// reset users view
	v, err := c.gui.View("users")
	if err != nil {
		return err
	}
	v.Clear()
	return nil
}

func (c *Client) isChannel(name string) bool {
	val, ok := c.channels[name]
	return ok && val
}

func (c *Client) addChannel(name string) error {
	maxX, maxY := c.gui.Size()
	if v, err := c.gui.SetView(name, maxX / 6 + 1, 0, maxX - 1, maxY - 4, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Wrap = true
		v.Autoscroll = true
	}
	c.gui.SetViewOnBottom(name)
	c.channels[name] = true
	c.appendChannelsList(name)
	return nil
}

func (c *Client) appendChannelsList(channel string) {
	c.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("channels")
		if err != nil {
			return err
		}
		fmt.Fprintf(v, "%v\n", channel)
		return nil
	})
}

func (c *Client) updateChannelsList() {
	c.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("channels")
		v.Clear()
		if err != nil {
			return err
		}
		for channel, value := range c.channels {
			if value {
				fmt.Fprintf(v, "%v\n", channel)
			}
		}
		return nil
	})
}

func (c *Client) readLoop() {
	reader := bufio.NewReader(c.conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		parsed, err := parseMessage(msg)
		if err != nil {
			fmt.Println(err)
			break
		}

		if handler, ok := systemCommandList[parsed.Command]; ok {
			handler(c, parsed)
		}
	}
}

func parseMessage(msg string) (*Message, error) {
	// imagine using regex
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
		parsed.Source = strings.Split(split[0][1:], "!")[0]
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
	// not super needed... yet
	return make([]string, 0), nil
}

func handleMsg(c *Client, msg string) error {
	fmt.Println(strings.Fields(msg))
	return nil
}

func writeToView(c *Client, msg *Message) error{
	var view string
	if (msg.Parameters[0] == c.user.Nick) {
		view = msg.Source
	} else {
		view = msg.Parameters[0]
		if view == "*" {
			view = msg.Source
		}
	}

	if !c.isChannel(view) {
		err := c.addChannel(view)
		if err != nil {
			return err
		}
	}
	/*if (!strings.HasPrefix(msg.Source, "#") || strings.Contains(strings.Join(msg.Parameters[1:], " "), c.user.Nick)) {
		fmt.Print("\a")
	}*/

	c.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(view)
		if err != nil {
			return err
		}

		fmt.Fprintf(v, "<%v> %v\n", msg.Source, strings.Join(msg.Parameters[1:], " "))
		return nil
	})

	return nil
}

func updateGuiFunc(gui *gocui.Gui, updater func(g *gocui.Gui) error) {
	gui.Update(func(g *gocui.Gui) error {
		return updater(g)
	})
}

func (c *Client) writeInputToScreen(msg string) {
	updateGuiFunc(c.gui, func(g *gocui.Gui) error {
		v, err := g.View(c.channel)
		if err != nil {
			return err
		}
		fmt.Fprintf(v, "<%v> %v\n", c.user.Nick, msg)
		return nil
	})
}
