// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

//!+broadcaster
type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

var kickedUsers = make([]chan string, 0)

var clients = make(map[client]bool) // all connected clients
func broadcaster() {
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func notKicked(userChannel chan<- string) bool {
	for i := 0; i < len(kickedUsers); i++ {
		if userChannel == kickedUsers[i] {
			return false
		}
	}
	return true
}

//!-broadcaster

//!+handleConn
type user struct {
	name string
	address       string
	channel  chan string
}

var users = make([]user, 0)
var adminAddress string


func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	input := bufio.NewScanner(conn)
	input.Scan()
	name := input.Text()
	var newUser user
	newUser.address = conn.RemoteAddr().String()
	newUser.channel = ch
	newUser.name = name
	users = append(users, newUser)
	
	ch <- "irc-server > Welcome to the Simple IRC Server"
	ch <- "irc-server > Your user " + name + " is successfully logged"


	messages <- "irc-server > New connected user " + name
	entering <- ch
	fmt.Println("irc-server > New connected user " + name)

	if len(users) == 1 {
		adminAddress = users[0].address
		ch <- "irc-server > Congrats, you were the first user"
		ch <- "irc-server > You're the new IRC Server ADMIN"
		fmt.Println("irc-server > " + name + " was promoted as the channel ADMIN")
	}
	go clientWriter(conn, ch)

	
	for input.Scan() {
		incoming := strings.Split(input.Text(), " ")
		if notKicked(ch) {
			switch incoming[0] {
			case "/users":
				lstOfUsrs := "irc-server > Users: "
				i := 0
				for ; i < len(users) - 1; i++ {
					lstOfUsrs += users[i].name + ", "
				}
				lstOfUsrs += users[i].name
				ch <- lstOfUsrs
			case "/msg":
				if(len(incoming) > 1) {
					for i, user := range users {
						if user.name == incoming[1] {
							dm := "Message from " + name + " > "
							for _, dmContent := range incoming[2:] {
								dm += (dmContent + " ")
							}
							user.channel <- dm
							break
						}
						if i == len(users)-1 {
							ch <- "irc-server> User not found"
						}
					}
				} else {
					ch <- "irc-server > ERROR: Missing argument. Correct user: /msg [user]"
				}
			case "/time":
				ch <- "irc-server > Local Time: " + "America/Mexico_City " + time.Now().Format("15:04")
			case "/user":
				if(len(incoming) > 1) {
					for i, user := range users {
						if user.name == incoming[1] {
							ip := strings.Split(user.address, ":")[0] // gets only the ip and not the port in the target user's address
							ch <- "irc-server > Username: " + user.name + " IP Address: " + ip
							break
						}
						if i == len(users) - 1 {
							ch <- "irc-server > User not found"
						}
					}
				} else {
					ch <- "irc-server > ERROR: Missing argument. Correct user: /user [user]"
				}
			case "/kick":
				if(len(incoming) > 1) {
					if isAdmin(ch) {
						foundUser := kickUser(incoming[1])
						if foundUser {
							messages <- "irc-server > " + incoming[1] + " was kicked from channel by the channel Admin"
						} else {
							ch <- "irc-server > User not found"
						}
					} else {
						ch <- "irc-server > You can't kick users, because you are not the admin"
					}
				} else {
					ch <- "irc-server > ERROR: Missing argument. Correct user: /kick [user]"
				}
			default:
				if notKicked(ch) {
					messages <- name + " > " + input.Text()
				}
			}
		}
	}

	leaving <- ch
	messages <- "irc-server > " + name + " left channel"
	fmt.Println("irc-server > " + name + " left")
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func isAdmin(channel chan string) bool {
	for i := 0; i < len(users); i++ {
		if (users[i].channel == channel) && (users[i].address == adminAddress) {
			return true
		}
	}
	return false
}

func kickUser(name string) bool{
	foundUser := false
	for i := 0; i < len(users); i++ {
		if (users[i].name == name) {
			users[i].channel <- "irc-server > You've been kicked from this channel"
			kickedUsers = append(kickedUsers, users[i].channel)
			leaving <- users[i].channel
			users = append(users[:i], users[i+1:]...)
			foundUser = true
		}
	}
	return foundUser
}

//!-handleConn

//!+main
func main() {
	if (len(os.Args) < 5) || (os.Args[1] != "-host") || (os.Args[3] != "-port") {
		log.Fatalf("How to execute: go run server.go -host [server ip address] -port [port]\n")
	}
	listener, err := net.Listen("tcp", os.Args[2] + ":" + os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("irc-server > Simple IRC Server started at localhost:9000")
	fmt.Println("irc-server > Ready for receiving new clients")

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main