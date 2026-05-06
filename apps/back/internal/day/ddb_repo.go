package day

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

const daysTable = "days"

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

func (r *DDBRepo) Get(ctx context.Context, id string) (*journeyv1.Day, error) {
	// Day ID alone doesn't give us the PK (tripID), so we scan with filter.
	// For production with large data, consider a GSI on day ID.
	// For now, iterate all items. Alternatively, encode tripID in day ID.
	// Use a simple approach: scan with filter (acceptable for small scale).
	filt := expression.Name("id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, fmt.Errorf("day ddb get build expr: %w", err)
	}
	out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 aws.String(r.table),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("day ddb get scan: %w", err)
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

func (r *DDBRepo) DeleteByTrip(ctx context.Context, tripID string) error {
	days, err := r.ListByTrip(ctx, tripID)
	if err != nil {
		return err
	}
	for _, d := range days {
		_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(r.table),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: ddb.TripKey(tripID)},
				"SK": &types.AttributeValueMemberS{Value: ddb.DayKey(d.GetId())},
			},
		})
		if err != nil {
			return fmt.Errorf("day ddb delete item %s: %w", d.GetId(), err)
		}
	}
	return nil
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
