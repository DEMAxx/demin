package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
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
	logg := logger.New(config.Logger.Level)

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
		fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		calendar,
	)

	latencyMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			next.ServeHTTP(w, r)
			latency := time.Since(startTime)
			logg.Info(fmt.Sprintf("Latency: %s", latency))
		})
	}

	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		dateTime := time.Now().Format(time.RFC3339)
		method := r.Method
		path := r.URL.Path
		httpVersion := r.Proto
		userAgent := r.Header.Get("User-Agent")

		logg.Info(fmt.Sprintf("Client IP: %s, DateTime: %s, Method: %s, Path: %s, HTTP Version: %s, User Agent: %s", clientIP, dateTime, method, path, httpVersion, userAgent))

		write, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			return
		}
		logg.Info(fmt.Sprintf("response: %d", write))
	})

	http.Handle("/hello", latencyMiddleware(helloHandler))

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
