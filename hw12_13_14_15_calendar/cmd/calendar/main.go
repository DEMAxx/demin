package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DEMAxx/demin/hw12_13_14_15_calendar/internal/app"
	"github.com/DEMAxx/demin/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/DEMAxx/demin/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/DEMAxx/demin/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logg := logger.New(ctx, config.Logger.Level, nil, true)

	switch config.Logger.Output {
	case "stderr":
		logg = logg.Output(os.Stderr)
	case "stdout":
		logg = logg.Output(os.Stdout)
	default:
		logg = logg.Output(os.Stdout) // Default to stdout if not specified
	}

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(
		logg,
		net.JoinHostPort(config.Server.Host, config.Server.Port),
		calendar,
	)

	ctx, cancel = signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	logg.Info("calendar is running...")

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()
}
