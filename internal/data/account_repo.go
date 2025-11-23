package data

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"jike_server/internal/biz"
)

type accountRepo struct {
	data *Data
	log  *log.Helper
}

func NewAccountRepo(data *Data, logger log.Logger) biz.AccountRepo {
	return &accountRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *accountRepo) AddAccount(ctx context.Context, a *biz.AccountModel) error {
	// 1. 一定要将 ctx 传递给数据库操作
	data := Account{}
	data.Biz2Data(a)
	err := r.data.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		// 2. 使用带上下文的日志记录
		r.log.WithContext(ctx).Errorf("AddAccount|Create fail, data:%v, err:%v", data, err)
		return err
	}
	return nil
}

func (r *accountRepo) GetAccountByMobile(ctx context.Context, mobile string) (bool, *biz.AccountModel, error) {
	data := Account{}
	err := r.data.db.WithContext(ctx).Where("mobile = ?", mobile).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil, nil
	}
	if err != nil {
		r.log.WithContext(ctx).Errorf("GetAccountByMobile|First fail, data:%v, err:%v", data, err)
		return false, nil, err
	}
	m := data.Data2Biz()
	return true, m, nil
}

func (r *accountRepo) GetAccountById(ctx context.Context, id int64) (bool, *biz.AccountModel, error) {
	data := Account{}
	err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil, nil
	}
	if err != nil {
		r.log.WithContext(ctx).Errorf("GetAccountById|First fail, data:%v, err:%v", data, err)
		return false, nil, err
	}
	m := data.Data2Biz()
	return true, m, nil
}
