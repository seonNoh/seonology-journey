package day

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/seonNoh/seonology-journey/apps/back/internal/repository/ddb"
)

const daysTable = "days"

// dayIDIndex is the GSI provisioned on the days table with `id` as HASH.
// It allows O(1) lookup by dayID without requiring the owning tripID.
const dayIDIndex = "DayIdIndex"

type dayItem struct {
	PK           string `dynamodbav:"PK"`
	SK           string `dynamodbav:"SK"`
	ID           string `dynamodbav:"id"`
	TripID       string `dynamodbav:"tripId"`
	DayNumber    int32  `dynamodbav:"dayNumber"`
	Date         string `dynamodbav:"date"`
	DayOfWeek    string `dynamodbav:"dayOfWeek,omitempty"`
	Region       string `dynamodbav:"region,omitempty"`
	Weather      string `dynamodbav:"weather,omitempty"`
	DailySummary string `dynamodbav:"dailySummary,omitempty"`
	CreatedAt    string `dynamodbav:"createdAt"`
	UpdatedAt    string `dynamodbav:"updatedAt"`
}

// DDBRepo is the DynamoDB implementation of day.Repository.
type DDBRepo struct {
	client *dynamodb.Client
	table  string
}

// NewDDBRepo creates a DDB-backed Day repository.
func NewDDBRepo(client *dynamodb.Client) *DDBRepo {
	return &DDBRepo{
		client: client,
		table:  ddb.TableName(daysTable),
	}
}

func (r *DDBRepo) Create(ctx context.Context, d *journeyv1.Day) error {
	item := dayToItem(d)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("day ddb create marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.table),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("day ddb create: %w", err)
	}
	return nil
}

// Get queries the DayIdIndex GSI for O(1) lookup by dayID.
// Replaces the previous full-table Scan.
func (r *DDBRepo) Get(ctx context.Context, id string) (*journeyv1.Day, error) {
	keyCond := expression.Key("id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("day ddb get build expr: %w", err)
	}
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		IndexName:                 aws.String(dayIDIndex),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("day ddb get query: %w", err)
	}
	if len(out.Items) == 0 {
		return nil, ErrNotFound
	}
	var item dayItem
	if err := attributevalue.UnmarshalMap(out.Items[0], &item); err != nil {
		return nil, fmt.Errorf("day ddb get unmarshal: %w", err)
	}
	return dayFromItem(&item), nil
}

func (r *DDBRepo) ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Day, error) {
	keyCond := expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value(ddb.TripKey(tripID))),
		expression.Key("SK").BeginsWith("DAY#"),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("day ddb list build expr: %w", err)
	}
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("day ddb list query: %w", err)
	}
	days := make([]*journeyv1.Day, 0, len(out.Items))
	for _, raw := range out.Items {
		var item dayItem
		if err := attributevalue.UnmarshalMap(raw, &item); err != nil {
			return nil, fmt.Errorf("day ddb list unmarshal: %w", err)
		}
		days = append(days, dayFromItem(&item))
	}
	return days, nil
}

func (r *DDBRepo) Update(ctx context.Context, d *journeyv1.Day) error {
	item := dayToItem(d)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("day ddb update marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.table),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("day ddb update: %w", err)
	}
	return nil
}

// DeleteByTrip removes every day for the trip in one BatchWriteItem call
// per 25-item chunk. Previous implementation issued one DeleteItem per day.
func (r *DDBRepo) DeleteByTrip(ctx context.Context, tripID string) error {
	days, err := r.ListByTrip(ctx, tripID)
	if err != nil {
		return err
	}
	if len(days) == 0 {
		return nil
	}
	sks := make([]string, 0, len(days))
	for _, d := range days {
		sks = append(sks, ddb.DayKey(d.GetId()))
	}
	return ddb.BatchDelete(ctx, r.client, r.table, ddb.BatchDeletePK(ddb.TripKey(tripID), sks))
}

func dayToItem(d *journeyv1.Day) *dayItem {
	item := &dayItem{
		PK:           ddb.TripKey(d.GetTripId()),
		SK:           ddb.DayKey(d.GetId()),
		ID:           d.GetId(),
		TripID:       d.GetTripId(),
		DayNumber:    d.GetDayNumber(),
		Date:         d.GetDate(),
		DayOfWeek:    d.GetDayOfWeek(),
		Region:       d.GetRegion(),
		Weather:      d.GetWeather(),
		DailySummary: d.GetDailySummary(),
	}
	if d.GetAudit() != nil {
		item.CreatedAt = d.GetAudit().GetCreatedAt().AsTime().Format(time.RFC3339)
		item.UpdatedAt = d.GetAudit().GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return item
}

func dayFromItem(item *dayItem) *journeyv1.Day {
	d := &journeyv1.Day{
		Id:           item.ID,
		TripId:       item.TripID,
		DayNumber:    item.DayNumber,
		Date:         item.Date,
		DayOfWeek:    item.DayOfWeek,
		Region:       item.Region,
		Weather:      item.Weather,
		DailySummary: item.DailySummary,
	}
	if item.CreatedAt != "" {
		if ct, err := time.Parse(time.RFC3339, item.CreatedAt); err == nil {
			d.Audit = &journeyv1.AuditTimestamps{CreatedAt: timestamppb.New(ct)}
		}
	}
	if d.Audit != nil && item.UpdatedAt != "" {
		if ut, err := time.Parse(time.RFC3339, item.UpdatedAt); err == nil {
			d.Audit.UpdatedAt = timestamppb.New(ut)
		}
	}
	return d
}
