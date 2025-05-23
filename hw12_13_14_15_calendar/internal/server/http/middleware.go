package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

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
