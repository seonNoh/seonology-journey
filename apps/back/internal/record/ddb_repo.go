package record

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

// idIndex is the GSI name used by the expenses / notes / checklists /
// reservations tables for single-attribute `id` HASH lookups. All four
// tables were provisioned with the same index name, which lets us share
// the small query helper below.
const idIndex = "IdIndex"

// queryByIDOnIndex does a single-item Query against the IdIndex GSI.
func queryByIDOnIndex(ctx context.Context, client *dynamodb.Client, table, id string) (map[string]types.AttributeValue, error) {
	keyCond := expression.Key("id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("%s id index build: %w", table, err)
	}
	out, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(table),
		IndexName:                 aws.String(idIndex),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("%s id index query: %w", table, err)
	}
	if len(out.Items) == 0 {
		return nil, ErrNotFound
	}
	return out.Items[0], nil
}

// ─── Expense DDB ─────────────────────────────────────────────────────────────

type expenseItem struct {
	PK            string `dynamodbav:"PK"` // TRIP#<tripID>
	SK            string `dynamodbav:"SK"` // EXPENSE#<id>
	ID            string `dynamodbav:"id"`
	TripID        string `dynamodbav:"tripId"`
	DayID         string `dynamodbav:"dayId,omitempty"`
	Category      int32  `dynamodbav:"category"`
	AmountVal     int64  `dynamodbav:"amountVal"`
	AmountCur     string `dynamodbav:"amountCur"`
	PaymentMethod int32  `dynamodbav:"paymentMethod"`
	Description   string `dynamodbav:"description,omitempty"`
	SpentAt       string `dynamodbav:"spentAt,omitempty"`
	CreatedAt     string `dynamodbav:"createdAt"`
	UpdatedAt     string `dynamodbav:"updatedAt"`
}

func expenseToItem(e *journeyv1.Expense) expenseItem {
	it := expenseItem{
		PK:            "TRIP#" + e.GetTripId(),
		SK:            "EXPENSE#" + e.GetId(),
		ID:            e.GetId(),
		TripID:        e.GetTripId(),
		DayID:         e.GetDayId(),
		Category:      int32(e.GetCategory()),
		PaymentMethod: int32(e.GetPaymentMethod()),
		Description:   e.GetDescription(),
	}
	if m := e.GetAmount(); m != nil {
		it.AmountVal = m.GetAmount()
		it.AmountCur = m.GetCurrency()
	}
	if t := e.GetSpentAt(); t != nil {
		it.SpentAt = t.AsTime().Format(time.RFC3339)
	}
	if a := e.GetAudit(); a != nil {
		it.CreatedAt = a.GetCreatedAt().AsTime().Format(time.RFC3339)
		it.UpdatedAt = a.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return it
}

func expenseFromItem(it expenseItem) *journeyv1.Expense {
	e := &journeyv1.Expense{
		Id:            it.ID,
		TripId:        it.TripID,
		DayId:         it.DayID,
		Category:      journeyv1.ExpenseCategory(it.Category),
		PaymentMethod: journeyv1.PaymentMethod(it.PaymentMethod),
		Description:   it.Description,
	}
	if it.AmountCur != "" || it.AmountVal != 0 {
		e.Amount = &journeyv1.Money{Currency: it.AmountCur, Amount: it.AmountVal}
	}
	if it.SpentAt != "" {
		if t, err := time.Parse(time.RFC3339, it.SpentAt); err == nil {
			e.SpentAt = timestamppb.New(t)
		}
	}
	ca, _ := time.Parse(time.RFC3339, it.CreatedAt)
	ua, _ := time.Parse(time.RFC3339, it.UpdatedAt)
	e.Audit = &journeyv1.AuditTimestamps{
		CreatedAt: timestamppb.New(ca),
		UpdatedAt: timestamppb.New(ua),
	}
	return e
}

// ExpenseDDBRepo implements ExpenseRepository backed by DynamoDB.
type ExpenseDDBRepo struct {
	client *dynamodb.Client
	table  string
}

func NewExpenseDDBRepo(client *dynamodb.Client) *ExpenseDDBRepo {
	return &ExpenseDDBRepo{client: client, table: ddb.TableName("expenses")}
}

func (r *ExpenseDDBRepo) Save(ctx context.Context, e *journeyv1.Expense) error {
	it := expenseToItem(e)
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return fmt.Errorf("expense marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      av,
	})
	return err
}

func (r *ExpenseDDBRepo) Get(ctx context.Context, id string) (*journeyv1.Expense, error) {
	raw, err := queryByIDOnIndex(ctx, r.client, r.table, id)
	if err != nil {
		return nil, err
	}
	var it expenseItem
	if err := attributevalue.UnmarshalMap(raw, &it); err != nil {
		return nil, err
	}
	return expenseFromItem(it), nil
}

func (r *ExpenseDDBRepo) Delete(ctx context.Context, id string) error {
	e, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "TRIP#" + e.GetTripId()},
			"SK": &types.AttributeValueMemberS{Value: "EXPENSE#" + id},
		},
	})
	return err
}

func (r *ExpenseDDBRepo) ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Expense, error) {
	keyCond := expression.Key("PK").Equal(expression.Value("TRIP#" + tripID)).
		And(expression.Key("SK").BeginsWith("EXPENSE#"))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("expense list: %w", err)
	}
	results := make([]*journeyv1.Expense, 0, len(out.Items))
	for _, item := range out.Items {
		var it expenseItem
		if err := attributevalue.UnmarshalMap(item, &it); err != nil {
			continue
		}
		results = append(results, expenseFromItem(it))
	}
	return results, nil
}

func (r *ExpenseDDBRepo) DeleteByTrip(ctx context.Context, tripID string) error {
	expenses, err := r.ListByTrip(ctx, tripID)
	if err != nil {
		return err
	}
	if len(expenses) == 0 {
		return nil
	}
	sks := make([]string, 0, len(expenses))
	for _, e := range expenses {
		sks = append(sks, "EXPENSE#"+e.GetId())
	}
	return ddb.BatchDelete(ctx, r.client, r.table, ddb.BatchDeletePK("TRIP#"+tripID, sks))
}

// ─── Note DDB ────────────────────────────────────────────────────────────────

type noteItem struct {
	PK        string `dynamodbav:"PK"` // TRIP#<tripID>
	SK        string `dynamodbav:"SK"` // NOTE#<id>
	ID        string `dynamodbav:"id"`
	TripID    string `dynamodbav:"tripId"`
	DayID     string `dynamodbav:"dayId,omitempty"`
	Content   string `dynamodbav:"content"`
	Mood      string `dynamodbav:"mood,omitempty"`
	CreatedAt string `dynamodbav:"createdAt"`
	UpdatedAt string `dynamodbav:"updatedAt"`
}

func noteToItem(n *journeyv1.Note) noteItem {
	it := noteItem{
		PK:      "TRIP#" + n.GetTripId(),
		SK:      "NOTE#" + n.GetId(),
		ID:      n.GetId(),
		TripID:  n.GetTripId(),
		DayID:   n.GetDayId(),
		Content: n.GetContent(),
		Mood:    n.GetMood(),
	}
	if a := n.GetAudit(); a != nil {
		it.CreatedAt = a.GetCreatedAt().AsTime().Format(time.RFC3339)
		it.UpdatedAt = a.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return it
}

func noteFromItem(it noteItem) *journeyv1.Note {
	n := &journeyv1.Note{
		Id:      it.ID,
		TripId:  it.TripID,
		DayId:   it.DayID,
		Content: it.Content,
		Mood:    it.Mood,
	}
	ca, _ := time.Parse(time.RFC3339, it.CreatedAt)
	ua, _ := time.Parse(time.RFC3339, it.UpdatedAt)
	n.Audit = &journeyv1.AuditTimestamps{
		CreatedAt: timestamppb.New(ca),
		UpdatedAt: timestamppb.New(ua),
	}
	return n
}

// NoteDDBRepo implements NoteRepository backed by DynamoDB.
type NoteDDBRepo struct {
	client *dynamodb.Client
	table  string
}

func NewNoteDDBRepo(client *dynamodb.Client) *NoteDDBRepo {
	return &NoteDDBRepo{client: client, table: ddb.TableName("notes")}
}

func (r *NoteDDBRepo) Save(ctx context.Context, n *journeyv1.Note) error {
	it := noteToItem(n)
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return fmt.Errorf("note marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      av,
	})
	return err
}

func (r *NoteDDBRepo) Get(ctx context.Context, id string) (*journeyv1.Note, error) {
	raw, err := queryByIDOnIndex(ctx, r.client, r.table, id)
	if err != nil {
		return nil, err
	}
	var it noteItem
	if err := attributevalue.UnmarshalMap(raw, &it); err != nil {
		return nil, err
	}
	return noteFromItem(it), nil
}

func (r *NoteDDBRepo) Delete(ctx context.Context, id string) error {
	n, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "TRIP#" + n.GetTripId()},
			"SK": &types.AttributeValueMemberS{Value: "NOTE#" + id},
		},
	})
	return err
}

func (r *NoteDDBRepo) ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Note, error) {
	keyCond := expression.Key("PK").Equal(expression.Value("TRIP#" + tripID)).
		And(expression.Key("SK").BeginsWith("NOTE#"))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("note list: %w", err)
	}
	results := make([]*journeyv1.Note, 0, len(out.Items))
	for _, item := range out.Items {
		var it noteItem
		if err := attributevalue.UnmarshalMap(item, &it); err != nil {
			continue
		}
		results = append(results, noteFromItem(it))
	}
	return results, nil
}

func (r *NoteDDBRepo) DeleteByTrip(ctx context.Context, tripID string) error {
	notes, err := r.ListByTrip(ctx, tripID)
	if err != nil {
		return err
	}
	if len(notes) == 0 {
		return nil
	}
	sks := make([]string, 0, len(notes))
	for _, n := range notes {
		sks = append(sks, "NOTE#"+n.GetId())
	}
	return ddb.BatchDelete(ctx, r.client, r.table, ddb.BatchDeletePK("TRIP#"+tripID, sks))
}

// ─── Checklist DDB ───────────────────────────────────────────────────────────

type checklistItem struct {
	PK        string `dynamodbav:"PK"` // TRIP#<tripID>
	SK        string `dynamodbav:"SK"` // CHECKLIST#<id>
	ID        string `dynamodbav:"id"`
	TripID    string `dynamodbav:"tripId"`
	Category  int32  `dynamodbav:"category"`
	Item      string `dynamodbav:"item"`
	IsChecked bool   `dynamodbav:"isChecked"`
	CreatedAt string `dynamodbav:"createdAt"`
	UpdatedAt string `dynamodbav:"updatedAt"`
}

func checklistToItem(c *journeyv1.ChecklistItem) checklistItem {
	it := checklistItem{
		PK:        "TRIP#" + c.GetTripId(),
		SK:        "CHECKLIST#" + c.GetId(),
		ID:        c.GetId(),
		TripID:    c.GetTripId(),
		Category:  int32(c.GetCategory()),
		Item:      c.GetItem(),
		IsChecked: c.GetIsChecked(),
	}
	if a := c.GetAudit(); a != nil {
		it.CreatedAt = a.GetCreatedAt().AsTime().Format(time.RFC3339)
		it.UpdatedAt = a.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return it
}

func checklistFromItem(it checklistItem) *journeyv1.ChecklistItem {
	c := &journeyv1.ChecklistItem{
		Id:        it.ID,
		TripId:    it.TripID,
		Category:  journeyv1.ChecklistCategory(it.Category),
		Item:      it.Item,
		IsChecked: it.IsChecked,
	}
	ca, _ := time.Parse(time.RFC3339, it.CreatedAt)
	ua, _ := time.Parse(time.RFC3339, it.UpdatedAt)
	c.Audit = &journeyv1.AuditTimestamps{
		CreatedAt: timestamppb.New(ca),
		UpdatedAt: timestamppb.New(ua),
	}
	return c
}

// ChecklistDDBRepo implements ChecklistRepository backed by DynamoDB.
type ChecklistDDBRepo struct {
	client *dynamodb.Client
	table  string
}

func NewChecklistDDBRepo(client *dynamodb.Client) *ChecklistDDBRepo {
	return &ChecklistDDBRepo{client: client, table: ddb.TableName("checklists")}
}

func (r *ChecklistDDBRepo) Save(ctx context.Context, c *journeyv1.ChecklistItem) error {
	it := checklistToItem(c)
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return fmt.Errorf("checklist marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      av,
	})
	return err
}

func (r *ChecklistDDBRepo) Get(ctx context.Context, id string) (*journeyv1.ChecklistItem, error) {
	raw, err := queryByIDOnIndex(ctx, r.client, r.table, id)
	if err != nil {
		return nil, err
	}
	var it checklistItem
	if err := attributevalue.UnmarshalMap(raw, &it); err != nil {
		return nil, err
	}
	return checklistFromItem(it), nil
}

func (r *ChecklistDDBRepo) Delete(ctx context.Context, id string) error {
	c, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "TRIP#" + c.GetTripId()},
			"SK": &types.AttributeValueMemberS{Value: "CHECKLIST#" + id},
		},
	})
	return err
}

func (r *ChecklistDDBRepo) ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.ChecklistItem, error) {
	keyCond := expression.Key("PK").Equal(expression.Value("TRIP#" + tripID)).
		And(expression.Key("SK").BeginsWith("CHECKLIST#"))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("checklist list: %w", err)
	}
	results := make([]*journeyv1.ChecklistItem, 0, len(out.Items))
	for _, item := range out.Items {
		var it checklistItem
		if err := attributevalue.UnmarshalMap(item, &it); err != nil {
			continue
		}
		results = append(results, checklistFromItem(it))
	}
	return results, nil
}

func (r *ChecklistDDBRepo) DeleteByTrip(ctx context.Context, tripID string) error {
	items, err := r.ListByTrip(ctx, tripID)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	sks := make([]string, 0, len(items))
	for _, c := range items {
		sks = append(sks, "CHECKLIST#"+c.GetId())
	}
	return ddb.BatchDelete(ctx, r.client, r.table, ddb.BatchDeletePK("TRIP#"+tripID, sks))
}

// ─── Reservation DDB ─────────────────────────────────────────────────────────

type reservationItem struct {
	PK              string `dynamodbav:"PK"` // TRIP#<tripID>
	SK              string `dynamodbav:"SK"` // RESERVATION#<id>
	ID              string `dynamodbav:"id"`
	TripID          string `dynamodbav:"tripId"`
	Type            int32  `dynamodbav:"type"`
	Vendor          string `dynamodbav:"vendor,omitempty"`
	ConfirmNumber   string `dynamodbav:"confirmNumber,omitempty"`
	ReservedAt      string `dynamodbav:"reservedAt,omitempty"`
	CostVal         int64  `dynamodbav:"costVal,omitempty"`
	CostCur         string `dynamodbav:"costCur,omitempty"`
	AttachmentS3Key string `dynamodbav:"attachmentS3Key,omitempty"`
	Notes           string `dynamodbav:"notes,omitempty"`
	CreatedAt       string `dynamodbav:"createdAt"`
	UpdatedAt       string `dynamodbav:"updatedAt"`
}

func reservationToItem(v *journeyv1.Reservation) reservationItem {
	it := reservationItem{
		PK:              "TRIP#" + v.GetTripId(),
		SK:              "RESERVATION#" + v.GetId(),
		ID:              v.GetId(),
		TripID:          v.GetTripId(),
		Type:            int32(v.GetType()),
		Vendor:          v.GetVendor(),
		ConfirmNumber:   v.GetConfirmNumber(),
		AttachmentS3Key: v.GetAttachmentS3Key(),
		Notes:           v.GetNotes(),
	}
	if t := v.GetReservedAt(); t != nil {
		it.ReservedAt = t.AsTime().Format(time.RFC3339)
	}
	if c := v.GetCost(); c != nil {
		it.CostVal = c.GetAmount()
		it.CostCur = c.GetCurrency()
	}
	if a := v.GetAudit(); a != nil {
		it.CreatedAt = a.GetCreatedAt().AsTime().Format(time.RFC3339)
		it.UpdatedAt = a.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return it
}

func reservationFromItem(it reservationItem) *journeyv1.Reservation {
	v := &journeyv1.Reservation{
		Id:              it.ID,
		TripId:          it.TripID,
		Type:            journeyv1.ReservationType(it.Type),
		Vendor:          it.Vendor,
		ConfirmNumber:   it.ConfirmNumber,
		AttachmentS3Key: it.AttachmentS3Key,
		Notes:           it.Notes,
	}
	if it.ReservedAt != "" {
		if t, err := time.Parse(time.RFC3339, it.ReservedAt); err == nil {
			v.ReservedAt = timestamppb.New(t)
		}
	}
	if it.CostCur != "" || it.CostVal != 0 {
		v.Cost = &journeyv1.Money{Currency: it.CostCur, Amount: it.CostVal}
	}
	ca, _ := time.Parse(time.RFC3339, it.CreatedAt)
	ua, _ := time.Parse(time.RFC3339, it.UpdatedAt)
	v.Audit = &journeyv1.AuditTimestamps{
		CreatedAt: timestamppb.New(ca),
		UpdatedAt: timestamppb.New(ua),
	}
	return v
}

// ReservationDDBRepo implements ReservationRepository backed by DynamoDB.
type ReservationDDBRepo struct {
	client *dynamodb.Client
	table  string
}

func NewReservationDDBRepo(client *dynamodb.Client) *ReservationDDBRepo {
	return &ReservationDDBRepo{client: client, table: ddb.TableName("reservations")}
}

func (r *ReservationDDBRepo) Save(ctx context.Context, v *journeyv1.Reservation) error {
	it := reservationToItem(v)
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return fmt.Errorf("reservation marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      av,
	})
	return err
}

func (r *ReservationDDBRepo) Get(ctx context.Context, id string) (*journeyv1.Reservation, error) {
	raw, err := queryByIDOnIndex(ctx, r.client, r.table, id)
	if err != nil {
		return nil, err
	}
	var it reservationItem
	if err := attributevalue.UnmarshalMap(raw, &it); err != nil {
		return nil, err
	}
	return reservationFromItem(it), nil
}

func (r *ReservationDDBRepo) Delete(ctx context.Context, id string) error {
	v, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "TRIP#" + v.GetTripId()},
			"SK": &types.AttributeValueMemberS{Value: "RESERVATION#" + id},
		},
	})
	return err
}

func (r *ReservationDDBRepo) ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Reservation, error) {
	keyCond := expression.Key("PK").Equal(expression.Value("TRIP#" + tripID)).
		And(expression.Key("SK").BeginsWith("RESERVATION#"))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("reservation list: %w", err)
	}
	results := make([]*journeyv1.Reservation, 0, len(out.Items))
	for _, item := range out.Items {
		var it reservationItem
		if err := attributevalue.UnmarshalMap(item, &it); err != nil {
			continue
		}
		results = append(results, reservationFromItem(it))
	}
	return results, nil
}

func (r *ReservationDDBRepo) DeleteByTrip(ctx context.Context, tripID string) error {
	items, err := r.ListByTrip(ctx, tripID)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	sks := make([]string, 0, len(items))
	for _, v := range items {
		sks = append(sks, "RESERVATION#"+v.GetId())
	}
	return ddb.BatchDelete(ctx, r.client, r.table, ddb.BatchDeletePK("TRIP#"+tripID, sks))
}
