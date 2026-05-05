// Package ddb provides the shared DynamoDB client and utilities for
// the seonology-journey-back service.
package ddb

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	once     sync.Once
	instance *dynamodb.Client
)

// TablePrefix is prepended to all logical table names.
const TablePrefix = "seonology-journey-"

// Client returns the singleton DynamoDB client. It initializes the client
// on first call using AWS SDK default credential chain with standard retry.
func Client(ctx context.Context) *dynamodb.Client {
	once.Do(func() {
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(region()),
			config.WithRetryer(func() aws.Retryer {
				return retry.NewStandard(func(o *retry.StandardOptions) {
					o.MaxAttempts = 5
				})
			}),
		)
		if err != nil {
			slog.Error("failed to load AWS config", "error", err)
			os.Exit(1)
		}

		opts := []func(*dynamodb.Options){}
		if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
			opts = append(opts, func(o *dynamodb.Options) {
				o.BaseEndpoint = aws.String(endpoint)
			})
		}

		instance = dynamodb.NewFromConfig(cfg, opts...)
	})
	return instance
}

// TableName returns the full table name with prefix.
func TableName(logical string) string {
	return TablePrefix + logical
}

func region() string {
	if r := os.Getenv("AWS_REGION"); r != "" {
		return r
	}
	return "ap-northeast-1"
}
