package external

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// Mount registers the /api/v1/nearby and /api/v1/transit routes.
func (s *Service) Mount(r chi.Router) {
	r.Get("/api/v1/nearby", s.handleNearby)
	r.Get("/api/v1/transit", s.handleTransit)
}

func (s *Service) handleNearby(w http.ResponseWriter, r *http.Request) {
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		http.Error(w, "lat required", http.StatusBadRequest)
		return
	}
	lng, err := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	if err != nil {
		http.Error(w, "lng required", http.StatusBadRequest)
		return
	}
	radius, _ := strconv.Atoi(r.URL.Query().Get("radius"))
	if radius == 0 {
		radius = 1000
	}
	placeType := r.URL.Query().Get("type")
	if placeType == "" {
		placeType = "restaurant"
	}

	results, err := s.Nearby(r.Context(), lat, lng, radius, placeType)
	if err == ErrRateLimited {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"results": results})
}

func (s *Service) handleTransit(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	originLat, err := strconv.ParseFloat(q.Get("origin_lat"), 64)
	if err != nil {
		http.Error(w, "origin_lat required", http.StatusBadRequest)
		return
	}
	originLng, err := strconv.ParseFloat(q.Get("origin_lng"), 64)
	if err != nil {
		http.Error(w, "origin_lng required", http.StatusBadRequest)
		return
	}
	destLat, err := strconv.ParseFloat(q.Get("dest_lat"), 64)
	if err != nil {
		http.Error(w, "dest_lat required", http.StatusBadRequest)
		return
	}
	destLng, err := strconv.ParseFloat(q.Get("dest_lng"), 64)
	if err != nil {
		http.Error(w, "dest_lng required", http.StatusBadRequest)
		return
	}

	var departure time.Time
	if ts := q.Get("departure_time"); ts != "" {
		if unix, err := strconv.ParseInt(ts, 10, 64); err == nil {
			departure = time.Unix(unix, 0)
		}
	}
	if departure.IsZero() {
		departure = time.Now()
	}

	results, err := s.Transit(r.Context(), originLat, originLng, destLat, destLng, departure)
	if err == ErrRateLimited {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"routes": results})
}
