// -*- Mode: Go; indent-tabs-mode: t -*-

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
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

	// Now run timedate-ctrl twice enable/disable
	cmd := exec.Command("timedatectl", "set-ntp", "true")
	err := cmd.Run()
	if err != nil {
		fmt.Println("timedatectl error cmd #1: %v", err)
	}

	cmd = exec.Command("timedatectl", "set-ntp", "false")
	err = cmd.Run()
	if err != nil {
		fmt.Println("timedatectl error cmd #2: %v", err)
	}

	cmd = exec.Command("timedatectl", "status")
	err = cmd.Run()
	if err != nil {
		fmt.Println("timedatectl error cmd #3: %v", err)
	}

	cmd = exec.Command("timedatectl", "set-timezone", "America/Tijuana")
	err = cmd.Run()
	if err != nil {
		fmt.Println("timedatectl error cmd #4: %v", err)
	}

	cmd = exec.Command("timedatectl", "set-time", "2017-10-31 18:17:16")
	err = cmd.Run()
	if err != nil {
		fmt.Println("timedatectl error cmd #5: %v", err)
	}

	fmt.Println("timedatectl all seemed to work")

	dir := os.Args[1]
	path := dir + "/sockets/socket"

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
