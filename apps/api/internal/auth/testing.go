package auth

import "context"

// TestContext - 테스트용 인증 컨텍스트 생성.
func TestContext(ctx context.Context, uid, uname string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, uid)
	ctx = context.WithValue(ctx, usernameKey, uname)
	return ctx
}
