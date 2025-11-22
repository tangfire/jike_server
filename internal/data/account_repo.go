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

// 在 data 层实现转换逻辑
func toBizModel(d Account) *biz.AccountModel {
	return &biz.AccountModel{
		Id:          d.Id,
		Mobile:      d.Mobile,
		Email:       d.Email,
		Password:    d.Password,
		Nickname:    d.Nickname,
		Avatar:      d.Avatar,
		Gender:      d.Gender,
		Birthday:    d.Birthday,
		Bio:         d.Bio,
		Status:      d.Status,
		LastLoginAt: d.LastLoginAt,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func toDataModel(b *biz.AccountModel) Account {
	return Account{
		Id:          b.Id,
		Mobile:      b.Mobile,
		Email:       b.Email,
		Password:    b.Password,
		Nickname:    b.Nickname,
		Avatar:      b.Avatar,
		Gender:      b.Gender,
		Birthday:    b.Birthday,
		Bio:         b.Bio,
		Status:      b.Status,
		LastLoginAt: b.LastLoginAt,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

func (r *accountRepo) AddAccount(ctx context.Context, a *biz.AccountModel) error {
	// 1. 一定要将 ctx 传递给数据库操作
	data := toDataModel(a)
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
	return true, toBizModel(data), nil
}
