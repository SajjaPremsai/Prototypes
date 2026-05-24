package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/redis/config"
	"github.com/redis/core"
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
				conClients -= 1
				log.Printf("Client disconnected: %s (Concurrent clients: %d)", c.RemoteAddr(), conClients)
				if err != io.EOF {
					break
				}
				log.Println("Error :", err)
				break
			}
			respond(cmd, c)
		}
	}
}

func readCommand(c net.Conn) (*core.RedisCmd, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf)

	if err != nil {
		return nil, err
	}

	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}

	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, c net.Conn) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmd *core.RedisCmd, c net.Conn){
	err := core.EvalAndRespond(cmd, c)
	if err != nil {
		respondError(err, c)
	}
}