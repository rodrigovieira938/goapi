package util

import (
	"fmt"
	"net"
)

func IsPortAvailable(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

func FindUsablePort(start int) int {
	for port := start; port < 65535; port++ {
		if IsPortAvailable(port) {
			return port
		}
	}
	return 0 // No available port found
}
