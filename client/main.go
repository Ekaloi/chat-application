package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to chat server")
	go readMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text() + "\n"
		fmt.Println()
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}

func readMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Server closed connection")
			os.Exit(0)
		}
		fmt.Print(msg)
	}
}
