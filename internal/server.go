package internal

import (
	"fmt"
	"net"
)

func Start(port int) {
	InitUsernames()
	InitMessages() // Initialize messages here

	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server is running and listening on %s...\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		client := &Client{
			Conn: conn,
		}

		// Handle each client in a separate goroutine
		go HandleClient(client)
	}
}

func InitMessages() {
	Messages = make([]*Message, 0)
}
