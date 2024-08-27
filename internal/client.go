package internal

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

type Client struct {
	Conn net.Conn
	Name string
}

var (
	Clients       []*Client
	ClientsMu     sync.Mutex
	ActiveClients int
	Usernames     map[string]bool
	UsernamesMu   sync.Mutex
)

// Create a map to stock the usernames
func InitUsernames() {
	Usernames = make(map[string]bool)
}

func HandleClient(client *Client) {
	// Possible nil client used after, just a warning, not an error
	defer client.Conn.Close()

	if client == nil {
		fmt.Println("Nil client passed to HandleClient")
		return
	}

	scanner := bufio.NewScanner(client.Conn)

	pinguin, _ := os.ReadFile("pinguin.txt")
	client.Conn.Write([]byte(pinguin))
	for {
		// Ask client for name
		client.Conn.Write([]byte("\n[ENTER YOUR NAME]: "))
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil && err != io.EOF {
				fmt.Printf("Error reading username: %v\n", err)
			}
			return
		}
		name := scanner.Text()
		UsernamesMu.Lock()
		if _, exists := Usernames[name]; !exists && name != "" {
			Usernames[name] = true
			UsernamesMu.Unlock()
			client.Name = name
			break
		}
		UsernamesMu.Unlock()
		client.Conn.Write([]byte("Username invalid or already used.\n"))
	}

	// Send previous messages to the client upon connection
	if len(Messages) != 0 {
		client.Conn.Write([]byte("History ------------------------------------------------\n"))
		for _, msg := range Messages {
			client.Conn.Write([]byte(msg.Text + "\n"))
		}
		client.Conn.Write([]byte("--------------------------------------------------------\n"))
	}

	ClientsMu.Lock()
	Clients = append(Clients, client)
	ActiveClients++
	ClientsMu.Unlock()

	// Notify all clients about the new connection
	BroadcastMessage(client, fmt.Sprintf("[SERVER INFO]: %s has joined the chat.", client.Name))

	fmt.Printf("%s connected\n", client.Name)

	for {
		message, err := bufio.NewReader(client.Conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading from %s: %v\n", client.Name, err)
			break
		}

		message = strings.TrimSpace(message)

		if strings.HasPrefix(message, "/quit") {
			client.Conn.Write([]byte("Press [ENTER] to return to basic terminal window."))
			break
		}

		if strings.HasPrefix(message, "/rename") {
			for {
				client.Conn.Write([]byte("Enter a new username: "))
				newScanner := bufio.NewScanner(client.Conn)
				if newScanner.Scan() {
					newName := strings.TrimSpace(newScanner.Text())
					if newName == "" {
						client.Conn.Write([]byte("New username cannot be empty.\n"))
						fmt.Printf("Handled Error: %s tried to change username to an invalid username\n", client.Name)
						continue
					}

					UsernamesMu.Lock()
					if _, exists := Usernames[newName]; !exists {
						delete(Usernames, client.Name)
						Usernames[newName] = true
						UsernamesMu.Unlock()

						oldName := client.Name
						client.Name = newName
						BroadcastMessage(client, fmt.Sprintf("[SERVER INFO]: %s has changed username to '%s'", oldName, client.Name))
						fmt.Printf("%s changed username to '%s'\n", oldName, client.Name)
						break
					} else {
						UsernamesMu.Unlock()
						client.Conn.Write([]byte("Username already taken. Please try a different one.\n"))
						fmt.Printf("Handled Error: %s tried to change username to '%s' (Username already taken)\n", client.Name, newName)
					}
				} else {
					if err := newScanner.Err(); err != nil && err != io.EOF {
						fmt.Printf("Error reading new username: %v\n", err)
					}
					break
				}
			}
			continue
		}

		// Broadcast message to all other clients
		if message != "" {
			BroadcastMessage(client, fmt.Sprintf("%s: %s", client.Name, message))
		}
	}

	// Notify all clients about the disconnection
	ClientsMu.Lock()
	Clients = removeClient(Clients, client)
	ActiveClients--
	ClientsMu.Unlock()

	BroadcastMessage(client, fmt.Sprintf("[SERVER INFO]: %s has left the chat.", client.Name))

	fmt.Printf("%s disconnected\n", client.Name)

	// Remove username from the map
	UsernamesMu.Lock()
	delete(Usernames, client.Name)
	UsernamesMu.Unlock()

	// If no more active clients, shut down the server
	if ActiveClients == 0 {
		fmt.Println("No more clients connected. Shutting down server.")
		os.Exit(0)
	}
}

func removeClient(clients []*Client, client *Client) []*Client {
	for i, c := range clients {
		if c == client {
			return append(clients[:i], clients[i+1:]...)
		}
	}
	return clients
}
