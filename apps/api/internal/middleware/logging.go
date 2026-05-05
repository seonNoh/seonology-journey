// Package middleware provides HTTP middleware for the api gateway.
package middleware

import (
	"log"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// StructuredLogger logs each request with traceId for Loki/promtail correlation.
func StructuredLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		traceID := chimw.GetReqID(r.Context())
		log.Printf(`{"level":"info","traceId":"%s","method":"%s","path":"%s","status":%d,"duration":"%s","remote":"%s"}`,
			traceID,
			r.Method,
			r.URL.Path,
			ww.Status(),
			time.Since(start).String(),
			r.RemoteAddr,
		)
	})
}
