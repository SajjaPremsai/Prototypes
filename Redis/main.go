package main

import (
	"flag"
	"log"

	"github.com/redis/config"
	"github.com/redis/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host for redis server to listen on")
	flag.IntVar(&config.Port, "port", 7379, "Port for redis server to listen on")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Printf("Starting Redis server on %s:%d", config.Host, config.Port)
	server.RunAsyncTCPServer()
}
