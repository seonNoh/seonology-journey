package meal

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

const mealsTable = "meals"

type mealItem struct {
	PK             string  `dynamodbav:"PK"`
	SK             string  `dynamodbav:"SK"`
	DayID          string  `dynamodbav:"dayId"`
	MealType       int32   `dynamodbav:"mealType"`
	Source         int32   `dynamodbav:"source"`
	RestaurantName string  `dynamodbav:"restaurantName,omitempty"`
	Menu           string  `dynamodbav:"menu,omitempty"`
	CostAmount     int64   `dynamodbav:"costAmount,omitempty"`
	CostCurrency   string  `dynamodbav:"costCurrency,omitempty"`
	Rating         int32   `dynamodbav:"rating,omitempty"`
	Review         string  `dynamodbav:"review,omitempty"`
	Latitude       float64 `dynamodbav:"latitude,omitempty"`
	Longitude      float64 `dynamodbav:"longitude,omitempty"`
	CreatedAt      string  `dynamodbav:"createdAt"`
	UpdatedAt      string  `dynamodbav:"updatedAt"`
}

// DDBRepo is the DynamoDB implementation of meal.Repository.
type DDBRepo struct {
	client *dynamodb.Client
	table  string
}

func NewDDBRepo(client *dynamodb.Client) *DDBRepo {
	return &DDBRepo{client: client, table: ddb.TableName(mealsTable)}
}

func (r *DDBRepo) Upsert(ctx context.Context, m *journeyv1.Meal) error {
	item := mealToItem(m)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("meal ddb upsert marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("meal ddb upsert: %w", err)
	}
	return nil
}

func (r *DDBRepo) ListByDay(ctx context.Context, dayID string) ([]*journeyv1.Meal, error) {
	keyCond := expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value(ddb.DayKey(dayID))),
		expression.Key("SK").BeginsWith("MEAL#"),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("meal ddb list build: %w", err)
	}
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("meal ddb list query: %w", err)
	}
	meals := make([]*journeyv1.Meal, 0, len(out.Items))
	for _, raw := range out.Items {
		var item mealItem
		if err := attributevalue.UnmarshalMap(raw, &item); err != nil {
			return nil, fmt.Errorf("meal ddb list unmarshal: %w", err)
		}
		meals = append(meals, mealFromItem(&item))
	}
	return meals, nil
}

func (r *DDBRepo) Delete(ctx context.Context, dayID string, mt journeyv1.MealType) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: ddb.DayKey(dayID)},
			"SK": &types.AttributeValueMemberS{Value: ddb.MealKey(mt.String())},
		},
	})
	if err != nil {
		return fmt.Errorf("meal ddb delete: %w", err)
	}
	return nil
}

// DeleteByDay removes all meal rows for the day in one BatchWriteItem sweep.
func (r *DDBRepo) DeleteByDay(ctx context.Context, dayID string) error {
	meals, err := r.ListByDay(ctx, dayID)
	if err != nil {
		return err
	}
	if len(meals) == 0 {
		return nil
	}
	sks := make([]string, 0, len(meals))
	for _, m := range meals {
		sks = append(sks, ddb.MealKey(m.GetMealType().String()))
	}
	return ddb.BatchDelete(ctx, r.client, r.table, ddb.BatchDeletePK(ddb.DayKey(dayID), sks))
}

func mealToItem(m *journeyv1.Meal) *mealItem {
	item := &mealItem{
		PK:             ddb.DayKey(m.GetDayId()),
		SK:             ddb.MealKey(m.GetMealType().String()),
		DayID:          m.GetDayId(),
		MealType:       int32(m.GetMealType()),
		Source:         int32(m.GetSource()),
		RestaurantName: m.GetRestaurantName(),
		Menu:           m.GetMenu(),
		Rating:         m.GetRating(),
		Review:         m.GetReview(),
	}
	if m.GetCost() != nil {
		item.CostAmount = m.GetCost().GetAmount()
		item.CostCurrency = m.GetCost().GetCurrency()
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

func mealFromItem(item *mealItem) *journeyv1.Meal {
	m := &journeyv1.Meal{
		DayId:          item.DayID,
		MealType:       journeyv1.MealType(item.MealType),
		Source:         journeyv1.MealSource(item.Source),
		RestaurantName: item.RestaurantName,
		Menu:           item.Menu,
		Rating:         item.Rating,
		Review:         item.Review,
	}
	if item.CostAmount != 0 || item.CostCurrency != "" {
		m.Cost = &journeyv1.Money{Amount: item.CostAmount, Currency: item.CostCurrency}
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
