package main

import (
	"bufio"
	"net"
	"os"
	"fmt"
	"chat/common"
)


func AttachClient(term *common.Terminal, client *Client) {
	go func() {
		for {
			term.Out <- <-client.In
		}
	}()

	go func() {
		for {
			client.Out <- <-term.In
		}
	}()
}

type Client struct {
	common.Terminal
	Name string
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		*common.NewTerminal(conn, conn),
		"Unnamed",
	}
}

func main() {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}

	client := NewClient(conn)
	fmt.Println("Enter your name:")
	name, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Name = name
	terminal := common.NewTerminal(os.Stdin, os.Stdout)
	AttachClient(terminal, client)
	client.Writer.WriteString(client.Name)
	client.Writer.Flush()
	client.Listen()
	go terminal.Write()
	terminal.Read()
}