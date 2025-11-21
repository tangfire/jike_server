package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type Account struct {
	Mobile string
	Code   string
}

type AccountRepo interface {
	Add(context.Context, *Account) error
}

type AccountUsecase struct {
	repo AccountRepo
	log  *log.Helper
}

func NewAccountUsecase(repo AccountRepo, logger log.Logger) *AccountUsecase {
	return &AccountUsecase{repo: repo, log: log.NewHelper(logger)}
}
