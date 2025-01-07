package utils

import (
	"fmt"
	"net"
)

// broadcastMessage sends the received message to all connected clients
func broadcastMessage(sender net.Conn, senderName, message string, sys bool, clear bool) {
	mu.Lock()
	defer mu.Unlock()

	movePriviousLine := "\033[F"
	removeLine := "\033[K"
	chatClear := "\x1b[3J\x1b[H\x1b[2J"

	if !sys {
		message = "[" + nowTime() + "][" + senderName + "]:" + message
	}

	log += message

	savelogs(message)

	for clientConn, client := range clients {
		if clear && clientConn == sender {
			_, err := clientConn.Write([]byte("\n" + chatClear + message))
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

		// Don't send the message back to the sender
		if !clear {

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
