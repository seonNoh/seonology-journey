package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// MarshalItem converts a Go struct (with dynamodbav tags) to a DynamoDB item map.
func MarshalItem(v any) (map[string]types.AttributeValue, error) {
	item, err := attributevalue.MarshalMap(v)
	if err != nil {
		return nil, fmt.Errorf("marshal DDB item: %w", err)
	}
	return item, nil
}

// UnmarshalItem converts a DynamoDB item map to a Go struct.
func UnmarshalItem(item map[string]types.AttributeValue, out any) error {
	if err := attributevalue.UnmarshalMap(item, out); err != nil {
		return fmt.Errorf("unmarshal DDB item: %w", err)
	}
	return nil
}

// MarshalList converts a slice of Go structs to a slice of DynamoDB item maps.
func MarshalList[T any](items []T) ([]map[string]types.AttributeValue, error) {
	result := make([]map[string]types.AttributeValue, 0, len(items))
	for _, item := range items {
		m, err := attributevalue.MarshalMap(item)
		if err != nil {
			return nil, fmt.Errorf("marshal DDB list item: %w", err)
		}
		result = append(result, m)
	}
	return result, nil
}

// UnmarshalList converts a slice of DynamoDB item maps to a slice of Go structs.
func UnmarshalList[T any](items []map[string]types.AttributeValue) ([]T, error) {
	result := make([]T, 0, len(items))
	for _, item := range items {
		var v T
		if err := attributevalue.UnmarshalMap(item, &v); err != nil {
			return nil, fmt.Errorf("unmarshal DDB item: %w", err)
		}
		result = append(result, v)
	}
	return result, nil
}
