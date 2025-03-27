package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var intervalFlag time.Duration

	pflag.DurationVarP(&intervalFlag, "timeout", "t", time.Second*10, "timeout of each event")

	pflag.Parse()

	args := pflag.Args()

	if len(args) != 2 {
		log.Fatal("Usage: go-telnet [--timeout=10s] host port")
	}

	client := NewTelnetClient(net.JoinHostPort(args[0], args[1]), intervalFlag, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	errChan := make(chan error, 2)

	go func() {
		errChan <- client.Send()
	}()

	go func() {
		errChan <- client.Receive()
	}()

	select {
	case <-sigChan:
		fmt.Fprintln(os.Stderr, "\n...Connection was closed by client")
	case err := <-errChan:
		if err != nil {
			fmt.Fprintln(os.Stderr, "...Connection was closed by server")
		}
	}
}
