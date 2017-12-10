package main

import (
	"net"
	"fmt"
	"strings"
	"chat/common"
)

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

type Server struct {
	clients []*Client
	joins   chan net.Conn
	in      chan string
}

func NewServer() *Server {
	return &Server{
		make([]*Client, 0),
		make(chan net.Conn),
		make(chan string),
	}
}

func (server *Server) Listen() {
	go func() {
		for {
			select {
			case data := <-server.in:
				server.Broadcast(data)
			case conn := <-server.joins:
				server.Join(conn)
			}
		}
	}()
}

func (server *Server) Broadcast(data string) {
	for _, client := range server.clients {
		client.Out <- data
	}
}

func (server *Server) Join(conn net.Conn) {
	client := NewClient(conn)
	client.Listen()
	client.Name = strings.TrimSuffix(<-client.In, "\n")
	server.clients = append(server.clients, client)
	go func() {
		for {
			server.in <- client.Name + ": " + <-client.In
		}
	}()
}

func main() {
	server := NewServer()
	server.Listen()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Join")
		server.joins <- conn
	}
}
