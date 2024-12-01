package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Client struct {
	ID          int
	Name        string
	Conn        net.Conn
	IsConnected bool
}
	
var (
	clients           = make(map[net.Conn]Client) // Store Client structs
	mu                sync.Mutex
	onlineClientCount = 0
	serverCap         = 10
	log               = ""
)

func main() {

	args := os.Args[1:]
	port := "8989"

	if len(args) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	if len(args) == 1 {
		port = args[0]
	}
	// Listen for incoming connections on port 8080
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listen.Close()

	fmt.Println("Listening on the port :" + port)

	clientID := 0 // Counter for assigning unique IDs to clients

	for {

		// Accept a new connection
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Increment client ID
		clientID++

		// Handle the new connection
		go handleClient(conn, clientID)
	}
}

func DeleteClient(conn net.Conn) {
	// Remove client from the map and close the connection
	mu.Lock()
	delete(clients, conn)
	mu.Unlock()
	conn.Close()
}

// handleClient manages client communication
func handleClient(conn net.Conn, id int) {
	defer DeleteClient(conn)

	// Send welcome message and prompt for name
	onlineClientCount++

	// Limit server cap
	if onlineClientCount > serverCap {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error blocking connection server is in cap:", err)
		}
		return
	}

	WelcomeInterface := "Welcome to TCP-Chat!\n" +
		"         _nnnn_\n" +
		"        dGGGGMMb\n" +
		"       @p~qp~~qMb\n" +
		"       M|@||@) M|\n" +
		"       @,----.JM|\n" +
		"      JS^\\__/  qKL\n" +
		"     dZP        qKRb\n" +
		"    dZP          qKKb\n" +
		"   fZP            SMMb\n" +
		"   HZM            MMMM\n" +
		"   FqM            MMMM\n" +
		" __| \".        |\\dS\"qML\n" +
		" |    `.       | `' \\Zq\n" +
		"_)      \\.___.,|     .'\n" +
		"\\____   )MMMMMP|   .'\n" +
		"     `-'       `--'\n"

	_, err := conn.Write([]byte(WelcomeInterface))
	if err != nil {
		fmt.Println("Error sending welcome message:", err)
		return
	}

	// Read client's name
	buf := make([]byte, 1024)
	var name string
	nameLen := 0
	for nameLen < 2 {

		// Write client current bar
		_, err = conn.Write([]byte("[ENTER YOUR NAME]: "))
		if err != nil {
			fmt.Println("Error while writing '[ENTER YOUR NAME]:':", err)
			return
		}
		nameLen, err = conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading client name:", err)
			onlineClientCount--
			return
		}
		name = string(buf[:nameLen])
		name = name[:len(name)-1] // Remove trailing newline character
	}

	//Load logs for client
	_, err = conn.Write([]byte(log))
	if err != nil {
		fmt.Println("Loading log:", err)
		return
	}

	// Show to all users the new user connection
	broadcastMessage(conn, "", name+" has joined our chat...\n", true)

	// Write client current bar
	_, err = conn.Write([]byte("[" + nowTime() + "][" + name + "]:"))
	if err != nil {
		fmt.Println("Error sending welcome message:", err)
		return
	}

	// Add the client to the map
	mu.Lock()
	clients[conn] = Client{
		ID:          id,
		Name:        name,
		Conn:        conn,
		IsConnected: true,
	}
	mu.Unlock()

	// fmt.Println("New client connected: "+name+" (ID: ", id, ", Address: "+conn.RemoteAddr().String()+")")

	// Handle client communication
	for {
		// Read data from the client
		n, err := conn.Read(buf)
		if err != nil {

			onlineClientCount--
			fmt.Printf("Client %s disconnected: %v\n", name, err)
			broadcastMessage(conn, "", name+" has left our chat...\n", true)
			return
		}

		// Write current user name bar
		_, err = conn.Write([]byte("[" + nowTime() + "][" + name + "]:"))
		if err != nil {
			fmt.Println("Error write current user bar name:", err)
			return
		}

		// Get the message from the client
		message := string(buf[:n])

		if message != "\n" {
			// Broadcast the message to all connected clients
			broadcastMessage(conn, name, message, false)
		}
	}
}

func nowTime() string {
	// Get the current time
	currentTime := time.Now()

	// Format the time in "YYYY-MM-DD HH:mm:ss"
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return formattedTime
}

// broadcastMessage sends the received message to all connected clients
func broadcastMessage(sender net.Conn, senderName, message string, sys bool) {
	mu.Lock()
	defer mu.Unlock()

	movePriviousLine := "\033[F"
	removeLine := "\033[K"

	if !sys {
		message = "[" + nowTime() + "][" + senderName + "]:" + message
	}

	log += message

	
	for clientConn, client := range clients {

		// Don't send the message back to the sender
		if clientConn != sender {

			// Broadcast message
			_, err := clientConn.Write([]byte("\n" + movePriviousLine + removeLine + message))
			if err != nil {
				fmt.Printf("Error writing to client %s: %v\n", client.Name, err)
			}

			// Write current user name bar
			_, err = clientConn.Write([]byte("[" + nowTime() + "][" + client.Name + "]:"))
			if err != nil {
				fmt.Println("Error write current user bar name:", err)
				return
			}
		}
	}
}
