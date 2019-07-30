package main

import (
	"fmt"
	"github.com/sambattalio/alpha-irc/client"
)

func main() {
	c := &client.Client{}

	u := &client.User{
		Server: "chat.freenode.net:6667",
		Nick: "student069client",
		User: "student069",
		Name: "sbattali",
	}

	if client.Connect(c, u) != nil {
		fmt.Println("Error initializing client")
	}
	for {

	}
}
