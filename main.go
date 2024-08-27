package main

import (
	"fmt"
	"netcat/internal"
	"os"
)

func main() {
	// Default port number
	defaultPort := 8989

	// Check if a port number is provided as a command-line argument
	var port int
	if len(os.Args) == 1 {
		// No port provided, use default
		port = defaultPort
	} else if len(os.Args) == 2 {
		// Parse the provided port number
		var err error
		port, err = internal.StringToInt(os.Args[1])
		if err != nil || port <= 0 || port > 65535 {
			fmt.Println("Invalid port number. Please provide a valid port number between 1 and 65535.")
			os.Exit(1)
		}
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}

	internal.Start(port)
}
