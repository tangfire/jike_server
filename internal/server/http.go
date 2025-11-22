package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	v1 "jike_server/api/account/v1"
	"jike_server/internal/conf"
	"jike_server/internal/middleware"
	"jike_server/internal/pkg/auth"
	"jike_server/internal/service"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, account *service.AccountService, jwt *auth.JWT, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			selector.Server(
				middleware.Auth(jwt), // 使用传入的 jwt 实例，而不是新建
			).
				Match(whiteListMatcher()).
				Build(),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
			handlers.ExposedHeaders([]string{"Content-Length"}),
			handlers.AllowCredentials(),
			handlers.MaxAge(3600),
		)),
		// 使用自定义的编码器
		http.ResponseEncoder(middleware.ResponseEncoder),
		http.ErrorEncoder(middleware.ErrorEncoder),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterAccountHTTPServer(srv, account)
	return srv
}

// whiteListMatcher 白名单匹配器
func whiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]bool)
	// /包名.服务名/方法名
	// 根据你的 proto 文件，服务名是 account.v1.Account
	whiteList["/account.v1.Account/Authorizations"] = true
	// 如果需要，可以添加其他公开接口
	// whiteList["/account.v1.Account/Register"] = true

	return func(ctx context.Context, operation string) bool {
		// 如果在白名单中，返回 false 表示不需要认证
		if _, ok := whiteList[operation]; ok {
			return false
		}
		// 不在白名单中的接口需要认证
		return true
	}
}
