package server

import (
	"context"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
)

// === Media ===

// GetUploadUrl implements JourneyService.
func (s *JourneyServer) GetUploadUrl(ctx context.Context, req *journeyv1.GetUploadUrlRequest) (*journeyv1.GetUploadUrlResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	resp, err := s.d.Media.GetUploadURL(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return resp, nil
}

// ConfirmUpload implements JourneyService.
func (s *JourneyServer) ConfirmUpload(ctx context.Context, req *journeyv1.ConfirmUploadRequest) (*journeyv1.ConfirmUploadResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	m, err := s.d.Media.ConfirmUpload(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ConfirmUploadResponse{Media: m}, nil
}

// ListMedia implements JourneyService.
func (s *JourneyServer) ListMedia(ctx context.Context, req *journeyv1.ListMediaRequest) (*journeyv1.ListMediaResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	cursor := ""
	limit := int32(0)
	if req.GetPage() != nil {
		cursor = req.GetPage().GetCursor()
		limit = req.GetPage().GetLimit()
	}
	out, next, err := s.d.Media.ListPage(ctx, req.GetTripId(), req.GetDayId(), cursor, limit)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListMediaResponse{
		Items: out,
		Page:  &journeyv1.PageInfo{NextCursor: next, HasMore: next != ""},
	}, nil
}

// DeleteMedia implements JourneyService.
func (s *JourneyServer) DeleteMedia(ctx context.Context, req *journeyv1.DeleteMediaRequest) (*journeyv1.DeleteMediaResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	prev, err := s.d.Media.Get(ctx, req.GetMediaId())
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.requireTripOwner(ctx, owner, prev.GetTripId()); err != nil {
		return nil, err
	}
	if err := s.d.Media.Delete(ctx, req.GetMediaId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DeleteMediaResponse{}, nil
}

// GetMediaUrl implements JourneyService.
func (s *JourneyServer) GetMediaUrl(ctx context.Context, req *journeyv1.GetMediaUrlRequest) (*journeyv1.GetMediaUrlResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	prev, err := s.d.Media.Get(ctx, req.GetMediaId())
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.requireTripOwner(ctx, owner, prev.GetTripId()); err != nil {
		return nil, err
	}
	return s.d.Media.URL(ctx, req.GetMediaId(), req.GetThumbnail())
}
