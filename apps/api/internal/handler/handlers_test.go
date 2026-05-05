package handler_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/seonNoh/seonology-journey/apps/api/internal/auth"
	"github.com/seonNoh/seonology-journey/apps/api/internal/handler"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mockJourneyClient - 테스트용 mock.
type mockJourneyClient struct {
	journeyv1.JourneyServiceClient
	createTripFn  func(ctx context.Context, in *journeyv1.CreateTripRequest, opts ...grpc.CallOption) (*journeyv1.CreateTripResponse, error)
	listTripsFn   func(ctx context.Context, in *journeyv1.ListTripsRequest, opts ...grpc.CallOption) (*journeyv1.ListTripsResponse, error)
	getTripFn     func(ctx context.Context, in *journeyv1.GetTripRequest, opts ...grpc.CallOption) (*journeyv1.GetTripResponse, error)
	deleteTripFn  func(ctx context.Context, in *journeyv1.DeleteTripRequest, opts ...grpc.CallOption) (*journeyv1.DeleteTripResponse, error)
	listDaysFn    func(ctx context.Context, in *journeyv1.ListDaysRequest, opts ...grpc.CallOption) (*journeyv1.ListDaysResponse, error)
	createShareFn func(ctx context.Context, in *journeyv1.CreateShareRequest, opts ...grpc.CallOption) (*journeyv1.CreateShareResponse, error)
	getShareFn    func(ctx context.Context, in *journeyv1.GetShareRequest, opts ...grpc.CallOption) (*journeyv1.GetShareResponse, error)
}

func (m *mockJourneyClient) CreateTrip(ctx context.Context, in *journeyv1.CreateTripRequest, opts ...grpc.CallOption) (*journeyv1.CreateTripResponse, error) {
	if m.createTripFn != nil {
		return m.createTripFn(ctx, in, opts...)
	}
	return &journeyv1.CreateTripResponse{}, nil
}

func (m *mockJourneyClient) ListTrips(ctx context.Context, in *journeyv1.ListTripsRequest, opts ...grpc.CallOption) (*journeyv1.ListTripsResponse, error) {
	if m.listTripsFn != nil {
		return m.listTripsFn(ctx, in, opts...)
	}
	return &journeyv1.ListTripsResponse{}, nil
}

func (m *mockJourneyClient) GetTrip(ctx context.Context, in *journeyv1.GetTripRequest, opts ...grpc.CallOption) (*journeyv1.GetTripResponse, error) {
	if m.getTripFn != nil {
		return m.getTripFn(ctx, in, opts...)
	}
	return &journeyv1.GetTripResponse{}, nil
}

func (m *mockJourneyClient) DeleteTrip(ctx context.Context, in *journeyv1.DeleteTripRequest, opts ...grpc.CallOption) (*journeyv1.DeleteTripResponse, error) {
	if m.deleteTripFn != nil {
		return m.deleteTripFn(ctx, in, opts...)
	}
	return &journeyv1.DeleteTripResponse{}, nil
}

func (m *mockJourneyClient) ListDays(ctx context.Context, in *journeyv1.ListDaysRequest, opts ...grpc.CallOption) (*journeyv1.ListDaysResponse, error) {
	if m.listDaysFn != nil {
		return m.listDaysFn(ctx, in, opts...)
	}
	return &journeyv1.ListDaysResponse{}, nil
}

func (m *mockJourneyClient) CreateShare(ctx context.Context, in *journeyv1.CreateShareRequest, opts ...grpc.CallOption) (*journeyv1.CreateShareResponse, error) {
	if m.createShareFn != nil {
		return m.createShareFn(ctx, in, opts...)
	}
	return &journeyv1.CreateShareResponse{}, nil
}

func (m *mockJourneyClient) GetShare(ctx context.Context, in *journeyv1.GetShareRequest, opts ...grpc.CallOption) (*journeyv1.GetShareResponse, error) {
	if m.getShareFn != nil {
		return m.getShareFn(ctx, in, opts...)
	}
	return &journeyv1.GetShareResponse{}, nil
}

// authedCtx - 인증된 컨텍스트 생성.
func authedCtx() context.Context {
	return auth.TestContext(context.Background(), "user-123", "testuser")
}

// setupRouter - 테스트용 라우터 구성.
func setupRouter(mock *mockJourneyClient) *chi.Mux {
	r := chi.NewRouter()
	api := handler.New(mock)
	api.Mount(r)
	return r
}

func TestHandlers(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		authed     bool
		setupMock  func(*mockJourneyClient)
		wantStatus int
		wantBody   string
	}{
		// Trip CRUD
		{
			name:       "POST /trips - success",
			method:     http.MethodPost,
			path:       "/api/v1/trips",
			body:       `{"title":"Tokyo Trip"}`,
			authed:     true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "POST /trips - unauthenticated",
			method:     http.MethodPost,
			path:       "/api/v1/trips",
			body:       `{"title":"Tokyo Trip"}`,
			authed:     false,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "GET /trips - success",
			method:     http.MethodGet,
			path:       "/api/v1/trips",
			authed:     true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "GET /trips/{tripId} - success",
			method:     http.MethodGet,
			path:       "/api/v1/trips/trip-001",
			authed:     true,
			wantStatus: http.StatusOK,
		},
		{
			name:   "GET /trips/{tripId} - not found",
			method: http.MethodGet,
			path:   "/api/v1/trips/trip-999",
			authed: true,
			setupMock: func(m *mockJourneyClient) {
				m.getTripFn = func(ctx context.Context, in *journeyv1.GetTripRequest, opts ...grpc.CallOption) (*journeyv1.GetTripResponse, error) {
					return nil, status.Error(codes.NotFound, "trip not found")
				}
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "trip not found",
		},
		{
			name:       "DELETE /trips/{tripId} - success",
			method:     http.MethodDelete,
			path:       "/api/v1/trips/trip-001",
			authed:     true,
			wantStatus: http.StatusOK,
		},
		{
			name:   "DELETE /trips/{tripId} - permission denied",
			method: http.MethodDelete,
			path:   "/api/v1/trips/trip-001",
			authed: true,
			setupMock: func(m *mockJourneyClient) {
				m.deleteTripFn = func(ctx context.Context, in *journeyv1.DeleteTripRequest, opts ...grpc.CallOption) (*journeyv1.DeleteTripResponse, error) {
					return nil, status.Error(codes.PermissionDenied, "not owner")
				}
			},
			wantStatus: http.StatusForbidden,
			wantBody:   "not owner",
		},
		// Day
		{
			name:       "GET /trips/{tripId}/days - success",
			method:     http.MethodGet,
			path:       "/api/v1/trips/trip-001/days",
			authed:     true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "GET /trips/{tripId}/days - unauthenticated",
			method:     http.MethodGet,
			path:       "/api/v1/trips/trip-001/days",
			authed:     false,
			wantStatus: http.StatusUnauthorized,
		},
		// Share (public endpoint)
		{
			name:       "GET /shares/{code} - public access",
			method:     http.MethodGet,
			path:       "/api/v1/shares/abc123",
			authed:     false,
			wantStatus: http.StatusOK,
		},
		{
			name:   "GET /shares/{code} - not found",
			method: http.MethodGet,
			path:   "/api/v1/shares/invalid",
			authed: false,
			setupMock: func(m *mockJourneyClient) {
				m.getShareFn = func(ctx context.Context, in *journeyv1.GetShareRequest, opts ...grpc.CallOption) (*journeyv1.GetShareResponse, error) {
					return nil, status.Error(codes.NotFound, "share not found")
				}
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "share not found",
		},
		// gRPC internal error
		{
			name:   "POST /trips - internal error",
			method: http.MethodPost,
			path:   "/api/v1/trips",
			body:   `{"title":"err"}`,
			authed: true,
			setupMock: func(m *mockJourneyClient) {
				m.createTripFn = func(ctx context.Context, in *journeyv1.CreateTripRequest, opts ...grpc.CallOption) (*journeyv1.CreateTripResponse, error) {
					return nil, status.Error(codes.Internal, "db error")
				}
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockJourneyClient{}
			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			r := setupRouter(mock)

			var bodyReader io.Reader
			if tt.body != "" {
				bodyReader = strings.NewReader(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bodyReader)
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			if tt.authed {
				req = req.WithContext(authedCtx())
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
			if tt.wantBody != "" && !strings.Contains(w.Body.String(), tt.wantBody) {
				t.Errorf("body = %q, want to contain %q", w.Body.String(), tt.wantBody)
			}
		})
	}
}
