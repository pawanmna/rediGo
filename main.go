package main

import (
	"flag"
	"log"

	"github.com/pawanmna/rediGo/config"
	"github.com/pawanmna/rediGo/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the rediGo server")
	flag.IntVar(&config.Port, "port", 7379, "port for the rediGo server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("starting the server")
	server.RunAsyncTCPServer()
}
