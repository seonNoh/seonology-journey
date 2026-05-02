// Package main is the entrypoint for the seonology-journey-back gRPC service.
//
// 実装は段階的に追加. 現時点ではビルド可能性の確保のみが目的.
package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	addr := os.Getenv("GRPC_LISTEN_ADDR")
	if addr == "" {
		addr = ":9090"
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	healthgrpc.RegisterHealthServer(srv, health.NewServer())
	reflection.Register(srv)
	log.Printf("seonology-journey-back gRPC listening on %s", addr)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
