package client

import (
	"context"
	"log/slog"
)

// AuthClient wraps gRPC calls to seonology-auth-back.
type AuthClient struct {
	conn *Conn
}

// NewAuthClient creates an AuthClient.
// Address should be the k8s service DNS: seonology-auth-back.seonology-auth:9090
func NewAuthClient(ctx context.Context, cfg Config) (*AuthClient, error) {
	conn, err := Dial(ctx, "auth", cfg)
	if err != nil {
		return nil, err
	}
	slog.Info("auth client connected", "addr", cfg.Address)
	return &AuthClient{conn: conn}, nil
}

// ValidateToken validates a bearer token via the auth service.
// Returns the user ID if valid.
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (userID string, err error) {
	result, err := c.conn.Execute(ctx, func(ctx context.Context) (any, error) {
		// TODO: call auth gRPC ValidateToken RPC
		_ = token
		return "", nil
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// Close closes the connection.
func (c *AuthClient) Close() error {
	return c.conn.Close()
}
