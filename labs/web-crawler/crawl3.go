// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
//
// Crawl3 adds support for depth limiting.
//
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"

	"github.com/todostreaming/gopl.io/ch5/links"
)

type paths struct {
	urls []string
	father int
}

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string, f int) paths {
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}
	pathsToReturn := paths{urls: list, father: f+1}
	return pathsToReturn
}

//!-sema

//!+
func main() {
	if len(os.Args) < 3 {
		log.Fatalf("How to execute: go run crawl3.go -depth=[int] [url]\n")
	}

	checkFlag := strings.Split(os.Args[1], "=")[0]
	if (checkFlag != "-depth") {
		log.Fatalf("depth should be an integer greater than zero")
	}

	depthString := strings.Split(os.Args[1], "=")[1]
	depth, err := strconv.Atoi(depthString)
	if err != nil {
		log.Fatalf("depth should be an integer greater than zero")
	}



	worklist := make(chan paths)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++
	firstpath := paths{urls: os.Args[2:], father: 1}

	go func() { worklist <- firstpath}()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		p := <-worklist
		for _, link := range p.urls {
			if !seen[link] {
				seen[link] = true
				n++
				if p.father <= depth {
					go func(link string) {
						worklist <- crawl(link, p.father)
					}(link)
				} else {
					return
				}
			}
		}
	}
}

//!-