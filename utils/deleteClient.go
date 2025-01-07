package utils

import "net"

// Remove client from the map and close the connection
func DeleteClient(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	delete(clients, conn)
	conn.Close()
}
