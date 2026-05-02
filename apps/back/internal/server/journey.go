// Package server - gRPC JourneyService 어댑터.
package server

import (
	"context"
	"errors"

	"github.com/seonNoh/seonology-journey/apps/back/internal/trip"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// JourneyServer - JourneyService 구현. 미구현 도메인은 Unimplemented 응답.
type JourneyServer struct {
	journeyv1.UnimplementedJourneyServiceServer
	Trip *trip.Service
}

// NewJourneyServer - 서버 생성.
func NewJourneyServer(tripSvc *trip.Service) *JourneyServer {
	return &JourneyServer{Trip: tripSvc}
}

// ownerFromCtx - 호출자 식별자 추출. apps/api 가 x-user-id 헤더로 전달.
func ownerFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}
	vals := md.Get("x-user-id")
	if len(vals) == 0 || vals[0] == "" {
		return "", status.Error(codes.Unauthenticated, "missing x-user-id")
	}
	return vals[0], nil
}

func mapErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, trip.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, trip.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

// CreateTrip implements JourneyService.
func (s *JourneyServer) CreateTrip(ctx context.Context, req *journeyv1.CreateTripRequest) (*journeyv1.CreateTripResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title required")
	}
	t, err := s.Trip.Create(ctx, owner, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.CreateTripResponse{Trip: t}, nil
}

// GetTrip implements JourneyService.
func (s *JourneyServer) GetTrip(ctx context.Context, req *journeyv1.GetTripRequest) (*journeyv1.GetTripResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	t, err := s.Trip.Get(ctx, owner, req.GetTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.GetTripResponse{Trip: t}, nil
}

// ListTrips implements JourneyService.
func (s *JourneyServer) ListTrips(ctx context.Context, req *journeyv1.ListTripsRequest) (*journeyv1.ListTripsResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	trips, err := s.Trip.List(ctx, owner, req.GetStatus())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListTripsResponse{Trips: trips, Page: &journeyv1.PageInfo{HasMore: false}}, nil
}

// UpdateTrip implements JourneyService.
func (s *JourneyServer) UpdateTrip(ctx context.Context, req *journeyv1.UpdateTripRequest) (*journeyv1.UpdateTripResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	t, err := s.Trip.Update(ctx, owner, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.UpdateTripResponse{Trip: t}, nil
}

// DeleteTrip implements JourneyService.
func (s *JourneyServer) DeleteTrip(ctx context.Context, req *journeyv1.DeleteTripRequest) (*journeyv1.DeleteTripResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.Trip.Delete(ctx, owner, req.GetTripId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DeleteTripResponse{}, nil
}
