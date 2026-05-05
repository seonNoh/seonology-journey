package trip_test

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/seonNoh/seonology-journey/apps/back/internal/trip"
)

func newTestClient(t *testing.T) *dynamodb.Client {
	t.Helper()
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("local", "local", "")),
	)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	// Verify connectivity.
	_, err = client.ListTables(context.Background(), &dynamodb.ListTablesInput{})
	if err != nil {
		t.Skipf("DynamoDB Local not available: %v", err)
	}
	return client
}

func TestDDBRepo_CreateAndGet(t *testing.T) {
	client := newTestClient(t)
	repo := trip.NewDDBRepo(client)
	ctx := context.Background()

	now := timestamppb.Now()
	tr := &journeyv1.Trip{
		Id:          "test-trip-001",
		OwnerId:     "user-001",
		Title:       "Tokyo Trip",
		Description: "Spring vacation",
		StartDate:   "2026-04-01",
		EndDate:     "2026-04-05",
		Status:      journeyv1.TripStatus_TRIP_STATUS_PLANNING,
		Destination: "Tokyo",
		CountryCode: "JP",
		TotalBudget: &journeyv1.Money{Amount: 500000, Currency: "JPY"},
		Audit:       &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}

	// Create
	if err := repo.Create(ctx, tr); err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Get
	got, err := repo.Get(ctx, "test-trip-001")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.GetTitle() != "Tokyo Trip" {
		t.Errorf("title = %q, want %q", got.GetTitle(), "Tokyo Trip")
	}
	if got.GetTotalBudget().GetAmount() != 500000 {
		t.Errorf("budget = %d, want 500000", got.GetTotalBudget().GetAmount())
	}
	if got.GetTotalBudget().GetCurrency() != "JPY" {
		t.Errorf("currency = %q, want JPY", got.GetTotalBudget().GetCurrency())
	}

	// Cleanup
	_ = repo.Delete(ctx, "test-trip-001")
}

func TestDDBRepo_ListByOwner(t *testing.T) {
	client := newTestClient(t)
	repo := trip.NewDDBRepo(client)
	ctx := context.Background()

	now := timestamppb.Now()
	trips := []*journeyv1.Trip{
		{Id: "list-001", OwnerId: "user-list", Title: "Trip A", StartDate: "2026-03-01", EndDate: "2026-03-03", Status: journeyv1.TripStatus_TRIP_STATUS_PLANNING, Audit: &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now}},
		{Id: "list-002", OwnerId: "user-list", Title: "Trip B", StartDate: "2026-04-01", EndDate: "2026-04-03", Status: journeyv1.TripStatus_TRIP_STATUS_ONGOING, Audit: &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now}},
		{Id: "list-003", OwnerId: "user-list", Title: "Trip C", StartDate: "2026-05-01", EndDate: "2026-05-03", Status: journeyv1.TripStatus_TRIP_STATUS_PLANNING, Audit: &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now}},
	}
	for _, tr := range trips {
		if err := repo.Create(ctx, tr); err != nil {
			t.Fatalf("Create %s: %v", tr.Id, err)
		}
	}

	// List all
	all, err := repo.ListByOwner(ctx, "user-list", journeyv1.TripStatus_TRIP_STATUS_UNSPECIFIED)
	if err != nil {
		t.Fatalf("ListByOwner all: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("got %d trips, want 3", len(all))
	}

	// List filtered by status
	planning, err := repo.ListByOwner(ctx, "user-list", journeyv1.TripStatus_TRIP_STATUS_PLANNING)
	if err != nil {
		t.Fatalf("ListByOwner planning: %v", err)
	}
	if len(planning) != 2 {
		t.Errorf("got %d planning trips, want 2", len(planning))
	}

	// Cleanup
	for _, tr := range trips {
		_ = repo.Delete(ctx, tr.Id)
	}
}

func TestDDBRepo_Update(t *testing.T) {
	client := newTestClient(t)
	repo := trip.NewDDBRepo(client)
	ctx := context.Background()

	now := timestamppb.Now()
	tr := &journeyv1.Trip{
		Id:        "update-001",
		OwnerId:   "user-upd",
		Title:     "Original",
		StartDate: "2026-06-01",
		EndDate:   "2026-06-05",
		Status:    journeyv1.TripStatus_TRIP_STATUS_PLANNING,
		Audit:     &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	if err := repo.Create(ctx, tr); err != nil {
		t.Fatalf("Create: %v", err)
	}

	tr.Title = "Updated"
	tr.Status = journeyv1.TripStatus_TRIP_STATUS_ONGOING
	if err := repo.Update(ctx, tr); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err := repo.Get(ctx, "update-001")
	if err != nil {
		t.Fatalf("Get after update: %v", err)
	}
	if got.GetTitle() != "Updated" {
		t.Errorf("title = %q, want Updated", got.GetTitle())
	}
	if got.GetStatus() != journeyv1.TripStatus_TRIP_STATUS_ONGOING {
		t.Errorf("status = %v, want ACTIVE", got.GetStatus())
	}

	// Cleanup
	_ = repo.Delete(ctx, "update-001")
}

func TestDDBRepo_Delete(t *testing.T) {
	client := newTestClient(t)
	repo := trip.NewDDBRepo(client)
	ctx := context.Background()

	now := timestamppb.Now()
	tr := &journeyv1.Trip{
		Id:        "del-001",
		OwnerId:   "user-del",
		Title:     "To Delete",
		StartDate: "2026-07-01",
		EndDate:   "2026-07-05",
		Status:    journeyv1.TripStatus_TRIP_STATUS_PLANNING,
		Audit:     &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	if err := repo.Create(ctx, tr); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := repo.Delete(ctx, "del-001"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := repo.Get(ctx, "del-001")
	if err == nil {
		t.Error("expected ErrNotFound after delete")
	}
}

func TestDDBRepo_GetNotFound(t *testing.T) {
	client := newTestClient(t)
	repo := trip.NewDDBRepo(client)
	ctx := context.Background()

	_, err := repo.Get(ctx, "nonexistent-id")
	if err == nil {
		t.Error("expected error for nonexistent trip")
	}
}
