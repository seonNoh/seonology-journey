// Package main is the entrypoint for the seonology-journey-api REST + WebSocket gateway.
package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/seonNoh/seonology-journey/apps/api/internal/auth"
	"github.com/seonNoh/seonology-journey/apps/api/internal/external"
	"github.com/seonNoh/seonology-journey/apps/api/internal/grpcclient"
	"github.com/seonNoh/seonology-journey/apps/api/internal/handler"
	"github.com/seonNoh/seonology-journey/apps/api/internal/ws"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	addr := getenv("HTTP_LISTEN_ADDR", ":8080")
	backAddr := getenv("BACK_GRPC_ADDR", "seonology-journey-back:9090")
	jwksURL := os.Getenv("KEYCLOAK_JWKS_URL") // 공란이면 dev 모드.
	issuer := os.Getenv("KEYCLOAK_ISSUER")
	aud := os.Getenv("KEYCLOAK_AUDIENCE")

	conn, err := grpcclient.Dial(backAddr)
	if err != nil {
		log.Fatalf("grpc dial: %v", err)
	}
	defer conn.Close() //nolint:errcheck
	jc := grpcclient.NewJourneyClient(conn)

	verifier, err := auth.NewVerifier(jwksURL, issuer, aud)
	if err != nil {
		log.Fatalf("jwks: %v", err)
	}

	hub := ws.New()
	api := handler.New(jc)
	extSvc := external.NewService(external.NewMemoryCache())

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(corsMiddleware)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		_, err := jc.GetExchangeRate(ctx, &journeyv1.GetExchangeRateRequest{Base: "KRW", Target: "JPY"})
		if err != nil {
			http.Error(w, "back unreachable", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ready"))
	})

	// 보호된 라우트 (REST + WS).
	r.Group(func(r chi.Router) {
		r.Use(verifier.Middleware())
		api.Mount(r)
		extSvc.Mount(r)
		// WebSocket: 프론트는 ?token=<JWT> 형식으로 붙일 수 있도록, 헤더와 query 둘 다 지원하려면
		// 클라이언트는 Authorization 헤더를 사용. 브라우저에서 직접 붙이려면 쿠키나 subprotocol 필요.
		r.Get("/ws/trips/{tripId}", hub.Handler(
			func(req *http.Request) (string, error) { return auth.UserID(req.Context()) },
			func(req *http.Request) string { return chi.URLParam(req, "tripId") },
		))
	})

	// Internal: back 의 PublishEvent 결과를 hub 로 전달하기 위해, api 자체는
	// 별도 폴링/스트림 구독이 없는 MVP 환경에서는 외부 호출자(back)가 이 endpoint 를 호출하도록 함.
	// 인증: API_INTERNAL_TOKEN 환경변수 비교.
	internalToken := os.Getenv("API_INTERNAL_TOKEN")
	r.Post("/internal/publish", func(w http.ResponseWriter, req *http.Request) {
		if internalToken == "" || req.Header.Get("X-Internal-Token") != internalToken {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		// back gRPC PublishEvent 직접 호출 후 hub 브로드캐스트.
		var ev journeyv1.RealtimeEvent
		body, err := io.ReadAll(http.MaxBytesReader(w, req.Body, 1<<20))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(body, &ev); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := metadata.AppendToOutgoingContext(req.Context(), "x-user-id", "internal")
		_, _ = jc.PublishEvent(ctx, &journeyv1.PublishEventRequest{Event: &ev})
		hub.Publish(req.Context(), &ev)
		w.WriteHeader(http.StatusNoContent)
	})

	log.Printf("seonology-journey-api HTTP listening on %s (back=%s, jwks=%v)", addr, backAddr, jwksURL != "")
	srv := &http.Server{Addr: addr, Handler: r, ReadHeaderTimeout: 10 * time.Second}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("listen: %v", err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func corsMiddleware(next http.Handler) http.Handler {
	allow := getenv("CORS_ALLOW_ORIGIN", "*")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allow)
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Internal-Token")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; connect-src 'self' https://*.seonology.com wss://*.seonology.com; img-src 'self' https://*.amazonaws.com data:; style-src 'self' 'unsafe-inline'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
