// Package ws - 여행 단위 WebSocket 허브.
//
// MVP: trip_id 별 클라이언트 set 을 관리. PublishEvent 가 호출되면
// 해당 trip 의 모든 구독자에게 RealtimeEvent 의 protojson 을 송신.
package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// Hub - trip 별 구독자 관리.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]struct{} // tripID -> set
}

// New - 신규 허브.
func New() *Hub { return &Hub{clients: make(map[string]map[*Client]struct{})} }

// Client - 단일 ws 연결.
type Client struct {
	conn   *websocket.Conn
	tripID string
	userID string
}

// register - 신규 등록.
func (h *Hub) register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c.tripID]; !ok {
		h.clients[c.tripID] = make(map[*Client]struct{})
	}
	h.clients[c.tripID][c] = struct{}{}
}

// unregister - 해제.
func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.clients[c.tripID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.clients, c.tripID)
		}
	}
}

// Publish - tripID 의 모든 클라이언트에게 이벤트 송신.
func (h *Hub) Publish(ctx context.Context, ev *journeyv1.RealtimeEvent) {
	h.mu.RLock()
	set := h.clients[ev.GetTripId()]
	clients := make([]*Client, 0, len(set))
	for c := range set {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	payload, err := protojson.Marshal(ev)
	if err != nil {
		return
	}
	var raw json.RawMessage = payload
	for _, c := range clients {
		cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		_ = wsjson.Write(cctx, c.conn, raw)
		cancel()
	}
}

// Handler - HTTP 업그레이드 후 trip 구독.
//
// 사용 예: GET /ws/trips/{tripId} (auth 미들웨어로 보호).
func (h *Hub) Handler(getUserID func(*http.Request) (string, error), getTripID func(*http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserID(r)
		if err != nil || userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		tripID := getTripID(r)
		if tripID == "" {
			http.Error(w, "trip required", http.StatusBadRequest)
			return
		}
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true, // CORS 는 reverse proxy 에서 처리.
		})
		if err != nil {
			return
		}
		c := &Client{conn: conn, tripID: tripID, userID: userID}
		h.register(c)
		defer h.unregister(c)
		defer conn.Close(websocket.StatusNormalClosure, "")

		ctx := r.Context()
		// ping 루프.
		go func() {
			t := time.NewTicker(25 * time.Second)
			defer t.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
					_ = conn.Ping(pctx)
					cancel()
				}
			}
		}()
		// 클라이언트로부터의 입력 수신 (현재는 무시).
		for {
			if _, _, err := conn.Read(ctx); err != nil {
				return
			}
		}
	}
}
