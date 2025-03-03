package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func readRoutine(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(conn)
	defer log.Printf("Finished readRoutine")

OUTER:
	for {
		select {
		case <-ctx.Done():
			log.Printf("readRoutine is done")
			return
		default:
			if !scanner.Scan() {
				log.Printf("CANNOT SCAN")
				break OUTER
			}
			text := scanner.Text()
			log.Printf("From server: %s", text)
		}
	}
}

func writeRoutine(ctx context.Context, conn net.Conn, wg *sync.WaitGroup, stdin chan string) {
	defer wg.Done()
	defer log.Printf("Finished writeRoutine")

	for {
		select {
		case <-ctx.Done():
			return
		case str := <-stdin:
			fmt.Printf("To server %v\n", str)

			conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		}

	}
}

func stdinScan() chan string {
	out := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		if scanner.Err() != nil {
			close(out)
		}
	}()
	return out
}

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	var err error
	var intervalFlag time.Duration

	pflag.DurationVarP(&intervalFlag, "timeout", "t", time.Second*10, "timeout of each event")

	pflag.Lookup("timeout").NoOptDefVal = "noconfig.json"

	pflag.Parse()

	dialer := &net.Dialer{}

	ctx, cancel := context.WithTimeout(context.Background(), intervalFlag)
	defer cancel()

	//ctxNotify, stop := signal.NotifyContext(ctx, syscall.SIGTERM)
	//defer stop()

	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort("localhost", "4242"))

	if err != nil {
		log.Fatalf("Can not connect: %v", err)
	}

	defer conn.Close()

	//client := NewTelnetClient(
	//	net.JoinHostPort("localhost", "4242"),
	//	intervalFlag,
	//	conn,
	//	conn,
	//)
	//
	//_ = client

	//conn.Write([]byte(fmt.Sprintf("%s\n", "You are connected")))

	stdin := stdinScan()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		readRoutine(ctx, conn, wg)
		cancel()
	}()

	wg.Add(1)
	go func() {
		writeRoutine(ctx, conn, wg, stdin)
	}()

	wg.Wait()
}
