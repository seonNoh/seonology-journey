package client

import (
	"context"
	"log/slog"
)

// MemberClient wraps gRPC calls to seonology-member-back.
type MemberClient struct {
	conn *Conn
}

// NewMemberClient creates a MemberClient.
// Address should be: seonology-member-back.seonology-member:9090
func NewMemberClient(ctx context.Context, cfg Config) (*MemberClient, error) {
	conn, err := Dial(ctx, "member", cfg)
	if err != nil {
		return nil, err
	}
	slog.Info("member client connected", "addr", cfg.Address)
	return &MemberClient{conn: conn}, nil
}

// GetUser retrieves user info by ID from the member service.
func (c *MemberClient) GetUser(ctx context.Context, userID string) (displayName string, err error) {
	result, err := c.conn.Execute(ctx, func(ctx context.Context) (any, error) {
		// TODO: call member gRPC GetUser RPC
		_ = userID
		return "", nil
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// Close closes the connection.
func (c *MemberClient) Close() error {
	return c.conn.Close()
}
