package utils

import (
	"fmt"
	"net"
	"strings"
)

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
		name = string(buf[:nameLen-1])
	}

	//Load logs for client
	_, err = conn.Write([]byte(log))
	if err != nil {
		fmt.Println("Loading log:", err)
		return
	}

	// Show to all users the new user connection
	broadcastMessage("", name+" has joined our chat...\n", true)

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

	// Handle client communication
	for {
		// Read data from the client
		n, err := conn.Read(buf)
		if err != nil {

			onlineClientCount--
			fmt.Printf("Client %s disconnected: %v\n", name, err)
			broadcastMessage("", name+" has left our chat...\n", true)
			return
		}

		// Get the message from the client
		message := string(buf[:n])

		// Handle flags
		isFlageName := strings.HasPrefix(message, "--name=")
		isFlageClear := message[:len(message)-1] == "--clear"

		oldName := name
		if message != "\n" {
			if isFlageName {
				lastName := strings.TrimPrefix(message[:n-1], "--name=")
				clients[conn] = Client{Name: lastName}
				name = lastName
			}
		}

		if message != "\n" {

			movePriviousLine := "\033[F"
			removeLine := "\033[K"
	
			// Write current user name bar
			_, err = conn.Write([]byte(movePriviousLine + removeLine))
			if err != nil {
				fmt.Println("Error write current user bar name:", err)
				return
			}

			if isFlageName {
				broadcastMessage(name, oldName+" has changed his name: "+name+"\n", true)
			} else if isFlageClear {
				chatClear := "\033[3J\033[2J\033[H"
				_, err := conn.Write([]byte(chatClear))
				if err != nil {
					fmt.Printf("Error writing to client %s: %v\n", name, err)
				}
				// Write current user name bar
				_, err = conn.Write([]byte("[" + nowTime() + "][" + name + "]:"))
				if err != nil {
					fmt.Println("Error write current user bar name:", err)
					return
				}
			} else {
				// Broadcast the message to all connected clients
				broadcastMessage(name, message, false)
			}
		} else {
			// Write client current bar
			_, err = conn.Write([]byte("[" + nowTime() + "][" + name + "]:"))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
