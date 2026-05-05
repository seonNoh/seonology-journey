// Package client provides gRPC client wrappers for external services
// (seonology-auth-back, seonology-member-back) with retry, timeout,
// and circuit breaker support.
package client

import (
	"context"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds connection settings for a gRPC client.
type Config struct {
	Address string
	Timeout time.Duration
}

// Conn wraps a gRPC connection with circuit breaker.
type Conn struct {
	cc *grpc.ClientConn
	cb *gobreaker.CircuitBreaker
}

// Dial creates a new gRPC client connection with circuit breaker.
func Dial(ctx context.Context, name string, cfg Config) (*Conn, error) {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	dialCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cc, err := grpc.DialContext(dialCtx, cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("client %s: dial %s: %w", name, cfg.Address, err)
	}

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
	})

	return &Conn{cc: cc, cb: cb}, nil
}

// ClientConn returns the underlying gRPC connection.
func (c *Conn) ClientConn() *grpc.ClientConn {
	return c.cc
}

// Execute runs fn within the circuit breaker.
func (c *Conn) Execute(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error) {
	return c.cb.Execute(func() (any, error) {
		return fn(ctx)
	})
}

// Close closes the underlying connection.
func (c *Conn) Close() error {
	return c.cc.Close()
}

// State returns the circuit breaker's current state.
func (c *Conn) State() gobreaker.State {
	return c.cb.State()
}
