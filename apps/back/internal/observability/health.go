package observability

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Status represents the health state of a dependency.
type Status string

const (
	StatusUp   Status = "up"
	StatusDown Status = "down"
)

// HealthChecker tracks the health status of service dependencies.
type HealthChecker struct {
	mu     sync.RWMutex
	checks map[string]Status
}

// NewHealthChecker creates a HealthChecker with no registered checks.
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{checks: make(map[string]Status)}
}

// Set updates the status of a named dependency.
func (h *HealthChecker) Set(name string, status Status) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = status
}

// IsHealthy returns true if all registered dependencies are up.
func (h *HealthChecker) IsHealthy() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, s := range h.checks {
		if s != StatusUp {
			return false
		}
	}
	return true
}

// Handler returns an http.Handler for /healthz that reports aggregate health.
func (h *HealthChecker) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.mu.RLock()
		snapshot := make(map[string]Status, len(h.checks))
		for k, v := range h.checks {
			snapshot[k] = v
		}
		h.mu.RUnlock()

		status := http.StatusOK
		if !h.IsHealthy() {
			status = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]any{
			"status": map[bool]string{true: "healthy", false: "unhealthy"}[status == http.StatusOK],
			"checks": snapshot,
		})
	})
}
