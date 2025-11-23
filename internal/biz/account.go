package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	pb "jike_server/api/account/v1"
	"jike_server/internal/middleware"
	"jike_server/internal/pkg/auth"
	"jike_server/internal/pkg/ecode"
	"time"
)

type AccountModel struct {
	Id          int64     `json:"id"`            // 用户ID
	Mobile      string    `json:"mobile"`        // 手机号
	Email       string    `json:"email"`         // 邮箱
	Password    string    `json:"password"`      // 密码
	Nickname    string    `json:"nickname"`      // 昵称
	Avatar      string    `json:"avatar"`        // 头像URL
	Gender      int8      `json:"gender"`        // 性别: 0-未知 1-男 2-女
	Birthday    time.Time `json:"birthday"`      // 生日
	Bio         string    `json:"bio"`           // 个人简介
	Status      int8      `json:"status"`        // 状态: 0-禁用 1-正常 2-冻结
	LastLoginAt time.Time `json:"last_login_at"` // 最后登录时间
	CreatedAt   time.Time `json:"created_at"`    // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`    // 更新时间
}

type AccountRepo interface {
	AddAccount(context.Context, *AccountModel) error
	GetAccountByMobile(context.Context, string) (bool, *AccountModel, error)
	GetAccountById(context.Context, int64) (bool, *AccountModel, error)
}

type AccountUsecase struct {
	repo AccountRepo
	jwt  *auth.JWT
	log  *log.Helper
}

func NewAccountUsecase(repo AccountRepo, jwt *auth.JWT, logger log.Logger) *AccountUsecase {
	return &AccountUsecase{repo: repo, jwt: jwt, log: log.NewHelper(logger)}
}

func (u *AccountUsecase) Authorizations(ctx context.Context, req *pb.AuthorizationsReq) (*pb.AuthorizationsResp, error) {
	// 1. 验证验证码
	isValid, err := u.VerifyCode(ctx, req.Mobile, req.Code)
	if err != nil {
		u.log.WithContext(ctx).Errorf("Authorizations|VerifyCode fail, req:%v, err:%v", req, err)
		return nil, err
	}
	if !isValid {
		return nil, errors.New("验证码错误")
	}

	// 2. 检查账号是否存在
	isExist, account, err := u.repo.GetAccountByMobile(ctx, req.Mobile)

	if err != nil {
		u.log.WithContext(ctx).Errorf("Authorizations|GetAccountByMobile fail, req:%v, err:%v", req, err)
		return nil, err
	}

	// 3. 如果账号不存在，则自动注册
	if !isExist {
		newAccount := &AccountModel{
			Mobile:    req.Mobile,
			Nickname:  "用户_" + req.Mobile[len(req.Mobile)-4:], // 默认昵称
			Status:    1,                                      // 正常状态
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := u.repo.AddAccount(ctx, newAccount); err != nil {
			u.log.WithContext(ctx).Errorf("Authorizations|AddAccount fail, req:%v, err:%v", req, err)
			return nil, err
		}

		// 重新获取新创建的账号信息
		isExist, account, err = u.repo.GetAccountByMobile(ctx, req.Mobile)
		if err != nil || !isExist {
			u.log.WithContext(ctx).Errorf("Authorizations|GetAccountByMobile after register fail, req:%v, err:%v", req, err)
			return nil, errors.New("注册失败")
		}
	}

	// 4. 更新最后登录时间
	// 这里可以添加更新最后登录时间的逻辑

	// 5. 生成JWT Token
	token, expiresAt, err := u.jwt.GenerateToken(account.Id, account.Mobile)
	if err != nil {
		u.log.WithContext(ctx).Errorf("Authorizations|GenerateToken fail, account:%v, err:%v", account, err)
		return nil, err
	}

	// 6. 返回响应
	return &pb.AuthorizationsResp{
		Token:     token,
		ExpiresAt: expiresAt,
		UserId:    account.Id,
		Mobile:    account.Mobile,
		Nickname:  account.Nickname,
	}, nil
}

func (u *AccountUsecase) VerifyCode(ctx context.Context, mobile, code string) (bool, error) {
	if code != "login1" {
		u.log.WithContext(ctx).Errorf("VerifyCode is not correct")
		return false, ecode.InvalidParamsWithMsg("验证码错误")
	}
	return true, nil
}

func (u *AccountUsecase) GetAccountInfo(ctx context.Context, req *pb.GetAccountReq) (*pb.GetAccountResp, error) {
	var resp *pb.GetAccountResp
	info, isExist := middleware.FromContext(ctx)
	if !isExist {
		u.log.WithContext(ctx).Errorf("GetAccountInfo|FromContext fail,req:%v", req)
		return resp, ecode.ErrUserNotFound
	}
	isExist, m, err := u.repo.GetAccountById(ctx, info.UserId)
	if err != nil {
		u.log.WithContext(ctx).Errorf("GetAccountInfo|GetAccountById fail, account:%v, err:%v", info, err)
		return resp, ecode.ErrUserNotFound
	}
	if !isExist {
		return resp, nil
	}
	return &pb.GetAccountResp{
		Birthday: m.Birthday.Format(time.DateOnly),
		Gender:   int32(m.Gender),
		Id:       m.Id,
		Mobile:   m.Mobile,
		Nickname: m.Nickname,
		Avatar:   m.Avatar,
	}, nil

}
