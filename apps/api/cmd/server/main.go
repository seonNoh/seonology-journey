// Package main is the entrypoint for the seonology-journey-api REST + WebSocket gateway.
package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("HTTP_LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	log.Printf("seonology-journey-api HTTP listening on %s", addr)
	srv := &http.Server{Addr: addr, Handler: mux}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
