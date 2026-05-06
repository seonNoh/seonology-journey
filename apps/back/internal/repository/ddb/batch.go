package ddb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// maxBatchWriteItems is the hard limit enforced by BatchWriteItem.
const maxBatchWriteItems = 25

// Key represents a DynamoDB primary key pair.
type Key struct {
	PK string
	SK string
}

// BatchDelete removes many items from the same table using BatchWriteItem,
// chunking requests at the 25-item limit and retrying any UnprocessedItems.
//
// Returns the first error encountered; callers should treat partial failures
// as retryable. On success, all keys have been deleted or never existed.
func BatchDelete(ctx context.Context, client *dynamodb.Client, table string, keys []Key) error {
	if len(keys) == 0 {
		return nil
	}
	for start := 0; start < len(keys); start += maxBatchWriteItems {
		end := start + maxBatchWriteItems
		if end > len(keys) {
			end = len(keys)
		}
		chunk := keys[start:end]
		writes := make([]types.WriteRequest, 0, len(chunk))
		for _, k := range chunk {
			writes = append(writes, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: k.PK},
						"SK": &types.AttributeValueMemberS{Value: k.SK},
					},
				},
			})
		}
		unprocessed := map[string][]types.WriteRequest{table: writes}
		// Retry unprocessed items up to 3 times. SDK standard retryer already
		// handles throttling on the whole call; this loop handles partial batches.
		for attempt := 0; attempt < 3 && len(unprocessed[table]) > 0; attempt++ {
			out, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
				RequestItems: unprocessed,
			})
			if err != nil {
				return fmt.Errorf("batch delete %s: %w", table, err)
			}
			unprocessed = out.UnprocessedItems
			if len(unprocessed[table]) == 0 {
				break
			}
		}
		if left := len(unprocessed[table]); left > 0 {
			return fmt.Errorf("batch delete %s: %d items unprocessed after retries", table, left)
		}
	}
	return nil
}

// BatchDeletePK builds Keys where every entry shares the same PK and only SK
// differs. Handy for cascading deletes within a partition.
func BatchDeletePK(pk string, sks []string) []Key {
	out := make([]Key, 0, len(sks))
	for _, sk := range sks {
		out = append(out, Key{PK: pk, SK: sk})
	}
	return out
}

// MustS is a tiny helper for AttributeValueMemberS keys in tests / simple code.
func MustS(v string) types.AttributeValue { return &types.AttributeValueMemberS{Value: v} }
