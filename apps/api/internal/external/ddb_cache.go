package external

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DDBCache implements Cache using a DynamoDB table with TTL.
type DDBCache struct {
	client *dynamodb.Client
	table  string
}

// NewDDBCache creates a DynamoDB-backed cache.
// The table must have PK (String) as partition key and a TTL attribute named "ttl".
func NewDDBCache(client *dynamodb.Client, table string) *DDBCache {
	return &DDBCache{client: client, table: table}
}

type ddbCacheItem struct {
	PK   string `dynamodbav:"PK"`
	Data []byte `dynamodbav:"data"`
	TTL  int64  `dynamodbav:"ttl"`
}

func (c *DDBCache) Get(ctx context.Context, key string) ([]byte, error) {
	out, err := c.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(c.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: key},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ddb cache get: %w", err)
	}
	if out.Item == nil {
		return nil, nil
	}
	var it ddbCacheItem
	if err := attributevalue.UnmarshalMap(out.Item, &it); err != nil {
		return nil, err
	}
	// Check if expired (DDB TTL deletion is eventual)
	if time.Now().Unix() > it.TTL {
		return nil, nil
	}
	return it.Data, nil
}

func (c *DDBCache) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	it := ddbCacheItem{
		PK:   key,
		Data: data,
		TTL:  time.Now().Add(ttl).Unix(),
	}
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return fmt.Errorf("ddb cache marshal: %w", err)
	}
	_, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.table),
		Item:      av,
	})
	return err
}
