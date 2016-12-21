// -*- Mode: Go; indent-tabs-mode: t -*-

package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func echoServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			log.Fatal("Read: ", err)
		}

		data := buf[0:nr]
		fmt.Println("Server got: ", string(data))
		_, err = c.Write(data)
		if err != nil {
			log.Fatal("Write: ", err)
		}
	}		
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Usage unix-domain-socket <socketdir>")
	}

	dir := os.Args[1]
	path := dir + "/socket"

	fmt.Println("socket path is ", path)

	l, err := net.Listen("unix", path)
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		// wait for connection
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}
		// handle the connection in a new goroutine.
		// the loop then returns to accepting, so that
		// multiple connections may be server conncurrently.
		go echoServer(conn)
	}
	
	fmt.Println("all done!")
}