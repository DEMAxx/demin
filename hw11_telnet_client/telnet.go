package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	out.Write([]byte(fmt.Sprintf("%s\n", "You are connected")))
	fmt.Printf("NewTelnetClient: You are connected")

	stdin := stdinScan()

	for {
		select {
		case str := <-stdin:
			log.Println("STDIN", str)
		case <-time.After(timeout - time.Second):
			log.Println("got timeout", timeout)
			in.Close()
			return nil
		}
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
func handleConnection(conn net.Conn) {
	fmt.Fprintf(conn, "Welcome to %s, friend from %s\n", conn.LocalAddr(), conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		log.Printf("RECEIVED: %s", text)
		if text == "quit" || text == "exit" {
			break
		}

		conn.Write([]byte(fmt.Sprintf("I have received '%s'\n", text)))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error happend on connection with %s: %v", conn.RemoteAddr(), err)
	}

	log.Printf("Closing connection with %s", conn.RemoteAddr())
}
