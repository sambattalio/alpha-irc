package main

import (
	"fmt"
	"github.com/sambattalio/alpha-irc/client"
)

func main() {
	u := client.User{
		Server: "chat.freenode.net:6667",
		Nick: "student069client",
		User: "student069",
		Name: "sbattali",
	}

	c := client.NewClient(u)

	if c.Connect() != nil {
		fmt.Println("Error initializing client")
	}
}
