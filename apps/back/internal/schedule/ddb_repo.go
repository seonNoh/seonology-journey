package schedule

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

const schedulesTable = "schedules"

type scheduleItem struct {
	PK              string  `dynamodbav:"PK"`
	SK              string  `dynamodbav:"SK"`
	ID              string  `dynamodbav:"id"`
	DayID           string  `dynamodbav:"dayId"`
	Order           int32   `dynamodbav:"order"`
	StartTime       string  `dynamodbav:"startTime,omitempty"`
	EndTime         string  `dynamodbav:"endTime,omitempty"`
	Title           string  `dynamodbav:"title"`
	Region          string  `dynamodbav:"region,omitempty"`
	Category        int32   `dynamodbav:"category"`
	Transport       int32   `dynamodbav:"transport,omitempty"`
	TransportDetail string  `dynamodbav:"transportDetail,omitempty"`
	CostAmount      int64   `dynamodbav:"costAmount,omitempty"`
	CostCurrency    string  `dynamodbav:"costCurrency,omitempty"`
	PlaceName       string  `dynamodbav:"placeName,omitempty"`
	Latitude        float64 `dynamodbav:"latitude,omitempty"`
	Longitude       float64 `dynamodbav:"longitude,omitempty"`
	Notes           string  `dynamodbav:"notes,omitempty"`
	IsCompleted     bool    `dynamodbav:"isCompleted"`
	CreatedAt       string  `dynamodbav:"createdAt"`
	UpdatedAt       string  `dynamodbav:"updatedAt"`
}

// DDBRepo is the DynamoDB implementation of schedule.Repository.
type DDBRepo struct {
	client *dynamodb.Client
	table  string
}

// NewDDBRepo creates a DDB-backed Schedule repository.
func NewDDBRepo(client *dynamodb.Client) *DDBRepo {
	return &DDBRepo{
		client: client,
		table:  ddb.TableName(schedulesTable),
	}
}

func (r *DDBRepo) Create(ctx context.Context, s *journeyv1.Schedule) error {
	item := scheduleToItem(s)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("schedule ddb create marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.table),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("schedule ddb create: %w", err)
	}
	return nil
}

func (r *DDBRepo) Get(ctx context.Context, id string) (*journeyv1.Schedule, error) {
	filt := expression.Name("id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, fmt.Errorf("schedule ddb get build: %w", err)
	}
	out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 aws.String(r.table),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("schedule ddb get scan: %w", err)
	}
	if len(out.Items) == 0 {
		return nil, ErrNotFound
	}
	var item scheduleItem
	if err := attributevalue.UnmarshalMap(out.Items[0], &item); err != nil {
		return nil, fmt.Errorf("schedule ddb get unmarshal: %w", err)
	}
	return scheduleFromItem(&item), nil
}

func (r *DDBRepo) ListByDay(ctx context.Context, dayID string) ([]*journeyv1.Schedule, error) {
	keyCond := expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value(ddb.DayKey(dayID))),
		expression.Key("SK").BeginsWith("SCHEDULE#"),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("schedule ddb list build: %w", err)
	}
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("schedule ddb list query: %w", err)
	}
	schedules := make([]*journeyv1.Schedule, 0, len(out.Items))
	for _, raw := range out.Items {
		var item scheduleItem
		if err := attributevalue.UnmarshalMap(raw, &item); err != nil {
			return nil, fmt.Errorf("schedule ddb list unmarshal: %w", err)
		}
		schedules = append(schedules, scheduleFromItem(&item))
	}
	return schedules, nil
}

func (r *DDBRepo) Update(ctx context.Context, s *journeyv1.Schedule) error {
	item := scheduleToItem(s)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("schedule ddb update marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.table),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("schedule ddb update: %w", err)
	}
	return nil
}

func (r *DDBRepo) Delete(ctx context.Context, id string) error {
	s, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: ddb.DayKey(s.GetDayId())},
			"SK": &types.AttributeValueMemberS{Value: ddb.ScheduleKey(id)},
		},
	})
	if err != nil {
		return fmt.Errorf("schedule ddb delete: %w", err)
	}
	return nil
}

// DeleteByDay deletes all schedules for a given day.
func (r *DDBRepo) DeleteByDay(ctx context.Context, dayID string) error {
	schedules, err := r.ListByDay(ctx, dayID)
	if err != nil {
		return err
	}
	for _, s := range schedules {
		_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(r.table),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: ddb.DayKey(dayID)},
				"SK": &types.AttributeValueMemberS{Value: ddb.ScheduleKey(s.GetId())},
			},
		})
		if err != nil {
			return fmt.Errorf("schedule ddb deleteByDay %s: %w", s.GetId(), err)
		}
	}
	return nil
}

func scheduleToItem(s *journeyv1.Schedule) *scheduleItem {
	item := &scheduleItem{
		PK:              ddb.DayKey(s.GetDayId()),
		SK:              ddb.ScheduleKey(s.GetId()),
		ID:              s.GetId(),
		DayID:           s.GetDayId(),
		Order:           s.GetOrder(),
		StartTime:       s.GetStartTime(),
		EndTime:         s.GetEndTime(),
		Title:           s.GetTitle(),
		Region:          s.GetRegion(),
		Category:        int32(s.GetCategory()),
		Transport:       int32(s.GetTransport()),
		TransportDetail: s.GetTransportDetail(),
		PlaceName:       s.GetPlaceName(),
		Notes:           s.GetNotes(),
		IsCompleted:     s.GetIsCompleted(),
	}
	if s.GetCost() != nil {
		item.CostAmount = s.GetCost().GetAmount()
		item.CostCurrency = s.GetCost().GetCurrency()
	}
	if s.GetLocation() != nil {
		item.Latitude = s.GetLocation().GetLatitude()
		item.Longitude = s.GetLocation().GetLongitude()
	}
	if s.GetAudit() != nil {
		item.CreatedAt = s.GetAudit().GetCreatedAt().AsTime().Format(time.RFC3339)
		item.UpdatedAt = s.GetAudit().GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return item
}

func scheduleFromItem(item *scheduleItem) *journeyv1.Schedule {
	s := &journeyv1.Schedule{
		Id:              item.ID,
		DayId:           item.DayID,
		Order:           item.Order,
		StartTime:       item.StartTime,
		EndTime:         item.EndTime,
		Title:           item.Title,
		Region:          item.Region,
		Category:        journeyv1.ScheduleCategory(item.Category),
		Transport:       journeyv1.TransportType(item.Transport),
		TransportDetail: item.TransportDetail,
		PlaceName:       item.PlaceName,
		Notes:           item.Notes,
		IsCompleted:     item.IsCompleted,
	}
	if item.CostAmount != 0 || item.CostCurrency != "" {
		s.Cost = &journeyv1.Money{Amount: item.CostAmount, Currency: item.CostCurrency}
	}
	if item.Latitude != 0 || item.Longitude != 0 {
		s.Location = &journeyv1.GeoPoint{Latitude: item.Latitude, Longitude: item.Longitude}
	}
	if item.CreatedAt != "" {
		if ct, err := time.Parse(time.RFC3339, item.CreatedAt); err == nil {
			s.Audit = &journeyv1.AuditTimestamps{CreatedAt: timestamppb.New(ct)}
		}
	}
	if s.Audit != nil && item.UpdatedAt != "" {
		if ut, err := time.Parse(time.RFC3339, item.UpdatedAt); err == nil {
			s.Audit.UpdatedAt = timestamppb.New(ut)
		}
	}
	return s
}
