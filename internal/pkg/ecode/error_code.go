package ecode

import "github.com/go-kratos/kratos/v2/errors"

var (
	// 通用错误码
	Success            = errors.New(0, "SUCCESS", "成功")
	InvalidParams      = errors.New(400, "INVALID_PARAMS", "请求参数错误")
	Unauthorized       = errors.New(401, "UNAUTHORIZED", "未授权")
	Forbidden          = errors.New(403, "FORBIDDEN", "禁止访问")
	NotFound           = errors.New(404, "NOT_FOUND", "资源不存在")
	InternalServer     = errors.New(500, "INTERNAL_SERVER", "服务器内部错误")
	ServiceUnavailable = errors.New(503, "SERVICE_UNAVAILABLE", "服务不可用")

	// 业务错误码 (从 10000 开始)
	ErrUserNotFound  = errors.New(10001, "USER_NOT_FOUND", "用户不存在")
	ErrUserExists    = errors.New(10002, "USER_EXISTS", "用户已存在")
	ErrWrongPassword = errors.New(10003, "WRONG_PASSWORD", "密码错误")
	ErrInvalidToken  = errors.New(10004, "INVALID_TOKEN", "无效的Token")
	ErrTokenExpired  = errors.New(10005, "TOKEN_EXPIRED", "Token已过期")
	ErrCodeExpired   = errors.New(10006, "CODE_EXPIRED", "验证码已过期")
	ErrCodeIncorrect = errors.New(10007, "CODE_INCORRECT", "验证码错误")
)

func InvalidParamsWithMsg(msg string) error {
	return errors.New(400, "INVALID_PARAMS", msg)
}
