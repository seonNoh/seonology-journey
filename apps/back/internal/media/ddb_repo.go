package media

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/seonNoh/seonology-journey/apps/back/internal/repository/ddb"
)

const (
	mediaTable   = "media"
	mediaIDIndex = "IdIndex"
)

type mediaItem struct {
	PK             string  `dynamodbav:"PK"`
	SK             string  `dynamodbav:"SK"`
	ID             string  `dynamodbav:"id"`
	TripID         string  `dynamodbav:"tripId"`
	DayID          string  `dynamodbav:"dayId,omitempty"`
	ScheduleID     string  `dynamodbav:"scheduleId,omitempty"`
	S3Key          string  `dynamodbav:"s3Key"`
	ThumbnailS3Key string  `dynamodbav:"thumbnailS3Key,omitempty"`
	MimeType       string  `dynamodbav:"mimeType"`
	Size           int64   `dynamodbav:"size"`
	TakenAt        string  `dynamodbav:"takenAt,omitempty"`
	Latitude       float64 `dynamodbav:"latitude,omitempty"`
	Longitude      float64 `dynamodbav:"longitude,omitempty"`
	Caption        string  `dynamodbav:"caption,omitempty"`
	CreatedAt      string  `dynamodbav:"createdAt"`
	UpdatedAt      string  `dynamodbav:"updatedAt"`
}

// DDBRepo is the DynamoDB implementation of media.Repository.
type DDBRepo struct {
	client *dynamodb.Client
	table  string
}

func NewDDBRepo(client *dynamodb.Client) *DDBRepo {
	return &DDBRepo{client: client, table: ddb.TableName(mediaTable)}
}

func (r *DDBRepo) Save(ctx context.Context, m *journeyv1.Media) error {
	item := mediaToItem(m)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("media ddb save marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("media ddb save: %w", err)
	}
	return nil
}

// Get resolves a media item by ID via the IdIndex GSI.
func (r *DDBRepo) Get(ctx context.Context, id string) (*journeyv1.Media, error) {
	keyCond := expression.Key("id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("media ddb get build: %w", err)
	}
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		IndexName:                 aws.String(mediaIDIndex),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("media ddb get query: %w", err)
	}
	if len(out.Items) == 0 {
		return nil, fmt.Errorf("media: not found")
	}
	var item mediaItem
	if err := attributevalue.UnmarshalMap(out.Items[0], &item); err != nil {
		return nil, fmt.Errorf("media ddb get unmarshal: %w", err)
	}
	return mediaFromItem(&item), nil
}

func (r *DDBRepo) ListByTrip(ctx context.Context, tripID, dayID string) ([]*journeyv1.Media, error) {
	keyCond := expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value(ddb.TripKey(tripID))),
		expression.Key("SK").BeginsWith("MEDIA#"),
	)
	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	if dayID != "" {
		builder = builder.WithFilter(expression.Name("dayId").Equal(expression.Value(dayID)))
	}
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("media ddb list build: %w", err)
	}
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(false),
	})
	if err != nil {
		return nil, fmt.Errorf("media ddb list query: %w", err)
	}
	result := make([]*journeyv1.Media, 0, len(out.Items))
	for _, raw := range out.Items {
		var item mediaItem
		if err := attributevalue.UnmarshalMap(raw, &item); err != nil {
			return nil, fmt.Errorf("media ddb list unmarshal: %w", err)
		}
		result = append(result, mediaFromItem(&item))
	}
	return result, nil
}

func (r *DDBRepo) Delete(ctx context.Context, id string) error {
	m, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	takenAt := ""
	if m.GetTakenAt() != nil {
		takenAt = m.GetTakenAt().AsTime().Format(time.RFC3339)
	}
	sk := ddb.MediaKey(takenAt, id)
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: ddb.TripKey(m.GetTripId())},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return fmt.Errorf("media ddb delete: %w", err)
	}
	return nil
}

// DeleteByTrip removes every media row for the trip using batched writes.
// The SK is rebuilt from each item's takenAt so we avoid re-reading individually.
func (r *DDBRepo) DeleteByTrip(ctx context.Context, tripID string) error {
	items, err := r.ListByTrip(ctx, tripID, "")
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	sks := make([]string, 0, len(items))
	for _, m := range items {
		takenAt := ""
		if m.GetTakenAt() != nil {
			takenAt = m.GetTakenAt().AsTime().Format(time.RFC3339)
		}
		sks = append(sks, ddb.MediaKey(takenAt, m.GetId()))
	}
	return ddb.BatchDelete(ctx, r.client, r.table, ddb.BatchDeletePK(ddb.TripKey(tripID), sks))
}

func (r *DDBRepo) CountByTrip(ctx context.Context, tripID string) (int, error) {
	keyCond := expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value(ddb.TripKey(tripID))),
		expression.Key("SK").BeginsWith("MEDIA#"),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return 0, fmt.Errorf("media ddb count build: %w", err)
	}
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Select:                    "COUNT",
	})
	if err != nil {
		return 0, fmt.Errorf("media ddb count: %w", err)
	}
	return int(out.Count), nil
}

func mediaToItem(m *journeyv1.Media) *mediaItem {
	takenAt := ""
	if m.GetTakenAt() != nil {
		takenAt = m.GetTakenAt().AsTime().Format(time.RFC3339)
	}
	item := &mediaItem{
		PK:             ddb.TripKey(m.GetTripId()),
		SK:             ddb.MediaKey(takenAt, m.GetId()),
		ID:             m.GetId(),
		TripID:         m.GetTripId(),
		DayID:          m.GetDayId(),
		ScheduleID:     m.GetScheduleId(),
		S3Key:          m.GetS3Key(),
		ThumbnailS3Key: m.GetThumbnailS3Key(),
		MimeType:       m.GetMimeType(),
		Size:           m.GetSize(),
		TakenAt:        takenAt,
		Caption:        m.GetCaption(),
	}
	if m.GetLocation() != nil {
		item.Latitude = m.GetLocation().GetLatitude()
		item.Longitude = m.GetLocation().GetLongitude()
	}
	if m.GetAudit() != nil {
		item.CreatedAt = m.GetAudit().GetCreatedAt().AsTime().Format(time.RFC3339)
		item.UpdatedAt = m.GetAudit().GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return item
}

func mediaFromItem(item *mediaItem) *journeyv1.Media {
	m := &journeyv1.Media{
		Id:             item.ID,
		TripId:         item.TripID,
		DayId:          item.DayID,
		ScheduleId:     item.ScheduleID,
		S3Key:          item.S3Key,
		ThumbnailS3Key: item.ThumbnailS3Key,
		MimeType:       item.MimeType,
		Size:           item.Size,
		Caption:        item.Caption,
	}
	if item.TakenAt != "" {
		if t, err := time.Parse(time.RFC3339, item.TakenAt); err == nil {
			m.TakenAt = timestamppb.New(t)
		}
	}
	if item.Latitude != 0 || item.Longitude != 0 {
		m.Location = &journeyv1.GeoPoint{Latitude: item.Latitude, Longitude: item.Longitude}
	}
	if item.CreatedAt != "" {
		if ct, err := time.Parse(time.RFC3339, item.CreatedAt); err == nil {
			m.Audit = &journeyv1.AuditTimestamps{CreatedAt: timestamppb.New(ct)}
		}
	}
	if m.Audit != nil && item.UpdatedAt != "" {
		if ut, err := time.Parse(time.RFC3339, item.UpdatedAt); err == nil {
			m.Audit.UpdatedAt = timestamppb.New(ut)
		}
	}
	return m
}

// ListByTripPage is the cursor-aware paged variant of ListByTrip.
func (r *DDBRepo) ListByTripPage(ctx context.Context, tripID, dayID, cursor string, limit int32) ([]*journeyv1.Media, string, error) {
	keyCond := expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value(ddb.TripKey(tripID))),
		expression.Key("SK").BeginsWith("MEDIA#"),
	)
	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	if dayID != "" {
		builder = builder.WithFilter(expression.Name("dayId").Equal(expression.Value(dayID)))
	}
	expr, err := builder.Build()
	if err != nil {
		return nil, "", fmt.Errorf("media ddb list page build: %w", err)
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	start, err := ddb.DecodeCursor(cursor)
	if err != nil {
		return nil, "", err
	}
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(false),
		Limit:                     aws.Int32(limit),
		ExclusiveStartKey:         start,
	})
	if err != nil {
		return nil, "", fmt.Errorf("media ddb list page query: %w", err)
	}
	result := make([]*journeyv1.Media, 0, len(out.Items))
	for _, raw := range out.Items {
		var item mediaItem
		if err := attributevalue.UnmarshalMap(raw, &item); err != nil {
			return nil, "", fmt.Errorf("media ddb list page unmarshal: %w", err)
		}
		result = append(result, mediaFromItem(&item))
	}
	next, err := ddb.EncodeCursor(out.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}
	return result, next, nil
}
