package util

import (
	"fmt"
	"net"
	"time"
)

// IsPortOpen returns true if port is open
func IsPortOpen(host string, port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func WaitForPort(host string, port int, to int) bool {

	timeout := time.After(time.Duration(to) * time.Second)
	tick := time.Tick(2 * time.Second)

	for {
		select {
		case <-timeout:
			return false
		case <-tick:
			if IsPortOpen(host, port) {
				return true
			}
		}
	}
	return false
}
