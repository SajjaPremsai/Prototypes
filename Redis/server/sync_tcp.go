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
				conClients--
				log.Printf("Client disconnected: %s (Concurrent clients: %d)", c.RemoteAddr(), conClients)
				if err != io.EOF {
					log.Println("read error", err)
				}
				break
			}

			if cmd == nil {
				continue
			}

			if err := respond(cmd, c); err != nil {
				respondError(err, c)
				break
			}
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
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err.Error())))
}

func respond(cmd *core.RedisCmd, c net.Conn) error {
	return core.EvalAndRespond(cmd, c)
}
