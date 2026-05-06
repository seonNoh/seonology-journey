// Package handler - REST → gRPC 변환 핸들러.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/seonNoh/seonology-journey/apps/api/internal/auth"
	"github.com/seonNoh/seonology-journey/apps/api/internal/grpcclient"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// API - REST 어댑터.
type API struct {
	J journeyv1.JourneyServiceClient
}

// New - 생성.
func New(j journeyv1.JourneyServiceClient) *API { return &API{J: j} }

// Mount - chi 라우터에 모든 v1 라우트 부착.
func (a *API) Mount(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		// Trip
		r.Post("/trips", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.CreateTripRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			return a.J.CreateTrip(reqctx(r), req)
		}))
		r.Get("/trips", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.ListTripsRequest{Page: pageReqFrom(r)}
			if v := r.URL.Query().Get("status"); v != "" {
				if i, err := strconv.Atoi(v); err == nil {
					req.Status = journeyv1.TripStatus(i)
				}
			}
			return a.J.ListTrips(reqctx(r), req)
		}))
		r.Get("/trips/{tripId}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.GetTrip(reqctx(r), &journeyv1.GetTripRequest{TripId: chi.URLParam(r, "tripId")})
		}))
		r.Patch("/trips/{tripId}", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.UpdateTripRequest{TripId: chi.URLParam(r, "tripId")}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.TripId = chi.URLParam(r, "tripId")
			return a.J.UpdateTrip(reqctx(r), req)
		}))
		r.Delete("/trips/{tripId}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.DeleteTrip(reqctx(r), &journeyv1.DeleteTripRequest{TripId: chi.URLParam(r, "tripId")})
		}))

		// Day
		r.Get("/trips/{tripId}/days", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListDays(reqctx(r), &journeyv1.ListDaysRequest{TripId: chi.URLParam(r, "tripId")})
		}))
		r.Patch("/days/{dayId}", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.UpdateDayRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.DayId = chi.URLParam(r, "dayId")
			return a.J.UpdateDay(reqctx(r), req)
		}))

		// Schedule
		r.Get("/days/{dayId}/schedules", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListSchedules(reqctx(r), &journeyv1.ListSchedulesRequest{DayId: chi.URLParam(r, "dayId")})
		}))
		r.Post("/days/{dayId}/schedules", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.CreateScheduleRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.DayId = chi.URLParam(r, "dayId")
			return a.J.CreateSchedule(reqctx(r), req)
		}))
		r.Patch("/schedules/{scheduleId}", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.UpdateScheduleRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.ScheduleId = chi.URLParam(r, "scheduleId")
			return a.J.UpdateSchedule(reqctx(r), req)
		}))
		r.Delete("/schedules/{scheduleId}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.DeleteSchedule(reqctx(r), &journeyv1.DeleteScheduleRequest{ScheduleId: chi.URLParam(r, "scheduleId")})
		}))
		r.Post("/days/{dayId}/schedules:reorder", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.ReorderSchedulesRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.DayId = chi.URLParam(r, "dayId")
			return a.J.ReorderSchedules(reqctx(r), req)
		}))

		// Meal
		r.Get("/days/{dayId}/meals", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListMeals(reqctx(r), &journeyv1.ListMealsRequest{DayId: chi.URLParam(r, "dayId")})
		}))
		r.Put("/days/{dayId}/meals", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.UpsertMealRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.DayId = chi.URLParam(r, "dayId")
			return a.J.UpsertMeal(reqctx(r), req)
		}))
		r.Delete("/days/{dayId}/meals/{mealType}", a.handle(func(r *http.Request) (proto.Message, error) {
			t, _ := strconv.Atoi(chi.URLParam(r, "mealType"))
			return a.J.DeleteMeal(reqctx(r), &journeyv1.DeleteMealRequest{
				DayId: chi.URLParam(r, "dayId"), MealType: journeyv1.MealType(t),
			})
		}))

		// Accommodation
		r.Get("/days/{dayId}/accommodation", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.GetAccommodation(reqctx(r), &journeyv1.GetAccommodationRequest{DayId: chi.URLParam(r, "dayId")})
		}))
		r.Put("/days/{dayId}/accommodation", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.UpsertAccommodationRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.DayId = chi.URLParam(r, "dayId")
			return a.J.UpsertAccommodation(reqctx(r), req)
		}))
		r.Delete("/days/{dayId}/accommodation", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.DeleteAccommodation(reqctx(r), &journeyv1.DeleteAccommodationRequest{DayId: chi.URLParam(r, "dayId")})
		}))

		// Expense
		r.Get("/trips/{tripId}/expenses", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListExpenses(reqctx(r), &journeyv1.ListExpensesRequest{TripId: chi.URLParam(r, "tripId")})
		}))
		r.Post("/trips/{tripId}/expenses", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.CreateExpenseRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.TripId = chi.URLParam(r, "tripId")
			return a.J.CreateExpense(reqctx(r), req)
		}))
		r.Patch("/expenses/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.UpdateExpenseRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.ExpenseId = chi.URLParam(r, "id")
			return a.J.UpdateExpense(reqctx(r), req)
		}))
		r.Delete("/expenses/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.DeleteExpense(reqctx(r), &journeyv1.DeleteExpenseRequest{ExpenseId: chi.URLParam(r, "id")})
		}))
		r.Get("/trips/{tripId}/expense-summary", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.GetExpenseSummary(reqctx(r), &journeyv1.GetExpenseSummaryRequest{TripId: chi.URLParam(r, "tripId")})
		}))

		// Note
		r.Get("/trips/{tripId}/notes", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListNotes(reqctx(r), &journeyv1.ListNotesRequest{TripId: chi.URLParam(r, "tripId")})
		}))
		r.Post("/trips/{tripId}/notes", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.CreateNoteRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.TripId = chi.URLParam(r, "tripId")
			return a.J.CreateNote(reqctx(r), req)
		}))
		r.Patch("/notes/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.UpdateNoteRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.NoteId = chi.URLParam(r, "id")
			return a.J.UpdateNote(reqctx(r), req)
		}))
		r.Delete("/notes/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.DeleteNote(reqctx(r), &journeyv1.DeleteNoteRequest{NoteId: chi.URLParam(r, "id")})
		}))

		// Checklist
		r.Get("/trips/{tripId}/checklist", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListChecklistItems(reqctx(r), &journeyv1.ListChecklistItemsRequest{TripId: chi.URLParam(r, "tripId")})
		}))
		r.Post("/trips/{tripId}/checklist", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.CreateChecklistItemRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.TripId = chi.URLParam(r, "tripId")
			return a.J.CreateChecklistItem(reqctx(r), req)
		}))
		r.Patch("/checklist/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.UpdateChecklistItemRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.ItemId = chi.URLParam(r, "id")
			return a.J.UpdateChecklistItem(reqctx(r), req)
		}))
		r.Delete("/checklist/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.DeleteChecklistItem(reqctx(r), &journeyv1.DeleteChecklistItemRequest{ItemId: chi.URLParam(r, "id")})
		}))

		// Reservation
		r.Get("/trips/{tripId}/reservations", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListReservations(reqctx(r), &journeyv1.ListReservationsRequest{TripId: chi.URLParam(r, "tripId")})
		}))
		r.Post("/trips/{tripId}/reservations", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.CreateReservationRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.TripId = chi.URLParam(r, "tripId")
			return a.J.CreateReservation(reqctx(r), req)
		}))
		r.Patch("/reservations/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.UpdateReservationRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.ReservationId = chi.URLParam(r, "id")
			return a.J.UpdateReservation(reqctx(r), req)
		}))
		r.Delete("/reservations/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.DeleteReservation(reqctx(r), &journeyv1.DeleteReservationRequest{ReservationId: chi.URLParam(r, "id")})
		}))

		// Companion
		r.Get("/trips/{tripId}/companions", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListCompanions(reqctx(r), &journeyv1.ListCompanionsRequest{TripId: chi.URLParam(r, "tripId")})
		}))
		r.Post("/trips/{tripId}/companions", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.AddCompanionRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.TripId = chi.URLParam(r, "tripId")
			return a.J.AddCompanion(reqctx(r), req)
		}))
		r.Patch("/trips/{tripId}/companions/{memberId}", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.UpdateCompanionRoleRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.TripId = chi.URLParam(r, "tripId")
			req.MemberId = chi.URLParam(r, "memberId")
			return a.J.UpdateCompanionRole(reqctx(r), req)
		}))
		r.Delete("/trips/{tripId}/companions/{memberId}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.RemoveCompanion(reqctx(r), &journeyv1.RemoveCompanionRequest{
				TripId: chi.URLParam(r, "tripId"), MemberId: chi.URLParam(r, "memberId"),
			})
		}))

		// Tag
		r.Get("/tags", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListTags(reqctx(r), &journeyv1.ListTagsRequest{})
		}))
		r.Post("/tags", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.CreateTagRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			return a.J.CreateTag(reqctx(r), req)
		}))
		r.Delete("/tags/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.DeleteTag(reqctx(r), &journeyv1.DeleteTagRequest{TagId: chi.URLParam(r, "id")})
		}))
		r.Get("/trips/{tripId}/tags", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListTripTags(reqctx(r), &journeyv1.ListTripTagsRequest{TripId: chi.URLParam(r, "tripId")})
		}))
		r.Put("/trips/{tripId}/tags/{tagId}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.AttachTag(reqctx(r), &journeyv1.AttachTagRequest{
				TripId: chi.URLParam(r, "tripId"), TagId: chi.URLParam(r, "tagId"),
			})
		}))
		r.Delete("/trips/{tripId}/tags/{tagId}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.DetachTag(reqctx(r), &journeyv1.DetachTagRequest{
				TripId: chi.URLParam(r, "tripId"), TagId: chi.URLParam(r, "tagId"),
			})
		}))

		// Template
		r.Get("/templates", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListTemplates(reqctx(r), &journeyv1.ListTemplatesRequest{})
		}))
		r.Post("/trips/{tripId}:save-as-template", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.SaveTripAsTemplateRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.TripId = chi.URLParam(r, "tripId")
			return a.J.SaveTripAsTemplate(reqctx(r), req)
		}))
		r.Post("/templates/{id}:create-trip", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.CreateTripFromTemplateRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.TemplateId = chi.URLParam(r, "id")
			return a.J.CreateTripFromTemplate(reqctx(r), req)
		}))
		r.Delete("/templates/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.DeleteTemplate(reqctx(r), &journeyv1.DeleteTemplateRequest{TemplateId: chi.URLParam(r, "id")})
		}))

		// FavoritePlace
		r.Get("/favorites", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListFavoritePlaces(reqctx(r), &journeyv1.ListFavoritePlacesRequest{})
		}))
		r.Post("/favorites", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.AddFavoritePlaceRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			return a.J.AddFavoritePlace(reqctx(r), req)
		}))
		r.Delete("/favorites/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.RemoveFavoritePlace(reqctx(r), &journeyv1.RemoveFavoritePlaceRequest{PlaceId: chi.URLParam(r, "id")})
		}))

		// Share
		r.Post("/shares", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.CreateShareRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			return a.J.CreateShare(reqctx(r), req)
		}))
		r.Delete("/shares/{code}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.RevokeShare(reqctx(r), &journeyv1.RevokeShareRequest{Code: chi.URLParam(r, "code")})
		}))
		r.Get("/shares/{code}", a.handlePublic(func(r *http.Request) (proto.Message, error) {
			return a.J.GetShare(r.Context(), &journeyv1.GetShareRequest{Code: chi.URLParam(r, "code")})
		}))

		// Media
		r.Post("/trips/{tripId}/media:upload-url", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.GetUploadUrlRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.TripId = chi.URLParam(r, "tripId")
			return a.J.GetUploadUrl(reqctx(r), req)
		}))
		r.Post("/trips/{tripId}/media:confirm", a.handle(func(r *http.Request) (proto.Message, error) {
			req := &journeyv1.ConfirmUploadRequest{}
			if err := decode(r, req); err != nil {
				return nil, err
			}
			req.TripId = chi.URLParam(r, "tripId")
			return a.J.ConfirmUpload(reqctx(r), req)
		}))
		r.Get("/trips/{tripId}/media", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.ListMedia(reqctx(r), &journeyv1.ListMediaRequest{
				TripId: chi.URLParam(r, "tripId"),
				DayId:  r.URL.Query().Get("day_id"),
				Page:   pageReqFrom(r),
			})
		}))
		r.Delete("/media/{id}", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.DeleteMedia(reqctx(r), &journeyv1.DeleteMediaRequest{MediaId: chi.URLParam(r, "id")})
		}))
		r.Get("/media/{id}/url", a.handle(func(r *http.Request) (proto.Message, error) {
			thumb := r.URL.Query().Get("thumbnail") == "true"
			return a.J.GetMediaUrl(reqctx(r), &journeyv1.GetMediaUrlRequest{MediaId: chi.URLParam(r, "id"), Thumbnail: thumb})
		}))

		// Statistics
		r.Get("/trips/{tripId}/statistics", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.GetTripStatistics(reqctx(r), &journeyv1.GetTripStatisticsRequest{TripId: chi.URLParam(r, "tripId")})
		}))
		r.Get("/statistics/yearly", a.handle(func(r *http.Request) (proto.Message, error) {
			y, _ := strconv.Atoi(r.URL.Query().Get("year"))
			return a.J.GetYearlyStatistics(reqctx(r), &journeyv1.GetYearlyStatisticsRequest{Year: int32(y)})
		}))

		// External
		r.Get("/external/geocode", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.Geocode(reqctx(r), &journeyv1.GeocodeRequest{Query: r.URL.Query().Get("q")})
		}))
		r.Get("/external/exchange-rate", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.GetExchangeRate(reqctx(r), &journeyv1.GetExchangeRateRequest{
				Base: r.URL.Query().Get("base"), Target: r.URL.Query().Get("target"), Date: r.URL.Query().Get("date"),
			})
		}))
		r.Get("/external/weather", a.handle(func(r *http.Request) (proto.Message, error) {
			return a.J.GetWeatherForecast(reqctx(r), &journeyv1.GetWeatherForecastRequest{Date: r.URL.Query().Get("date")})
		}))
	})
}

// reqctx - 인증 컨텍스트에 x-user-id metadata 부착.
func reqctx(r *http.Request) context.Context {
	uid, _ := auth.UserID(r.Context())
	return grpcclient.WithUser(r.Context(), uid)
}

// handle - 인증 필요한 핸들러 wrapper.
func (a *API) handle(fn func(*http.Request) (proto.Message, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := auth.UserID(r.Context()); err != nil {
			writeErr(w, status.Error(codes.Unauthenticated, "unauthenticated"))
			return
		}
		resp, err := fn(r)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeResp(w, resp)
	}
}

// handlePublic - 인증 없이 호출 (e.g. GetShare).
func (a *API) handlePublic(fn func(*http.Request) (proto.Message, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := fn(r)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeResp(w, resp)
	}
}

func decode(r *http.Request, m proto.Message) error {
	body, err := io.ReadAll(http.MaxBytesReader(nil, r.Body, 1<<20))
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return nil
	}
	return protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(body, m)
}

func writeResp(w http.ResponseWriter, m proto.Message) {
	w.Header().Set("Content-Type", "application/json")
	b, err := protojson.MarshalOptions{EmitUnpopulated: true, UseProtoNames: false}.Marshal(m)
	if err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func writeErr(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.NotFound:
			code = http.StatusNotFound
		case codes.PermissionDenied:
			code = http.StatusForbidden
		case codes.Unauthenticated:
			code = http.StatusUnauthorized
		case codes.InvalidArgument, codes.FailedPrecondition:
			code = http.StatusBadRequest
		case codes.AlreadyExists:
			code = http.StatusConflict
		case codes.Unimplemented:
			code = http.StatusNotImplemented
		case codes.OK:
			code = http.StatusOK
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// 컴파일 안정용.
var _ = errors.Is

// pageReqFrom extracts cursor/limit query params into a PageRequest.
// Returns nil when neither is supplied so callers omit the field in the
// gRPC request and let the server apply its default limit.
func pageReqFrom(r *http.Request) *journeyv1.PageRequest {
	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")
	if cursor == "" && limitStr == "" {
		return nil
	}
	pr := &journeyv1.PageRequest{Cursor: cursor}
	if limitStr != "" {
		if n, err := strconv.Atoi(limitStr); err == nil && n > 0 && n <= 100 {
			pr.Limit = int32(n)
		}
	}
	return pr
}
