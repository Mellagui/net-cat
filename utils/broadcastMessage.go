package utils

import (
	"fmt"
)

// broadcastMessage sends the received message to all connected clients
func broadcastMessage(senderName, message string, sys bool) {
	mu.Lock()
	defer mu.Unlock()

	movePriviousLine := "\033[F"
	removeLine := "\033[K"

	if !sys {
		message = "[" + nowTime() + "][" + senderName + "]:" + message
	}

	log += message

	savelogs(message)

	for clientConn, client := range clients {
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
