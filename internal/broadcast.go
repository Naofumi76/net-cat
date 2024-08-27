package internal

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Sender    string
	Text      string
	Timestamp time.Time
}

var (
	Messages   []*Message
	MessagesMu sync.Mutex
)

func AddMessage(sender *Client, message string) {
	if sender == nil {
		fmt.Println("Error: Attempted to add message with nil sender")
		return
	}

	MessagesMu.Lock()
	defer MessagesMu.Unlock()

	msg := &Message{
		Sender:    sender.Name,
		Text:      message,
		Timestamp: time.Now(),
	}

	Messages = append(Messages, msg)
}

func GetMessages() []*Message {
	MessagesMu.Lock()
	defer MessagesMu.Unlock()

	// Return a copy of Messages to avoid race conditions
	messagesCopy := make([]*Message, len(Messages))
	copy(messagesCopy, Messages)

	return messagesCopy
}

func BroadcastMessage(sender *Client, message string) {
	ClientsMu.Lock()
	defer ClientsMu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	message = fmt.Sprintf("[%s] %s", timestamp, message)

	// Store message if sender is not nil
	if sender != nil {
		AddMessage(sender, message)
	}

	// Broadcast message to all clients
	for _, client := range Clients {
		if client != nil && sender != nil && client != sender {
			client.Conn.Write([]byte(message + "\n"))
		}
	}
}
