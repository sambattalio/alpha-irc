package client

import (
	"fmt"
	"net"
	"bufio"
)

type User struct {
	Server, Nick, User, Name string
}

type Client struct {
	conn net.Conn
	err chan string
	write chan string
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
	go readMessage(c)

	fmt.Fprintf(c.conn, "NICK %v\r\n", u.Nick)
	fmt.Fprintf(c.conn, "USER %v - * :%v\r\n", u.User, u.Name)


	return nil
}

func readMessage(c *Client) {
	reader := bufio.NewReader(c.conn)

	// loop while connected
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(msg)
	}
}
