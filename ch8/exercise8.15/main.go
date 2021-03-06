//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

//!+broadcaster
type client struct {
	ch   chan<- string // an outgoing message channel
	name string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-entering:
			clients[cli] = true
			cli.ch <- "Current online users:"
			for user := range clients {
				cli.ch <- user.name
			}

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string, 5) // use buffered channel here
	go clientWriter(conn, ch)

	ch <- "Please input your name:"
	input := bufio.NewScanner(conn)
	input.Scan()
	who := input.Text()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- client{ch: ch, name: who}

	timer := time.NewTimer(time.Minute)
	go func() {
		<-timer.C
		conn.Close()
	}()

	for input.Scan() {
		messages <- who + ": " + input.Text()
		timer.Reset(time.Minute)
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- client{ch: ch, name: who}
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

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
