// Package external provides Google Places (Nearby) and Directions (Transit)
// proxy handlers with DynamoDB-backed response caching and daily rate limiting.
package external

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

// PlaceResult is a normalized nearby place response.
type PlaceResult struct {
	PlaceID  string   `json:"place_id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Lat      float64  `json:"lat"`
	Lng      float64  `json:"lng"`
	Rating   float64  `json:"rating"`
	Types    []string `json:"types"`
	PhotoRef string   `json:"photo_ref,omitempty"`
}

// RouteResult is a normalized transit route response.
type RouteResult struct {
	Summary  string      `json:"summary"`
	Duration string      `json:"duration"`
	Distance string      `json:"distance"`
	Steps    []RouteStep `json:"steps"`
}

// RouteStep represents one leg of a transit route.
type RouteStep struct {
	Mode         string `json:"mode"`
	Instructions string `json:"instructions"`
	Duration     string `json:"duration"`
	LineName     string `json:"line_name,omitempty"`
}

// Cache is an interface for response caching.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte, ttl time.Duration) error
}

// RateLimiter tracks daily API call count.
type RateLimiter struct {
	mu    sync.Mutex
	count int
	day   string
	limit int
}

// NewRateLimiter creates a rate limiter with the given daily limit.
func NewRateLimiter(dailyLimit int) *RateLimiter {
	return &RateLimiter{limit: dailyLimit}
}

// Allow checks if a request is allowed under the daily limit.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	today := time.Now().UTC().Format("2006-01-02")
	if rl.day != today {
		rl.day = today
		rl.count = 0
	}
	if rl.count >= rl.limit {
		return false
	}
	rl.count++
	return true
}

// Service handles external API proxy calls.
type Service struct {
	apiKey      string
	httpClient  *http.Client
	cache       Cache
	rateLimiter *RateLimiter
}

// NewService creates an external API service.
func NewService(cache Cache) *Service {
	limit, _ := strconv.Atoi(os.Getenv("EXTERNAL_API_DAILY_LIMIT"))
	if limit == 0 {
		limit = 1000
	}
	return &Service{
		apiKey:      os.Getenv("GOOGLE_MAPS_API_KEY"),
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		cache:       cache,
		rateLimiter: NewRateLimiter(limit),
	}
}

// Nearby queries Google Places Nearby Search API.
func (s *Service) Nearby(ctx context.Context, lat, lng float64, radius int, placeType string) ([]PlaceResult, error) {
	if !s.rateLimiter.Allow() {
		return nil, ErrRateLimited
	}

	cacheKey := cacheKeyNearby(lat, lng, radius, placeType)
	if s.cache != nil {
		if data, err := s.cache.Get(ctx, cacheKey); err == nil && data != nil {
			var results []PlaceResult
			if json.Unmarshal(data, &results) == nil {
				return results, nil
			}
		}
	}

	u := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%f,%f&radius=%d&type=%s&key=%s",
		lat, lng, radius, url.QueryEscape(placeType), s.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nearby: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	var raw struct {
		Results []struct {
			PlaceID  string `json:"place_id"`
			Name     string `json:"name"`
			Vicinity string `json:"vicinity"`
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
			Rating float64  `json:"rating"`
			Types  []string `json:"types"`
			Photos []struct {
				Ref string `json:"photo_reference"`
			} `json:"photos"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("nearby: parse: %w", err)
	}

	results := make([]PlaceResult, 0, len(raw.Results))
	for _, r := range raw.Results {
		pr := PlaceResult{
			PlaceID: r.PlaceID,
			Name:    r.Name,
			Address: r.Vicinity,
			Lat:     r.Geometry.Location.Lat,
			Lng:     r.Geometry.Location.Lng,
			Rating:  r.Rating,
			Types:   r.Types,
		}
		if len(r.Photos) > 0 {
			pr.PhotoRef = r.Photos[0].Ref
		}
		results = append(results, pr)
	}

	if s.cache != nil {
		data, _ := json.Marshal(results)
		_ = s.cache.Set(ctx, cacheKey, data, 30*time.Minute)
	}

	return results, nil
}

// Transit queries Google Directions API for transit mode.
func (s *Service) Transit(ctx context.Context, originLat, originLng, destLat, destLng float64, departureTime time.Time) ([]RouteResult, error) {
	if !s.rateLimiter.Allow() {
		return nil, ErrRateLimited
	}

	cacheKey := cacheKeyTransit(originLat, originLng, destLat, destLng, departureTime)
	if s.cache != nil {
		if data, err := s.cache.Get(ctx, cacheKey); err == nil && data != nil {
			var results []RouteResult
			if json.Unmarshal(data, &results) == nil {
				return results, nil
			}
		}
	}

	u := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/directions/json?origin=%f,%f&destination=%f,%f&mode=transit&departure_time=%d&key=%s",
		originLat, originLng, destLat, destLng, departureTime.Unix(), s.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("transit: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	var raw struct {
		Routes []struct {
			Summary string `json:"summary"`
			Legs    []struct {
				Duration struct {
					Text string `json:"text"`
				} `json:"duration"`
				Distance struct {
					Text string `json:"text"`
				} `json:"distance"`
				Steps []struct {
					TravelMode       string `json:"travel_mode"`
					HTMLInstructions string `json:"html_instructions"`
					Duration         struct {
						Text string `json:"text"`
					} `json:"duration"`
					TransitDetails struct {
						Line struct {
							ShortName string `json:"short_name"`
							Name      string `json:"name"`
						} `json:"line"`
					} `json:"transit_details"`
				} `json:"steps"`
			} `json:"legs"`
		} `json:"routes"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("transit: parse: %w", err)
	}

	results := make([]RouteResult, 0, len(raw.Routes))
	for _, route := range raw.Routes {
		rr := RouteResult{Summary: route.Summary}
		if len(route.Legs) > 0 {
			rr.Duration = route.Legs[0].Duration.Text
			rr.Distance = route.Legs[0].Distance.Text
			for _, step := range route.Legs[0].Steps {
				rs := RouteStep{
					Mode:         step.TravelMode,
					Instructions: step.HTMLInstructions,
					Duration:     step.Duration.Text,
				}
				if step.TransitDetails.Line.ShortName != "" {
					rs.LineName = step.TransitDetails.Line.ShortName
				} else {
					rs.LineName = step.TransitDetails.Line.Name
				}
				rr.Steps = append(rr.Steps, rs)
			}
		}
		results = append(results, rr)
	}

	if s.cache != nil {
		data, _ := json.Marshal(results)
		_ = s.cache.Set(ctx, cacheKey, data, 1*time.Hour)
	}

	return results, nil
}

// ErrRateLimited is returned when the daily API quota is exceeded.
var ErrRateLimited = fmt.Errorf("external: daily rate limit exceeded")

func cacheKeyNearby(lat, lng float64, radius int, placeType string) string {
	raw := fmt.Sprintf("nearby:%f:%f:%d:%s", lat, lng, radius, placeType)
	h := sha256.Sum256([]byte(raw))
	return "nearby:" + hex.EncodeToString(h[:8])
}

func cacheKeyTransit(oLat, oLng, dLat, dLng float64, t time.Time) string {
	// Round departure time to 15-min window for cache hits.
	rounded := t.Truncate(15 * time.Minute)
	raw := fmt.Sprintf("transit:%f:%f:%f:%f:%d", oLat, oLng, dLat, dLng, rounded.Unix())
	h := sha256.Sum256([]byte(raw))
	return "transit:" + hex.EncodeToString(h[:8])
}
