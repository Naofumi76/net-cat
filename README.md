
# Net-Cat

  

Net-Cat is a simple TCP chat server written in Go. It allows multiple clients to connect and chat with each other in real time. The server handles client connections, manages usernames, and supports basic chat commands.

  

## Features

  

-  **TCP Chat Server**: Listens for incoming TCP connections on a specified port.

-  **Username Management**: Clients can set and change their usernames.

-  **Message History**: Clients receive a history of messages upon connection.

-  **Broadcasting**: Messages are broadcasted to all connected clients.

-  **Basic Commands**: Includes commands for quitting and renaming.

   

## Build the Project:

  
```bash
go build
```
  

## Usage

  

Run the server with an optional port number. If no port is provided, the server will use the default port 8989.

  

```bash
./net-cat [port]
```
  

**Example:**

  

To run the server on port 9090, use:

  

```bash
./net-cat 9090
```
If no port number is provided:

```bash
./net-cat
```
## How to connect to the net-cat chat: 

To connect to the chat in a local environment, use the command : 
```bash 
nc localhost [port]
```

To connect to the chat with an ip, use the command :
```bash 
nc [ip] [port]
```

## Commands

  

/quit: Exit the chat. Clients will be prompted to press [ENTER] to return to the terminal.

/rename: Change the username. The client will be prompted to enter a new username.

  

## Configuration

  

The server does not require any external configuration files. It reads the port number from the command-line arguments and expects the pinguin.txt file for initial client greeting.

  

## Project Structure

  

- `main.go`: Entry point for the application, handles command-line arguments and starts the server.

- `internal/server.go`: Manages server operations, including client handling and broadcasting messages.

- `internal/client.go`: Handles individual client connections, username management, and message processing.

- `internal/util.go`: Utility functions such as converting strings to integers.

- `internal/broadcast.go`: Manages message history and broadcasting to clients.

  

## License

  

This project is licensed under the MIT License. See the LICENSE file for details.
