package irc

import (
	"fmt"
	"strings"

	"net"
	"bufio"
)

type User struct {
        Server, Nick, User, Name string
}

type Client struct {
        conn net.Conn
        user User

        channel string
        channels map[string]bool

	MsgStream chan *Message
	InpStream chan *Message

        users []string
}

type Message struct {
        Tags       []string
        Source     string
        Command    string
        Parameters []string
}

// NewClient initializes and returns a pointer to a new Client struct
func NewClient(u User) *Client {
	c := Client{
		user: u,
		channels: make(map[string]bool),
		users: make([]string, 0),
	}
	return &c
}

// Connect creates the server conneciton and attempts to register with
// the set Nick / User in the User Struct
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.user.Server)
        if err != nil {
                return err
        }
        c.conn = conn

        go c.readLoop()

        fmt.Fprintf(c.conn, "NICK %v\r\n", c.user.Nick)
        fmt.Fprintf(c.conn, "USER %v - * :%v\r\n", c.user.User, c.user.Name)

	/*
	for name, ok := range channels {
		// Connect to preset channels TODO
	}
	*/

        return nil
}

// Used to add a channel to the map and
func (c *Client) AddChannel(name string) error {
	c.channels[name] = true
}

// Used to remove channel from map
func (c *Client) RemoveChannel(name string) error {
	c.channels[name] = false
}

// This is the "main" loop where it constantly reads incomming messages
// and delivers them to the proper handlers
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

		messageStream <- parsed
	}
}

// Used to parse the message into proper formatting
// Decided against regex b/c it can have so many different
// setups ¯\_(ツ)_/¯
func parseMessage(msg string) (*Message, error) {
	parsed := &Message{}
	if (len(msg) == 0) {
		return parsed, errors.New("Len0")
	}

	split := strings.Fields(msg)

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

// TODO: handle tags
func parseTags(s string) ([]string, error) {
	return make([]string, 0), nil
}

