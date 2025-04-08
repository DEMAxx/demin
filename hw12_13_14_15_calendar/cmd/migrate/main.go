package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/DEMAxx/demin/hw12_13_14_15_calendar/internal/config"
	"github.com/DEMAxx/demin/hw12_13_14_15_calendar/migrations"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	conf := config.NewConfig(configFile)

	address := fmt.Sprintf("%s:%s", conf.DB.Host, conf.DB.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to start TCP server: %v", err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Printf("failed to close listener: %v", err)
		}
	}(listener)

	log.Printf("TCP server listening on %s", address)

	if err := migrations.Run(&migrations.Config{
		User:     conf.DB.User,
		Password: conf.DB.Password,
		Host:     conf.DB.Host,
		Port:     conf.DB.Port,
		Name:     conf.DB.Name,
	}); err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}(conn)
}
