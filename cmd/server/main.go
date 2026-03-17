package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var clients = make(map[net.Conn]string)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server started on port 9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	clients[conn] = conn.RemoteAddr().String()
	fmt.Println("Client connected:", clients[conn])

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Printf("[%s]: %s\n", clients[conn], msg)
		broadcast(conn, fmt.Sprintf("[%s]: %s\n", clients[conn], msg))
	}

	delete(clients, conn)
	fmt.Println("Client disconnected:", clients[conn])
}

func broadcast(sender net.Conn, message string) {
	for conn := range clients {
		if conn != sender {
			fmt.Fprint(conn, message)
		}
	}
}
