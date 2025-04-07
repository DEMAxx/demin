package main

import (
	"flag"
	"fmt"

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

	fmt.Println("user", conf.DB.User)
	fmt.Println("password", conf.DB.Password)
	fmt.Println("host", conf.DB.Host)
	fmt.Println("port", conf.DB.Port)
	fmt.Println("name", conf.DB.Name)

	if err := migrations.Run(&migrations.Config{
		User:     conf.DB.User,
		Password: conf.DB.Password,
		Host:     conf.DB.Host,
		Port:     conf.DB.Port,
		Name:     conf.DB.Name,
	}); err != nil {
		panic(err)
	}
}
