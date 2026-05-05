package day

import (
	"context"
	"testing"
	"time"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
)

func fixedNow() time.Time {
	return time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
}

func newTestService() *Service {
	svc := NewService(NewMemoryRepo())
	svc.now = fixedNow
	return svc
}

func TestGenerateForTrip(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	days, err := svc.GenerateForTrip(ctx, "trip-1", "2026-04-01", "2026-04-05")
	if err != nil {
		t.Fatalf("GenerateForTrip: %v", err)
	}
	if len(days) != 5 {
		t.Fatalf("expected 5 days, got %d", len(days))
	}

	// day_number + date + day_of_week
	want := []struct {
		num  int32
		date string
		dow  string
	}{
		{1, "2026-04-01", "수"},
		{2, "2026-04-02", "목"},
		{3, "2026-04-03", "금"},
		{4, "2026-04-04", "토"},
		{5, "2026-04-05", "일"},
	}
	for i, w := range want {
		d := days[i]
		if d.GetDayNumber() != w.num {
			t.Errorf("day[%d] number: got %d, want %d", i, d.GetDayNumber(), w.num)
		}
		if d.GetDate() != w.date {
			t.Errorf("day[%d] date: got %s, want %s", i, d.GetDate(), w.date)
		}
		if d.GetDayOfWeek() != w.dow {
			t.Errorf("day[%d] day_of_week: got %s, want %s", i, d.GetDayOfWeek(), w.dow)
		}
		if d.GetTripId() != "trip-1" {
			t.Errorf("day[%d] trip_id: got %s", i, d.GetTripId())
		}
		if d.GetId() == "" {
			t.Errorf("day[%d] id empty", i)
		}
	}
}

func TestGenerateForTrip_ReplacesExisting(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	// First generation
	_, err := svc.GenerateForTrip(ctx, "trip-1", "2026-04-01", "2026-04-03")
	if err != nil {
		t.Fatalf("first generate: %v", err)
	}

	// Second generation should replace
	days, err := svc.GenerateForTrip(ctx, "trip-1", "2026-04-10", "2026-04-12")
	if err != nil {
		t.Fatalf("second generate: %v", err)
	}
	if len(days) != 3 {
		t.Fatalf("expected 3 days, got %d", len(days))
	}
	if days[0].GetDate() != "2026-04-10" {
		t.Fatalf("expected new start date, got %s", days[0].GetDate())
	}

	// List should only have the new 3 days
	list, err := svc.List(ctx, "trip-1")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 after replace, got %d", len(list))
	}
}

func TestGenerateForTrip_InvalidDates(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	// end before start
	_, err := svc.GenerateForTrip(ctx, "trip-1", "2026-04-05", "2026-04-01")
	if err == nil {
		t.Fatal("expected error for end before start")
	}

	// invalid format
	_, err = svc.GenerateForTrip(ctx, "trip-1", "2026/04/01", "2026/04/05")
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestGetAndList(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	days, _ := svc.GenerateForTrip(ctx, "trip-1", "2026-04-01", "2026-04-02")

	got, err := svc.Get(ctx, days[0].GetId())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.GetDate() != "2026-04-01" {
		t.Fatalf("Get date: %s", got.GetDate())
	}

	// Not found
	_, err = svc.Get(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	list, err := svc.List(ctx, "trip-1")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("List len=%d, want 2", len(list))
	}
	// sorted by day_number
	if list[0].GetDayNumber() > list[1].GetDayNumber() {
		t.Fatal("List not sorted by day_number")
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	days, _ := svc.GenerateForTrip(ctx, "trip-1", "2026-04-01", "2026-04-01")
	dayID := days[0].GetId()

	updated, err := svc.Update(ctx, &journeyv1.UpdateDayRequest{
		DayId:        dayID,
		Region:       "Shibuya",
		Weather:      "Sunny",
		DailySummary: "Great day exploring",
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.GetRegion() != "Shibuya" {
		t.Errorf("region: got %s", updated.GetRegion())
	}
	if updated.GetWeather() != "Sunny" {
		t.Errorf("weather: got %s", updated.GetWeather())
	}
	if updated.GetDailySummary() != "Great day exploring" {
		t.Errorf("summary: got %s", updated.GetDailySummary())
	}
}

func TestUpdate_NotFound(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.Update(ctx, &journeyv1.UpdateDayRequest{
		DayId:  "nonexistent",
		Region: "Tokyo",
	})
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteByTrip(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	svc.GenerateForTrip(ctx, "trip-1", "2026-04-01", "2026-04-03")
	svc.GenerateForTrip(ctx, "trip-2", "2026-05-01", "2026-05-02")

	if err := svc.DeleteByTrip(ctx, "trip-1"); err != nil {
		t.Fatalf("DeleteByTrip: %v", err)
	}

	list1, _ := svc.List(ctx, "trip-1")
	if len(list1) != 0 {
		t.Fatalf("trip-1 days should be deleted, got %d", len(list1))
	}

	// trip-2 unaffected
	list2, _ := svc.List(ctx, "trip-2")
	if len(list2) != 2 {
		t.Fatalf("trip-2 days should remain, got %d", len(list2))
	}
}
