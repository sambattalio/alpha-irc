package client

import (
	"fmt"
	"strings"
	"errors"
	"net"
	"bufio"
)

type User struct {
	Server, Nick, User, Name string
}

type Client struct {
	conn net.Conn
}

type Message struct {
	Tags       []string
	Source     string
	Command    string
	Parameters []string
}

func Connect(c *Client, u *User) error {
	fmt.Printf("Creating connection to %v\n", u.Server)

	conn, err := net.Dial("tcp", u.Server)
	if err != nil {
		fmt.Println("Error dialing connection")
		return err
	}
	c.conn = conn

	fmt.Println("Successfully connected!")

	fmt.Println("Creating readmessage goroutine")
	go readLoop(c)

	fmt.Fprintf(c.conn, "NICK %v\r\n", u.Nick)
	fmt.Fprintf(c.conn, "USER %v - * :%v\r\n", u.User, u.Name)


	return nil
}

func readLoop(c *Client) {
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

		fmt.Printf("Tags: %q\n", parsed.Tags)
		fmt.Printf("Source: %v\n", parsed.Source)
		fmt.Printf("command: %v\n", parsed.Command)
		fmt.Printf("params: %q\n", parsed.Parameters)

		fmt.Println("Handling command")

		if (parsed.Command == "PING") {
			fmt.Println("Sending pong")
			fmt.Fprintf(c.conn, "PONG %v\r\n", parsed.Parameters[0])
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

func handle_msg(c *Client, msg string) error {
	fmt.Println(strings.Fields(msg))
	return nil
}
