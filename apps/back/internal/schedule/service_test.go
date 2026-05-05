package schedule

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

func TestCreate(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	sch, err := svc.Create(ctx, &journeyv1.CreateScheduleRequest{
		DayId:     "day-1",
		StartTime: "09:00",
		EndTime:   "10:00",
		Title:     "Visit temple",
		Region:    "Asakusa",
		Category:  journeyv1.ScheduleCategory_SCHEDULE_CATEGORY_SIGHTSEEING,
		PlaceName: "Senso-ji",
		Cost:      &journeyv1.Money{Currency: "JPY", Amount: 0},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if sch.GetId() == "" {
		t.Fatal("expected id")
	}
	if sch.GetOrder() != 1 {
		t.Fatalf("expected order=1, got %d", sch.GetOrder())
	}
	if sch.GetTitle() != "Visit temple" {
		t.Fatalf("title: %s", sch.GetTitle())
	}
}

func TestCreate_AutoIncrementOrder(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	s1, _ := svc.Create(ctx, &journeyv1.CreateScheduleRequest{
		DayId: "day-1", Title: "First",
	})
	s2, _ := svc.Create(ctx, &journeyv1.CreateScheduleRequest{
		DayId: "day-1", Title: "Second",
	})
	s3, _ := svc.Create(ctx, &journeyv1.CreateScheduleRequest{
		DayId: "day-1", Title: "Third",
	})

	if s1.GetOrder() != 1 || s2.GetOrder() != 2 || s3.GetOrder() != 3 {
		t.Fatalf("order auto-increment failed: %d, %d, %d", s1.GetOrder(), s2.GetOrder(), s3.GetOrder())
	}
}

func TestCreate_DayIDRequired(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.Create(ctx, &journeyv1.CreateScheduleRequest{
		Title: "No day",
	})
	if err == nil {
		t.Fatal("expected error for missing day_id")
	}
}

func TestListByDay(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	svc.Create(ctx, &journeyv1.CreateScheduleRequest{DayId: "day-1", Title: "A"})
	svc.Create(ctx, &journeyv1.CreateScheduleRequest{DayId: "day-1", Title: "B"})
	svc.Create(ctx, &journeyv1.CreateScheduleRequest{DayId: "day-2", Title: "C"})

	list, err := svc.List(ctx, "day-1")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2, got %d", len(list))
	}
	// sorted by order
	if list[0].GetOrder() > list[1].GetOrder() {
		t.Fatal("not sorted by order")
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	sch, _ := svc.Create(ctx, &journeyv1.CreateScheduleRequest{
		DayId:    "day-1",
		Title:    "Original",
		Category: journeyv1.ScheduleCategory_SCHEDULE_CATEGORY_SIGHTSEEING,
	})

	updated, err := svc.Update(ctx, &journeyv1.UpdateScheduleRequest{
		ScheduleId:  sch.GetId(),
		Title:       "Updated",
		Category:    journeyv1.ScheduleCategory_SCHEDULE_CATEGORY_SHOPPING,
		IsCompleted: true,
		Cost:        &journeyv1.Money{Currency: "JPY", Amount: 3000},
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.GetTitle() != "Updated" {
		t.Errorf("title: %s", updated.GetTitle())
	}
	if updated.GetCategory() != journeyv1.ScheduleCategory_SCHEDULE_CATEGORY_SHOPPING {
		t.Errorf("category: %v", updated.GetCategory())
	}
	if !updated.GetIsCompleted() {
		t.Error("is_completed should be true")
	}
	if updated.GetCost().GetAmount() != 3000 {
		t.Errorf("cost: %v", updated.GetCost())
	}
}

func TestUpdate_NotFound(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.Update(ctx, &journeyv1.UpdateScheduleRequest{
		ScheduleId: "nonexistent",
		Title:      "x",
	})
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	sch, _ := svc.Create(ctx, &journeyv1.CreateScheduleRequest{
		DayId: "day-1", Title: "To delete",
	})

	if err := svc.Delete(ctx, sch.GetId()); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	list, _ := svc.List(ctx, "day-1")
	if len(list) != 0 {
		t.Fatalf("expected empty after delete, got %d", len(list))
	}
}

func TestDelete_NotFound(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	if err := svc.Delete(ctx, "nonexistent"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestReorder(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	s1, _ := svc.Create(ctx, &journeyv1.CreateScheduleRequest{DayId: "day-1", Title: "A"})
	s2, _ := svc.Create(ctx, &journeyv1.CreateScheduleRequest{DayId: "day-1", Title: "B"})
	s3, _ := svc.Create(ctx, &journeyv1.CreateScheduleRequest{DayId: "day-1", Title: "C"})

	// Reverse order
	reordered, err := svc.Reorder(ctx, "day-1", []string{s3.GetId(), s2.GetId(), s1.GetId()})
	if err != nil {
		t.Fatalf("Reorder: %v", err)
	}
	if len(reordered) != 3 {
		t.Fatalf("expected 3, got %d", len(reordered))
	}
	// After reorder, order should be 1,2,3 with titles C,B,A
	for _, s := range reordered {
		switch s.GetTitle() {
		case "C":
			if s.GetOrder() != 1 {
				t.Errorf("C order: %d, want 1", s.GetOrder())
			}
		case "B":
			if s.GetOrder() != 2 {
				t.Errorf("B order: %d, want 2", s.GetOrder())
			}
		case "A":
			if s.GetOrder() != 3 {
				t.Errorf("A order: %d, want 3", s.GetOrder())
			}
		}
	}
}

func TestReorder_WrongDay(t *testing.T) {
	t.Parallel()
	svc := newTestService()
	ctx := context.Background()

	s1, _ := svc.Create(ctx, &journeyv1.CreateScheduleRequest{DayId: "day-1", Title: "A"})

	_, err := svc.Reorder(ctx, "day-2", []string{s1.GetId()})
	if err == nil {
		t.Fatal("expected error for wrong day")
	}
}
