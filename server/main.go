package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var clients = make(map[net.Conn]*Client)
var mutex = sync.Mutex{}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Chat server started on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		mutex.Lock()
		clients[conn] = newClient("")
		mutex.Unlock()
		go handleClient(conn)
	}

}

func handleClient(conn net.Conn) {
	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		mutex.Lock()
		err = handleMessage(msg, conn)
		mutex.Unlock()
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}
		mutex.Lock()
		clients[conn].addMessage(msg)
		mutex.Unlock()

		mutex.Lock()
		idMessage := fmt.Sprintf("%s: %s", clients[conn].name, msg)
		fmt.Println(clients[conn].name + "hey")
		for client := range clients {
			if client != conn {
				client.Write([]byte(idMessage))
			}
		}
		mutex.Unlock()
	}
}

type Client struct {
	name     string
	messages []string
	isConn   bool
}

func newClient(n string) *Client {
	return &Client{name: n, messages: make([]string, 10), isConn: true}
}

func (c *Client) addMessage(message string) {
	c.messages = append(c.messages, message)
}

func handleMessage(message string, conn net.Conn) error {
	client := clients[conn]

	if client.name == "" && len(message) > 0 && rune(message[0]) != '/' {
		resp := []byte("Need to use comnmand /setuser <USERNAME> before you can send a message\n")
		conn.Write(resp)
		return fmt.Errorf("User account not set")
	}

	if len(message) > 0 && rune(message[0]) == '/' {
		m := strings.Split(message, " ")
		if m[0] == "/setuser" {
			//There is new line at the end of m[1]
			client.name = m[1][:len(m[1])-1]
		}
	}
	return nil
}
