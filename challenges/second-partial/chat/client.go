// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"io"
	"log"
	"net"
	"os"
	"fmt"
)

var user string

//!+
func main() {
	if (len(os.Args) < 5) || (os.Args[1] != "-user") || (os.Args[3] != "-server") {
		log.Fatalf("How to execute: go run client.go -user [user] -server [ip]:[port]\n")
	}

	user = os.Args[2]
	server := os.Args[4]

	conn, err := net.Dial("tcp",server)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.WriteString(conn, user+"\n")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		fmt.Println("Server was shutdown")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
