package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"io"
	account_v1 "jike_server/api/account/v1"
	article_v1 "jike_server/api/article/v1"
	"jike_server/internal/conf"
	"jike_server/internal/middleware"
	"jike_server/internal/pkg/auth"
	"jike_server/internal/service"
	"os"
	"path/filepath"
	"time"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, account *service.AccountService, article *service.ArticleService, jwt *auth.JWT, logger log.Logger) *http.Server {
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
		// 注册 multipart/form-data 解码器
		//http.RequestDecoder(multipartRequestDecoder),
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
	account_v1.RegisterAccountHTTPServer(srv, account)
	article_v1.RegisterArticleHTTPServer(srv, article)

	// 添加独立的上传文件路由
	route := srv.Route("/")
	route.POST("/v1/upload", uploadFileHandler)

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

// uploadFileHandler 处理文件上传的独立HTTP处理器
func uploadFileHandler(ctx http.Context) error {
	req := ctx.Request()

	// 设置最大内存限制为10MB
	if err := req.ParseMultipartForm(10 << 20); err != nil {
		return ctx.String(400, "Failed to parse form: "+err.Error())
	}

	// 获取文件名
	fileName := req.FormValue("name")
	if fileName == "" {
		fileName = "unnamed"
	}

	// 获取文件
	file, handler, err := req.FormFile("image") // 注意这里改为 "image" 以匹配你的curl
	if err != nil {
		return ctx.String(400, "Failed to get file: "+err.Error())
	}
	defer file.Close()

	// 安全检查：限制文件大小
	if handler.Size > 10<<20 { // 10MB
		return ctx.String(400, "File too large")
	}

	// 创建上传目录
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// 安全的文件名
	timestamp := time.Now().UnixNano()
	safeFileName := fmt.Sprintf("upload_%d_%s", timestamp, handler.Filename)
	filePath := filepath.Join(uploadDir, safeFileName)

	// 保存文件
	f, err := os.Create(filePath)
	if err != nil {
		return ctx.String(500, "Failed to create file: "+err.Error())
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return ctx.String(500, "Failed to save file: "+err.Error())
	}

	// 返回成功响应
	return ctx.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "File uploaded successfully",
		"data": map[string]string{
			"url":  "/uploads/" + safeFileName,
			"name": fileName,
		},
		"timestamp": time.Now().Unix(),
	})
}
