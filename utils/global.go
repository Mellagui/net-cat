package utils

import (
	"net"
	"os"
	"sync"
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
	startingTime      = nowTime()
	logFile           *os.File
)
