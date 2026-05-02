// Package auth - Keycloak JWT 검증 미들웨어.
package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// ctxKey - context key 타입.
type ctxKey string

const (
	userIDKey   ctxKey = "uid"
	usernameKey ctxKey = "uname"
)

// Verifier - JWT 검증기.
type Verifier struct {
	jwks    keyfunc.Keyfunc
	issuer  string
	aud     string
	disable bool
}

// NewVerifier - JWKS URL 로 keyfunc 로딩.
func NewVerifier(jwksURL, issuer, aud string) (*Verifier, error) {
	if jwksURL == "" {
		return &Verifier{disable: true, issuer: issuer, aud: aud}, nil
	}
	k, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		return nil, err
	}
	return &Verifier{jwks: k, issuer: issuer, aud: aud}, nil
}

// Middleware - Authorization: Bearer <token> 검증.
func (v *Verifier) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Health / WS 진입은 별도 체인. 이 미들웨어는 보호 대상에만 attach.
			tok := bearerToken(r)
			if tok == "" {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			if v.disable {
				// 개발 모드: sub 만 추출.
				claims := jwt.MapClaims{}
				_, _, err := jwt.NewParser().ParseUnverified(tok, claims)
				if err != nil {
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}
				ctx := withClaims(r.Context(), claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			parsed, err := jwt.Parse(tok, v.jwks.Keyfunc, jwt.WithLeeway(30*time.Second))
			if err != nil || !parsed.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			claims, ok := parsed.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid claims", http.StatusUnauthorized)
				return
			}
			if v.issuer != "" {
				iss, _ := claims["iss"].(string)
				if iss != v.issuer {
					http.Error(w, "invalid issuer", http.StatusUnauthorized)
					return
				}
			}
			if v.aud != "" && !audMatches(claims["aud"], v.aud) {
				http.Error(w, "invalid audience", http.StatusUnauthorized)
				return
			}
			ctx := withClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
		return ""
	}
	return strings.TrimSpace(h[7:])
}

func audMatches(a interface{}, want string) bool {
	switch v := a.(type) {
	case string:
		return v == want
	case []interface{}:
		for _, x := range v {
			if s, ok := x.(string); ok && s == want {
				return true
			}
		}
	}
	return false
}

func withClaims(ctx context.Context, claims jwt.MapClaims) context.Context {
	uid, _ := claims["sub"].(string)
	uname, _ := claims["preferred_username"].(string)
	ctx = context.WithValue(ctx, userIDKey, uid)
	ctx = context.WithValue(ctx, usernameKey, uname)
	return ctx
}

// UserID - 인증된 사용자 ID 반환.
func UserID(ctx context.Context) (string, error) {
	if v, _ := ctx.Value(userIDKey).(string); v != "" {
		return v, nil
	}
	return "", errors.New("auth: no user")
}

// Username - 인증된 사용자명.
func Username(ctx context.Context) string {
	v, _ := ctx.Value(usernameKey).(string)
	return v
}
