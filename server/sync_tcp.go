// Package server provides the user to run synchronous TCP echo server
package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/pawanmna/rediGo/config"
	"github.com/pawanmna/rediGo/core"
)

// Replaced net.Conn with io.Writer/Reader
func respondError(err error, c io.Writer) {
	c.Write((fmt.Appendf(nil, "-%s\r\n", err)))
}

func readCommand(c io.Reader) (*core.RedisCmd, error) {
	buf := make([]byte, 512) // the max size of the message passed to the server can be 512 bytes, make() create a slice

	n, err := c.Read(buf[:]) // Read return the size and error and add the command to buf slice
	if err != nil {
		return nil, err
	}
	token, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}
	return &core.RedisCmd{
		Cmd:  strings.ToUpper(token[0]),
		Args: token[1:],
	}, nil
}

func respond(cmd *core.RedisCmd, c net.Conn) error {
	resp, err := core.Eval(cmd)
	if err != nil {
		respondError(err, c)
		return err
	}

	_, err = c.Write(resp)
	return err
}

func RunSyncTCPServer() {
	log.Println("starting a synchronous TCP server on", config.Host, config.Port)

	conClients := 0                                                         // number of concurrent clients connected to the server at the moment
	ln, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port)) // listens tcp connection with port and host from config.go or input flags
	if err != nil {
		log.Println(err)
	}

	for {
		conn, err := ln.Accept() // accepts the connection
		if err != nil {
			panic(err)
		}

		conClients += 1 // increment the concurrent clients, it will be never >1. because the server won't accept new connection unless older one is disconnected

		log.Println("connected to new client at address:", conn.RemoteAddr(), "concurrent clients", conClients)
		for {
			cmd, err := readCommand(conn) // read command sent by the client
			if err != nil {
				conn.Close()    // if any error close the connection
				conClients -= 1 // decrement the number of concurrent clients
				log.Println("disconnected client", conn.RemoteAddr(), "concurrent clients", conClients)
				if err == io.EOF {
					break // break the inner inf for loop when the connection is broken, it will let the server accept new connections
				}
				log.Println(err)
			}
			log.Println("commands", cmd)
			err = respond(cmd, conn) // echo the commands to the client
			if err != nil {
				log.Print("err write: ", err)
			}
		}

	}
}
