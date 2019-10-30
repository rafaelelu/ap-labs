package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("How to execute: ./clockWall [city]=[server ip]:[port]")
	}

	var wg sync.WaitGroup

	for i:= 1; i < len(os.Args); i++ {
		argSlice := strings.Split(os.Args[i], "=")
		city := argSlice[0]
		server := argSlice[1]

		conn, err := net.Dial("tcp",server)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		go connCopy(conn, city, &wg)
	}
	wg.Wait()
}

func connCopy(conn net.Conn, city string, wg *sync.WaitGroup) {
	defer wg.Done()
	for ;; {
		_, err := io.Copy(os.Stdout, conn)
		if err == nil {
			break
		}
	}
	log.Println("Lost connection with " + city)
}