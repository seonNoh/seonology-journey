package ddb

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// EncodeCursor encodes a DynamoDB LastEvaluatedKey into a URL-safe base64 string.
// Returns empty string if the key is nil (no more pages).
func EncodeCursor(key map[string]types.AttributeValue) (string, error) {
	if key == nil {
		return "", nil
	}
	// Convert to a simple map[string]string for portability.
	simple := make(map[string]string, len(key))
	for k, v := range key {
		switch av := v.(type) {
		case *types.AttributeValueMemberS:
			simple[k] = av.Value
		case *types.AttributeValueMemberN:
			simple[k] = av.Value
		default:
			return "", fmt.Errorf("unsupported attribute type for cursor key %q", k)
		}
	}
	data, err := json.Marshal(simple)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

// DecodeCursor decodes a base64 cursor string back into a DynamoDB ExclusiveStartKey.
// Returns nil if cursor is empty.
func DecodeCursor(cursor string) (map[string]types.AttributeValue, error) {
	if cursor == "" {
		return nil, nil
	}
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor encoding: %w", err)
	}
	var simple map[string]string
	if err := json.Unmarshal(data, &simple); err != nil {
		return nil, fmt.Errorf("invalid cursor payload: %w", err)
	}
	result := make(map[string]types.AttributeValue, len(simple))
	for k, v := range simple {
		result[k] = &types.AttributeValueMemberS{Value: v}
	}
	return result, nil
}
