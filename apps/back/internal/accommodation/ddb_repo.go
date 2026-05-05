package accommodation

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/seonNoh/seonology-journey/apps/back/internal/repository/ddb"
)

const accommodationsTable = "accommodations"

type accomItem struct {
	PK            string  `dynamodbav:"PK"`
	SK            string  `dynamodbav:"SK"`
	DayID         string  `dynamodbav:"dayId"`
	Name          string  `dynamodbav:"name"`
	CheckInTime   string  `dynamodbav:"checkInTime,omitempty"`
	CheckOutTime  string  `dynamodbav:"checkOutTime,omitempty"`
	CostAmount    int64   `dynamodbav:"costAmount,omitempty"`
	CostCurrency  string  `dynamodbav:"costCurrency,omitempty"`
	Amenities     string  `dynamodbav:"amenities,omitempty"`
	Address       string  `dynamodbav:"address,omitempty"`
	Latitude      float64 `dynamodbav:"latitude,omitempty"`
	Longitude     float64 `dynamodbav:"longitude,omitempty"`
	ReservationID string  `dynamodbav:"reservationId,omitempty"`
	CreatedAt     string  `dynamodbav:"createdAt"`
	UpdatedAt     string  `dynamodbav:"updatedAt"`
}

// DDBRepo is the DynamoDB implementation of accommodation.Repository.
type DDBRepo struct {
	client *dynamodb.Client
	table  string
}

func NewDDBRepo(client *dynamodb.Client) *DDBRepo {
	return &DDBRepo{client: client, table: ddb.TableName(accommodationsTable)}
}

func (r *DDBRepo) Upsert(ctx context.Context, a *journeyv1.Accommodation) error {
	item := accomToItem(a)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("accommodation ddb upsert marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("accommodation ddb upsert: %w", err)
	}
	return nil
}

func (r *DDBRepo) Get(ctx context.Context, dayID string) (*journeyv1.Accommodation, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: ddb.DayKey(dayID)},
			"SK": &types.AttributeValueMemberS{Value: ddb.AccommodationSK},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("accommodation ddb get: %w", err)
	}
	if out.Item == nil {
		return nil, ErrNotFound
	}
	var item accomItem
	if err := attributevalue.UnmarshalMap(out.Item, &item); err != nil {
		return nil, fmt.Errorf("accommodation ddb get unmarshal: %w", err)
	}
	return accomFromItem(&item), nil
}

func (r *DDBRepo) Delete(ctx context.Context, dayID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: ddb.DayKey(dayID)},
			"SK": &types.AttributeValueMemberS{Value: ddb.AccommodationSK},
		},
	})
	if err != nil {
		return fmt.Errorf("accommodation ddb delete: %w", err)
	}
	return nil
}

func accomToItem(a *journeyv1.Accommodation) *accomItem {
	item := &accomItem{
		PK:            ddb.DayKey(a.GetDayId()),
		SK:            ddb.AccommodationSK,
		DayID:         a.GetDayId(),
		Name:          a.GetName(),
		CheckInTime:   a.GetCheckInTime(),
		CheckOutTime:  a.GetCheckOutTime(),
		Amenities:     a.GetAmenities(),
		Address:       a.GetAddress(),
		ReservationID: a.GetReservationId(),
	}
	if a.GetCost() != nil {
		item.CostAmount = a.GetCost().GetAmount()
		item.CostCurrency = a.GetCost().GetCurrency()
	}
	if a.GetLocation() != nil {
		item.Latitude = a.GetLocation().GetLatitude()
		item.Longitude = a.GetLocation().GetLongitude()
	}
	if a.GetAudit() != nil {
		item.CreatedAt = a.GetAudit().GetCreatedAt().AsTime().Format(time.RFC3339)
		item.UpdatedAt = a.GetAudit().GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return item
}

func accomFromItem(item *accomItem) *journeyv1.Accommodation {
	a := &journeyv1.Accommodation{
		DayId:         item.DayID,
		Name:          item.Name,
		CheckInTime:   item.CheckInTime,
		CheckOutTime:  item.CheckOutTime,
		Amenities:     item.Amenities,
		Address:       item.Address,
		ReservationId: item.ReservationID,
	}
	if item.CostAmount != 0 || item.CostCurrency != "" {
		a.Cost = &journeyv1.Money{Amount: item.CostAmount, Currency: item.CostCurrency}
	}
	if item.Latitude != 0 || item.Longitude != 0 {
		a.Location = &journeyv1.GeoPoint{Latitude: item.Latitude, Longitude: item.Longitude}
	}
	if item.CreatedAt != "" {
		if ct, err := time.Parse(time.RFC3339, item.CreatedAt); err == nil {
			a.Audit = &journeyv1.AuditTimestamps{CreatedAt: timestamppb.New(ct)}
		}
	}
	if a.Audit != nil && item.UpdatedAt != "" {
		if ut, err := time.Parse(time.RFC3339, item.UpdatedAt); err == nil {
			a.Audit.UpdatedAt = timestamppb.New(ut)
		}
	}
	return a
}
