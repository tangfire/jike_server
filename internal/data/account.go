package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"jike_server/internal/biz"
)

type accountRepo struct {
	data *Data
	log  *log.Helper
}

func NewAccountRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *accountRepo) Add(ctx context.Context, a *biz.Account) error {
	return nil

}
