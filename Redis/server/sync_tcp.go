package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/redis/config"
)

func RunSyncTCPServer() {
	log.Printf("Starting a synchronous TCP server on %s:%d", config.Host, config.Port)
	var conClients int
	lsnr, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		panic(err)
	}
	for {
		c, err := lsnr.Accept()
		if err != nil {
			panic(err)
		}
		conClients++
		log.Printf("Client connected: %s (Concurrent clients: %d)", c.RemoteAddr(), conClients)

		for {
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				conClients--
				log.Printf("Client disconnected: %s (Concurrent clients: %d)", c.RemoteAddr(), conClients)
				if err != io.EOF {
					log.Println("Error reading command:", err)
				}
				break
			}
			log.Println("Received command:", cmd)
			if err = respond(cmd, c); err != nil {
				log.Println("Error writing response:", err)
				break
			}
		}
	}
}

func readCommand(c net.Conn) (string, error) {
	buf := make([]byte, 512)
	n, err := c.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}
