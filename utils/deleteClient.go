package utils

import "net"

func DeleteClient(conn net.Conn) {
	// Remove client from the map and close the connection
	mu.Lock()
	delete(clients, conn)
	mu.Unlock()
	conn.Close()
}
