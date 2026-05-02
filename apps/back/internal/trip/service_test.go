package trip

import (
	"context"
	"testing"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
)

func TestServiceCRUD(t *testing.T) {
	t.Parallel()
	svc := NewService(NewMemoryRepo())
	ctx := context.Background()
	owner := "u-1"

	created, err := svc.Create(ctx, owner, &journeyv1.CreateTripRequest{
		Title:       "Tokyo 5 days",
		Description: "Spring trip",
		StartDate:   "2026-04-01",
		EndDate:     "2026-04-05",
		Destination: "Tokyo",
		CountryCode: "JP",
		TotalBudget: &journeyv1.Money{Currency: "JPY", Amount: 200000},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.GetId() == "" {
		t.Fatalf("expected id")
	}
	if created.GetStatus() != journeyv1.TripStatus_TRIP_STATUS_PLANNING {
		t.Fatalf("unexpected status: %v", created.GetStatus())
	}

	got, err := svc.Get(ctx, owner, created.GetId())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.GetTitle() != "Tokyo 5 days" {
		t.Fatalf("Get title mismatch: %s", got.GetTitle())
	}

	if _, err := svc.Get(ctx, "u-other", created.GetId()); err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	list, err := svc.List(ctx, owner, journeyv1.TripStatus_TRIP_STATUS_UNSPECIFIED)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("List len=%d", len(list))
	}

	updated, err := svc.Update(ctx, owner, &journeyv1.UpdateTripRequest{
		TripId: created.GetId(),
		Title:  "Tokyo 6 days",
		Status: journeyv1.TripStatus_TRIP_STATUS_ONGOING,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.GetTitle() != "Tokyo 6 days" || updated.GetStatus() != journeyv1.TripStatus_TRIP_STATUS_ONGOING {
		t.Fatalf("update not applied: %+v", updated)
	}

	if err := svc.Delete(ctx, owner, created.GetId()); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := svc.Get(ctx, owner, created.GetId()); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
