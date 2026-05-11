// Package server provides the user to run asynchronous TCP echo server
package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/pawanmna/rediGo/config"
)

var concurrentClients int64

func handleConn(conn net.Conn) {
	defer conn.Close()

	n := atomic.AddInt64(&concurrentClients, 1)
	log.Println("connected to new client at address:", conn.RemoteAddr(), "concurrent clients", n)

	defer func() {
		n := atomic.AddInt64(&concurrentClients, -1)
		log.Println("disconnected client", conn.RemoteAddr(), "concurrent clients", n)
	}()

	for {
		cmd, err := readCommand(conn)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("read error:", err)
			return
		}

		log.Println("command:", cmd)

		if err := respond(cmd, conn); err != nil {
			log.Println("write/eval error:", err)
			return
		}
	}
}

func RunAsyncTCPServer() {
	addr := config.Host + ":" + strconv.Itoa(config.Port)
	log.Println("starting an asynchronous TCP server on", addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		go handleConn(conn)
	}
}
