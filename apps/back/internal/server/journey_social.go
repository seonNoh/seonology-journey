package server

import (
	"context"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
)

// === Companion ===

// AddCompanion implements JourneyService.
func (s *JourneyServer) AddCompanion(ctx context.Context, req *journeyv1.AddCompanionRequest) (*journeyv1.AddCompanionResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	c, err := s.d.Companion.Add(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.AddCompanionResponse{Companion: c}, nil
}

// ListCompanions implements JourneyService.
func (s *JourneyServer) ListCompanions(ctx context.Context, req *journeyv1.ListCompanionsRequest) (*journeyv1.ListCompanionsResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	out, err := s.d.Companion.List(ctx, req.GetTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListCompanionsResponse{Companions: out}, nil
}

// UpdateCompanionRole implements JourneyService.
func (s *JourneyServer) UpdateCompanionRole(ctx context.Context, req *journeyv1.UpdateCompanionRoleRequest) (*journeyv1.UpdateCompanionRoleResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	c, err := s.d.Companion.UpdateRole(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.UpdateCompanionRoleResponse{Companion: c}, nil
}

// RemoveCompanion implements JourneyService.
func (s *JourneyServer) RemoveCompanion(ctx context.Context, req *journeyv1.RemoveCompanionRequest) (*journeyv1.RemoveCompanionResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	if err := s.d.Companion.Remove(ctx, req.GetTripId(), req.GetMemberId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.RemoveCompanionResponse{}, nil
}

// === Tag ===

// CreateTag implements JourneyService.
func (s *JourneyServer) CreateTag(ctx context.Context, req *journeyv1.CreateTagRequest) (*journeyv1.CreateTagResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	t, err := s.d.Tag.Create(ctx, owner, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.CreateTagResponse{Tag: t}, nil
}

// ListTags implements JourneyService.
func (s *JourneyServer) ListTags(ctx context.Context, _ *journeyv1.ListTagsRequest) (*journeyv1.ListTagsResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	out, err := s.d.Tag.List(ctx, owner)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListTagsResponse{Tags: out}, nil
}

// DeleteTag implements JourneyService.
func (s *JourneyServer) DeleteTag(ctx context.Context, req *journeyv1.DeleteTagRequest) (*journeyv1.DeleteTagResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	// Verify the tag belongs to the caller by attempting a direct PK read.
	// A missing row surfaces as NotFound; mismatched owner is indistinguishable
	// from a missing tag, which matches the authorization intent.
	if _, err := s.d.Tag.Get(ctx, owner, req.GetTagId()); err != nil {
		return nil, mapErr(err)
	}
	if err := s.d.Tag.Delete(ctx, owner, req.GetTagId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DeleteTagResponse{}, nil
}

// AttachTag implements JourneyService.
func (s *JourneyServer) AttachTag(ctx context.Context, req *journeyv1.AttachTagRequest) (*journeyv1.AttachTagResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	if err := s.d.Tag.Attach(ctx, owner, req.GetTripId(), req.GetTagId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.AttachTagResponse{}, nil
}

// DetachTag implements JourneyService.
func (s *JourneyServer) DetachTag(ctx context.Context, req *journeyv1.DetachTagRequest) (*journeyv1.DetachTagResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	if err := s.d.Tag.Detach(ctx, req.GetTripId(), req.GetTagId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DetachTagResponse{}, nil
}

// ListTripTags implements JourneyService.
func (s *JourneyServer) ListTripTags(ctx context.Context, req *journeyv1.ListTripTagsRequest) (*journeyv1.ListTripTagsResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	out, err := s.d.Tag.ListByTrip(ctx, req.GetTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListTripTagsResponse{Tags: out}, nil
}

// === Template ===

// SaveTripAsTemplate implements JourneyService.
func (s *JourneyServer) SaveTripAsTemplate(ctx context.Context, req *journeyv1.SaveTripAsTemplateRequest) (*journeyv1.SaveTripAsTemplateResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	t, err := s.d.Template.Save(ctx, owner, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.SaveTripAsTemplateResponse{Template: t}, nil
}

// CreateTripFromTemplate implements JourneyService.
//
// MVP: source trip 의 메타만 복제. 일정/지출 등 상세 복제는 추후 Wave 에서.
func (s *JourneyServer) CreateTripFromTemplate(ctx context.Context, req *journeyv1.CreateTripFromTemplateRequest) (*journeyv1.CreateTripFromTemplateResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	tpl, err := s.d.Template.Get(ctx, owner, req.GetTemplateId())
	if err != nil {
		return nil, mapErr(err)
	}
	src, err := s.d.Trip.Get(ctx, owner, tpl.GetSourceTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	createReq := &journeyv1.CreateTripRequest{
		Title:       req.GetTitle(),
		Description: src.GetDescription(),
		StartDate:   req.GetStartDate(),
		EndDate:     src.GetEndDate(),
		Destination: src.GetDestination(),
		CountryCode: src.GetCountryCode(),
		TotalBudget: src.GetTotalBudget(),
	}
	t, err := s.d.Trip.Create(ctx, owner, createReq)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.CreateTripFromTemplateResponse{TripId: t.GetId()}, nil
}

// ListTemplates implements JourneyService.
func (s *JourneyServer) ListTemplates(ctx context.Context, _ *journeyv1.ListTemplatesRequest) (*journeyv1.ListTemplatesResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	out, err := s.d.Template.List(ctx, owner)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListTemplatesResponse{Templates: out}, nil
}

// DeleteTemplate implements JourneyService.
func (s *JourneyServer) DeleteTemplate(ctx context.Context, req *journeyv1.DeleteTemplateRequest) (*journeyv1.DeleteTemplateResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.d.Template.Get(ctx, owner, req.GetTemplateId()); err != nil {
		return nil, mapErr(err)
	}
	if err := s.d.Template.Delete(ctx, owner, req.GetTemplateId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DeleteTemplateResponse{}, nil
}

// === FavoritePlace ===

// AddFavoritePlace implements JourneyService.
func (s *JourneyServer) AddFavoritePlace(ctx context.Context, req *journeyv1.AddFavoritePlaceRequest) (*journeyv1.AddFavoritePlaceResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	p, err := s.d.Favorite.Add(ctx, owner, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.AddFavoritePlaceResponse{Place: p}, nil
}

// ListFavoritePlaces implements JourneyService.
func (s *JourneyServer) ListFavoritePlaces(ctx context.Context, _ *journeyv1.ListFavoritePlacesRequest) (*journeyv1.ListFavoritePlacesResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	out, err := s.d.Favorite.List(ctx, owner)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListFavoritePlacesResponse{Places: out}, nil
}

// RemoveFavoritePlace implements JourneyService.
func (s *JourneyServer) RemoveFavoritePlace(ctx context.Context, req *journeyv1.RemoveFavoritePlaceRequest) (*journeyv1.RemoveFavoritePlaceResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.d.Favorite.Remove(ctx, owner, req.GetPlaceId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.RemoveFavoritePlaceResponse{}, nil
}

// === Share ===

// CreateShare implements JourneyService.
func (s *JourneyServer) CreateShare(ctx context.Context, req *journeyv1.CreateShareRequest) (*journeyv1.CreateShareResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	sh, err := s.d.Share.Create(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.CreateShareResponse{Share: sh}, nil
}

// GetShare implements JourneyService - 비인증 허용 (코드 보유자 누구나 조회 가능).
func (s *JourneyServer) GetShare(ctx context.Context, req *journeyv1.GetShareRequest) (*journeyv1.GetShareResponse, error) {
	sh, err := s.d.Share.Get(ctx, req.GetCode())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.GetShareResponse{Share: sh}, nil
}

// RevokeShare implements JourneyService.
func (s *JourneyServer) RevokeShare(ctx context.Context, req *journeyv1.RevokeShareRequest) (*journeyv1.RevokeShareResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	sh, err := s.d.Share.Get(ctx, req.GetCode())
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.requireTripOwner(ctx, owner, sh.GetTripId()); err != nil {
		return nil, err
	}
	if err := s.d.Share.Revoke(ctx, req.GetCode()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.RevokeShareResponse{}, nil
}
