package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// === Statistics ===

// GetTripStatistics implements JourneyService.
//
// Uses a fan-out worker pool across days: each day contributes a
// schedule count, meal count, and region label. Fetches for distinct
// days run concurrently so overall latency is O(1) in the number of
// days plus a small expense/media aggregate pair, instead of O(N) as
// the previous implementation did.
func (s *JourneyServer) GetTripStatistics(ctx context.Context, req *journeyv1.GetTripStatisticsRequest) (*journeyv1.GetTripStatisticsResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	days, _ := s.d.Day.List(ctx, req.GetTripId())
	stats := &journeyv1.TripStatistics{TripId: req.GetTripId(), TotalDays: int32(len(days))}

	// Aggregations protected by a mutex; writes are small and coarse.
	var (
		mu            sync.Mutex
		scheduleTotal int32
		mealTotal     int32
		regions       = map[string]struct{}{}
	)

	// Cap concurrency so a trip with 30+ days doesn't spam DynamoDB.
	const maxParallel = 8
	sem := make(chan struct{}, maxParallel)
	var wg sync.WaitGroup

	for _, d := range days {
		d := d
		if d.GetRegion() != "" {
			mu.Lock()
			regions[d.GetRegion()] = struct{}{}
			mu.Unlock()
		}
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			scCount := int32(0)
			if scs, err := s.d.Schedule.List(ctx, d.GetId()); err == nil {
				scCount = int32(len(scs))
			}
			meCount := int32(0)
			if ms, err := s.d.Meal.List(ctx, d.GetId()); err == nil {
				meCount = int32(len(ms))
			}
			mu.Lock()
			scheduleTotal += scCount
			mealTotal += meCount
			mu.Unlock()
		}()
	}

	// Expense summary and media count run in parallel with the day fan-out.
	var (
		expenseWg      sync.WaitGroup
		expenseSummary *journeyv1.ExpenseSummary
		mediaCount     int
	)
	if s.d.Expense != nil {
		expenseWg.Add(1)
		go func() {
			defer expenseWg.Done()
			if sum, err := s.d.Expense.Summary(ctx, req.GetTripId()); err == nil {
				expenseSummary = sum
			}
		}()
	}
	if s.d.MediaRepo != nil {
		expenseWg.Add(1)
		go func() {
			defer expenseWg.Done()
			if n, err := s.d.MediaRepo.CountByTrip(ctx, req.GetTripId()); err == nil {
				mediaCount = n
			}
		}()
	}

	wg.Wait()
	expenseWg.Wait()

	stats.TotalSchedules = scheduleTotal
	stats.TotalMeals = mealTotal
	stats.VisitedRegions = int32(len(regions))
	stats.TotalPhotos = int32(mediaCount)
	if expenseSummary != nil {
		stats.TotalExpense = expenseSummary.GetGrandTotal()
	}
	return &journeyv1.GetTripStatisticsResponse{Stats: stats}, nil
}

// GetYearlyStatistics implements JourneyService - MVP stub.
func (s *JourneyServer) GetYearlyStatistics(ctx context.Context, req *journeyv1.GetYearlyStatisticsRequest) (*journeyv1.GetYearlyStatisticsResponse, error) {
	if _, err := ownerFromCtx(ctx); err != nil {
		return nil, err
	}
	return &journeyv1.GetYearlyStatisticsResponse{Stats: &journeyv1.YearlyStatistics{Year: req.GetYear()}}, nil
}

// === External (stubs) ===

// === External (geocoding via OpenStreetMap Nominatim) ===

// nominatimClient holds a shared HTTP client used for Nominatim calls.
// Nominatim의 사용 정책상 user-agent 명시와 1 req/sec 정도의 트래픽 제한이
// 권장된다. 우리 서비스에서는 사용자 입력 typeahead로 호출되므로 클라이언트
// 측에서 디바운스를 둔다.
var nominatimClient = &http.Client{Timeout: 6 * time.Second}

const nominatimUA = "seonology-journey/1.0 (+https://github.com/seonNoh/seonology-journey)"

type nominatimItem struct {
	PlaceID     int64  `json:"place_id"`
	OSMID       int64  `json:"osm_id"`
	OSMType     string `json:"osm_type"`
	DisplayName string `json:"display_name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Class       string `json:"class"`
	Address     struct {
		Name        string `json:"name"`
		Amenity     string `json:"amenity"`
		Shop        string `json:"shop"`
		Tourism     string `json:"tourism"`
		Attraction  string `json:"attraction"`
		Restaurant  string `json:"restaurant"`
		Hotel       string `json:"hotel"`
		Road        string `json:"road"`
		Suburb      string `json:"suburb"`
		City        string `json:"city"`
		Town        string `json:"town"`
		Village     string `json:"village"`
		County      string `json:"county"`
		State       string `json:"state"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		Postcode    string `json:"postcode"`
	} `json:"address"`
}

func (n *nominatimItem) bestName() string {
	candidates := []string{
		n.Name,
		n.Address.Name,
		n.Address.Amenity,
		n.Address.Shop,
		n.Address.Tourism,
		n.Address.Attraction,
		n.Address.Restaurant,
		n.Address.Hotel,
	}
	for _, c := range candidates {
		if c != "" {
			return c
		}
	}
	if n.DisplayName != "" {
		// "Cafe ABC, Road, City, Country" 형태에서 첫 토큰을 이름으로 사용.
		if i := strings.Index(n.DisplayName, ","); i > 0 {
			return strings.TrimSpace(n.DisplayName[:i])
		}
		return n.DisplayName
	}
	return ""
}

func (n *nominatimItem) toProto() *journeyv1.GeocodePlace {
	lat, _ := strconv.ParseFloat(n.Lat, 64)
	lon, _ := strconv.ParseFloat(n.Lon, 64)
	pid := fmt.Sprintf("osm:%s/%d", n.OSMType, n.OSMID)
	if n.OSMType == "" || n.OSMID == 0 {
		pid = fmt.Sprintf("nominatim:%d", n.PlaceID)
	}
	return &journeyv1.GeocodePlace{
		PlaceId:  pid,
		Name:     n.bestName(),
		Address:  n.DisplayName,
		Location: &journeyv1.GeoPoint{Latitude: lat, Longitude: lon},
	}
}

func nominatimDo(ctx context.Context, path string, q url.Values) ([]nominatimItem, error) {
	q.Set("format", "jsonv2")
	q.Set("addressdetails", "1")
	endpoint := "https://nominatim.openstreetmap.org" + path + "?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", nominatimUA)
	req.Header.Set("Accept-Language", "ja,ko,en")
	resp, err := nominatimClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("nominatim http %d", resp.StatusCode)
	}
	// /search 는 배열, /reverse 는 단일 객체이므로 호출처가 분기한다.
	if strings.HasPrefix(path, "/reverse") {
		var single nominatimItem
		if err := json.NewDecoder(resp.Body).Decode(&single); err != nil {
			return nil, err
		}
		return []nominatimItem{single}, nil
	}
	var items []nominatimItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// Geocode implements JourneyService.
//
// OpenStreetMap Nominatim 공개 엔드포인트를 사용한다. 키가 필요 없고 일본/한국
// 주소·POI 검색이 모두 가능하다. 사용량이 늘면 추후 자체 호스팅 또는
// Google Places API 로 교체한다.
func (s *JourneyServer) Geocode(ctx context.Context, req *journeyv1.GeocodeRequest) (*journeyv1.GeocodeResponse, error) {
	q := strings.TrimSpace(req.GetQuery())
	if q == "" {
		return &journeyv1.GeocodeResponse{}, nil
	}
	values := url.Values{}
	values.Set("q", q)
	values.Set("limit", "8")
	items, err := nominatimDo(ctx, "/search", values)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "geocode upstream: %v", err)
	}
	out := make([]*journeyv1.GeocodePlace, 0, len(items))
	for i := range items {
		out = append(out, items[i].toProto())
	}
	return &journeyv1.GeocodeResponse{Places: out}, nil
}

// ReverseGeocode implements JourneyService.
func (s *JourneyServer) ReverseGeocode(ctx context.Context, req *journeyv1.ReverseGeocodeRequest) (*journeyv1.ReverseGeocodeResponse, error) {
	loc := req.GetLocation()
	if loc == nil {
		return nil, status.Error(codes.InvalidArgument, "location required")
	}
	values := url.Values{}
	values.Set("lat", strconv.FormatFloat(loc.GetLatitude(), 'f', 7, 64))
	values.Set("lon", strconv.FormatFloat(loc.GetLongitude(), 'f', 7, 64))
	items, err := nominatimDo(ctx, "/reverse", values)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "reverse geocode upstream: %v", err)
	}
	if len(items) == 0 {
		return &journeyv1.ReverseGeocodeResponse{}, nil
	}
	return &journeyv1.ReverseGeocodeResponse{Place: items[0].toProto()}, nil
}

// GetExchangeRate implements JourneyService - 정적 stub (KRW->JPY=0.11 등).
func (s *JourneyServer) GetExchangeRate(_ context.Context, req *journeyv1.GetExchangeRateRequest) (*journeyv1.GetExchangeRateResponse, error) {
	rate := 1.0
	switch req.GetBase() + "->" + req.GetTarget() {
	case "KRW->JPY":
		rate = 0.11
	case "JPY->KRW":
		rate = 9.0
	case "KRW->USD":
		rate = 0.00072
	case "USD->KRW":
		rate = 1380.0
	}
	return &journeyv1.GetExchangeRateResponse{
		Base:      req.GetBase(),
		Target:    req.GetTarget(),
		Rate:      rate,
		FetchedAt: timestamppb.New(time.Now().UTC()),
	}, nil
}

// GetWeatherForecast implements JourneyService.
func (s *JourneyServer) GetWeatherForecast(_ context.Context, _ *journeyv1.GetWeatherForecastRequest) (*journeyv1.GetWeatherForecastResponse, error) {
	return nil, status.Error(codes.Unimplemented, "weather not configured")
}

// === Realtime (stub) ===

// PublishEvent implements JourneyService - hub 미연결 환경에서는 noop.
func (s *JourneyServer) PublishEvent(_ context.Context, _ *journeyv1.PublishEventRequest) (*journeyv1.PublishEventResponse, error) {
	return &journeyv1.PublishEventResponse{}, nil
}
