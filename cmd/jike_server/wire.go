//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"jike_server/internal/biz"
	"jike_server/internal/conf"
	"jike_server/internal/data"
	"jike_server/internal/pkg/auth"
	"jike_server/internal/server"
	"jike_server/internal/service"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Auth, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, provideJWT, newApp))
}

// 提供 JWT 实例
func provideJWT(authConf *conf.Auth) *auth.JWT {
	return auth.NewJWT(authConf.Jwt.Secret, authConf.Jwt.Issuer, authConf.Jwt.ExpiresIn)
}
