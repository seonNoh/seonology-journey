package trip

import (
	"context"
	"errors"
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

const tripsTable = "trips"

// tripItem is the DynamoDB item representation of a Trip.
//
// Keys:
//
//	PK = USER#{ownerID}
//	SK = TRIP#{tripID}
//	GSI1PK = TRIP#{tripID}     (direct lookup by tripID)
//	GSI1SK = METADATA
type tripItem struct {
	PK          string `dynamodbav:"PK"`
	SK          string `dynamodbav:"SK"`
	GSI1PK      string `dynamodbav:"GSI1PK"`
	GSI1SK      string `dynamodbav:"GSI1SK"`
	ID          string `dynamodbav:"id"`
	OwnerID     string `dynamodbav:"ownerId"`
	Title       string `dynamodbav:"title"`
	Description string `dynamodbav:"description,omitempty"`
	StartDate   string `dynamodbav:"startDate"`
	EndDate     string `dynamodbav:"endDate"`
	Status      int32  `dynamodbav:"status"`
	Destination string `dynamodbav:"destination,omitempty"`
	CountryCode string `dynamodbav:"countryCode,omitempty"`
	CoverImage  string `dynamodbav:"coverImageUrl,omitempty"`
	BudgetAmt   int64  `dynamodbav:"budgetAmount,omitempty"`
	BudgetCur   string `dynamodbav:"budgetCurrency,omitempty"`
	CreatedAt   string `dynamodbav:"createdAt"`
	UpdatedAt   string `dynamodbav:"updatedAt"`
}

// DDBRepo is the DynamoDB implementation of Repository.
type DDBRepo struct {
	client *dynamodb.Client
	table  string
}

// NewDDBRepo creates a new DDB-backed Trip repository.
func NewDDBRepo(client *dynamodb.Client) *DDBRepo {
	return &DDBRepo{
		client: client,
		table:  ddb.TableName(tripsTable),
	}
}

func (r *DDBRepo) Create(ctx context.Context, t *journeyv1.Trip) error {
	item := toItem(t)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("trip ddb create marshal: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.table),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return fmt.Errorf("trip already exists: %s", t.GetId())
		}
		return fmt.Errorf("trip ddb create: %w", err)
	}
	return nil
}

// Get looks up a trip by ID using GSI1 (GSI1PK = TRIP#{id}).
// Each lookup consumes 1 RCU for the GSI plus 1 RCU for the base item since
// GSI1 projects ALL attributes; this is effectively a single-shot read.
func (r *DDBRepo) Get(ctx context.Context, id string) (*journeyv1.Trip, error) {
	keyCond := expression.KeyAnd(
		expression.Key("GSI1PK").Equal(expression.Value(ddb.TripKey(id))),
		expression.Key("GSI1SK").Equal(expression.Value("METADATA")),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("trip ddb get build expr: %w", err)
	}

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("trip ddb get query: %w", err)
	}
	if len(out.Items) == 0 {
		return nil, ErrNotFound
	}

	var item tripItem
	if err := attributevalue.UnmarshalMap(out.Items[0], &item); err != nil {
		return nil, fmt.Errorf("trip ddb get unmarshal: %w", err)
	}
	return fromItem(&item), nil
}

// ListByOwner returns trips belonging to ownerID using the base table partition.
// This replaces the previous full-table Scan: each user's trips live under a
// single partition so Query on PK is O(N_user), not O(N_total).
func (r *DDBRepo) ListByOwner(ctx context.Context, ownerID string, status journeyv1.TripStatus) ([]*journeyv1.Trip, error) {
	keyCond := expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value(ddb.UserKey(ownerID))),
		expression.Key("SK").BeginsWith("TRIP#"),
	)
	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	if status != journeyv1.TripStatus_TRIP_STATUS_UNSPECIFIED {
		builder = builder.WithFilter(expression.Name("status").Equal(expression.Value(int32(status))))
	}
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("trip ddb list build expr: %w", err)
	}

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("trip ddb list query: %w", err)
	}

	trips := make([]*journeyv1.Trip, 0, len(out.Items))
	for _, raw := range out.Items {
		var item tripItem
		if err := attributevalue.UnmarshalMap(raw, &item); err != nil {
			continue
		}
		trips = append(trips, fromItem(&item))
	}
	return trips, nil
}

func (r *DDBRepo) Update(ctx context.Context, t *journeyv1.Trip) error {
	item := toItem(t)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("trip ddb update marshal: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.table),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return ErrNotFound
		}
		return fmt.Errorf("trip ddb update: %w", err)
	}
	return nil
}

// Delete does a GSI1 lookup to resolve ownerID, then issues a single
// DeleteItem against the base table. Two operations instead of one,
// but still O(1).
func (r *DDBRepo) Delete(ctx context.Context, id string) error {
	t, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: ddb.UserKey(t.GetOwnerId())},
			"SK": &types.AttributeValueMemberS{Value: ddb.TripKey(id)},
		},
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return ErrNotFound
		}
		return fmt.Errorf("trip ddb delete: %w", err)
	}
	return nil
}

// toItem converts a proto Trip to a DynamoDB item.
func toItem(t *journeyv1.Trip) *tripItem {
	item := &tripItem{
		PK:          ddb.UserKey(t.GetOwnerId()),
		SK:          ddb.TripKey(t.GetId()),
		GSI1PK:      ddb.TripKey(t.GetId()),
		GSI1SK:      "METADATA",
		ID:          t.GetId(),
		OwnerID:     t.GetOwnerId(),
		Title:       t.GetTitle(),
		Description: t.GetDescription(),
		StartDate:   t.GetStartDate(),
		EndDate:     t.GetEndDate(),
		Status:      int32(t.GetStatus()),
		Destination: t.GetDestination(),
		CountryCode: t.GetCountryCode(),
		CoverImage:  t.GetCoverImageUrl(),
	}

	if t.GetTotalBudget() != nil {
		item.BudgetAmt = t.GetTotalBudget().GetAmount()
		item.BudgetCur = t.GetTotalBudget().GetCurrency()
	}

	if t.GetAudit() != nil {
		item.CreatedAt = t.GetAudit().GetCreatedAt().AsTime().Format("2006-01-02T15:04:05Z")
		item.UpdatedAt = t.GetAudit().GetUpdatedAt().AsTime().Format("2006-01-02T15:04:05Z")
	}

	return item
}

// fromItem converts a DynamoDB item to a proto Trip.
func fromItem(item *tripItem) *journeyv1.Trip {
	t := &journeyv1.Trip{
		Id:            item.ID,
		OwnerId:       item.OwnerID,
		Title:         item.Title,
		Description:   item.Description,
		StartDate:     item.StartDate,
		EndDate:       item.EndDate,
		Status:        journeyv1.TripStatus(item.Status),
		Destination:   item.Destination,
		CountryCode:   item.CountryCode,
		CoverImageUrl: item.CoverImage,
	}

	if item.BudgetAmt != 0 || item.BudgetCur != "" {
		t.TotalBudget = &journeyv1.Money{
			Amount:   item.BudgetAmt,
			Currency: item.BudgetCur,
		}
	}

	if item.CreatedAt != "" {
		if ct, err := parseTime(item.CreatedAt); err == nil {
			t.Audit = &journeyv1.AuditTimestamps{CreatedAt: ct}
		}
	}
	if t.Audit != nil && item.UpdatedAt != "" {
		if ut, err := parseTime(item.UpdatedAt); err == nil {
			t.Audit.UpdatedAt = ut
		}
	}

	return t
}

func parseTime(s string) (*timestamppb.Timestamp, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return timestamppb.New(t), nil
}

// ListByOwnerPage is the cursor-aware variant. Uses the same base partition
// as ListByOwner and emits a base64 cursor derived from DynamoDB's
// LastEvaluatedKey so callers can resume without reading prior pages.
func (r *DDBRepo) ListByOwnerPage(ctx context.Context, ownerID string, status journeyv1.TripStatus, cursor string, limit int32) ([]*journeyv1.Trip, string, error) {
	keyCond := expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value(ddb.UserKey(ownerID))),
		expression.Key("SK").BeginsWith("TRIP#"),
	)
	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	if status != journeyv1.TripStatus_TRIP_STATUS_UNSPECIFIED {
		builder = builder.WithFilter(expression.Name("status").Equal(expression.Value(int32(status))))
	}
	expr, err := builder.Build()
	if err != nil {
		return nil, "", fmt.Errorf("trip ddb list page build expr: %w", err)
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
		Limit:                     aws.Int32(limit),
		ExclusiveStartKey:         start,
	})
	if err != nil {
		return nil, "", fmt.Errorf("trip ddb list page query: %w", err)
	}

	trips := make([]*journeyv1.Trip, 0, len(out.Items))
	for _, raw := range out.Items {
		var item tripItem
		if err := attributevalue.UnmarshalMap(raw, &item); err != nil {
			continue
		}
		trips = append(trips, fromItem(&item))
	}
	next, err := ddb.EncodeCursor(out.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}
	return trips, next, nil
}
