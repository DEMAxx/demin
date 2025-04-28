package internalhttp

import (
	"context"
	"fmt"
	"github.com/DEMAxx/demin/hw12_13_14_15_calendar/events/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func ValidationInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	switch r := req.(type) {
	case *pb.EventCreateRequest:
		if err := validateEventCreateRequest(r); err != nil {
			return nil, status.Errorf(400, "validation error: %v", err)
		}
		// Add cases for other request types as needed
	}
	return handler(ctx, req)
}

func LoggingMiddleware(next http.Handler, logg Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		dateTime := time.Now().Format(time.RFC3339)
		method := r.Method
		path := r.URL.Path
		httpVersion := r.Proto
		userAgent := r.Header.Get("User-Agent")

		logg.Info(
			fmt.Sprintf(
				"Client IP: %s, DateTime: %s, Method: %s, Path: %s, HTTP Version: %s, User Agent: %s",
				clientIP, dateTime, method, path, httpVersion, userAgent,
			),
		)

		next.ServeHTTP(w, r)
	})
}
