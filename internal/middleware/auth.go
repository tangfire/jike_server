package middleware

import (
	"context"
	"jike_server/internal/pkg/auth"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// AuthKey 用于从 context 中获取认证信息
type AuthKey struct{}

// AuthInfo 认证信息
type AuthInfo struct {
	UserId int64
	Mobile string
}

func Auth(jwtAuth *auth.JWT) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 从传输层获取信息
			if tr, ok := transport.FromServerContext(ctx); ok {
				// HTTP 传输
				if ht, ok := tr.(*http.Transport); ok {
					authHeader := ht.RequestHeader().Get("Authorization")
					if authHeader == "" {
						return nil, auth.ErrTokenInvalid
					}

					// 解析 Bearer Token
					parts := strings.SplitN(authHeader, " ", 2)
					if !(len(parts) == 2 && parts[0] == "Bearer") {
						return nil, auth.ErrTokenInvalid
					}

					// 验证 Token
					claims, err := jwtAuth.ParseToken(parts[1])
					if err != nil {
						return nil, err
					}

					// 将用户信息存入上下文
					authInfo := &AuthInfo{
						UserId: claims.UserId,
						Mobile: claims.Mobile,
					}
					ctx = context.WithValue(ctx, AuthKey{}, authInfo)
				}

				// 如果是 gRPC 传输，可以类似处理 metadata
			}

			return handler(ctx, req)
		}
	}
}

// FromContext 从上下文中获取认证信息
func FromContext(ctx context.Context) (*AuthInfo, bool) {
	authInfo, ok := ctx.Value(AuthKey{}).(*AuthInfo)
	return authInfo, ok
}

// MustFromContext 从上下文中获取认证信息，如果没有则 panic
func MustFromContext(ctx context.Context) *AuthInfo {
	authInfo, ok := ctx.Value(AuthKey{}).(*AuthInfo)
	if !ok {
		panic("auth info not found in context")
	}
	return authInfo
}
