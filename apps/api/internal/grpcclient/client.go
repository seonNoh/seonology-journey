// Package grpcclient - back gRPC 클라이언트 다이얼.
package grpcclient

import (
	"context"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Dial - back gRPC 다이얼.
func Dial(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// NewJourneyClient - JourneyService 클라이언트.
func NewJourneyClient(conn *grpc.ClientConn) journeyv1.JourneyServiceClient {
	return journeyv1.NewJourneyServiceClient(conn)
}

// WithUser - x-user-id metadata 부착.
func WithUser(ctx context.Context, userID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "x-user-id", userID)
}
