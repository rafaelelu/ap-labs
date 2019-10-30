// Clock2 is a concurrent TCP server that periodically writes the time.
package main

import (
	"io"
	"log"
	"net"
	"time"
	"os"
	"strings"
)

func handleConn(c net.Conn) {
	defer c.Close()

	timeZone := os.Getenv("TZ")
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, err := io.WriteString(c, timeZone+strings.Repeat(" ", 10-len(timeZone)) + " : "+time.Now().In(location).Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("How to execute: TZ=[timezone] go run clock2.go -port [port]")
	}

	port := os.Args[2]
	serverAddress := "localhost:" + port


	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
}