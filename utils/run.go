package utils

import (
	"fmt"
	"net"
	"os"
)

func RunServer() {
	args := os.Args[1:]
	port := "8989"

	var err error
	os.Mkdir("logs", 0777)
	logFile, err = os.Create("logs/" + startingTime + ".txt")
	if err != nil {
		fmt.Print(err)
		return
	}
	defer logFile.Close()
	savelogs("Server Started at: " + startingTime + "\n")
	if len(args) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(0)
	}

	if len(args) == 1 {
		port = args[0]
	}
	// Listen for incoming connections on port 8080
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(0)
	}
	defer listen.Close()

	fmt.Println("Listening on the port :" + port)
	go exit()

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
