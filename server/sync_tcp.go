package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/pawanmna/rediGo/config"
)

func readCommand(c net.Conn) (string, error) {
	var buf []byte = make([]byte, 512)

	n, err := c.Read(buf[:])
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

func RunSyncTCPServer() {
	log.Println("starting a synchronous TCP server on", config.Host, config.Port)

	con_clients := 0 // number of concurrent clients connected to the server at the moment
	ln, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		log.Println(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		con_clients += 1

		log.Println("connected to new client at address:", conn.RemoteAddr(), "concurrent clients", con_clients)
		cmd, err := readCommand(conn)
		if err != nil {
			conn.Close()
			con_clients -= 1
			log.Println("disconnected client", conn.RemoteAddr(), "concurrent clients", con_clients)
			if err == io.EOF {
				break
			}
			log.Println(err)
		}
		log.Println("commands", cmd)
		if err = respond(cmd, conn); err != nil {
			log.Println("err write: ", err)
		}

	}
}
