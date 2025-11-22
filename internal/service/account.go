package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "jike_server/api/account/v1"
	"jike_server/internal/biz"
)

type AccountService struct {
	pb.UnimplementedAccountServer
	log *log.Helper
	uc  *biz.AccountUsecase
}

func NewAccountService(uc *biz.AccountUsecase, logger log.Logger) *AccountService {
	return &AccountService{uc: uc, log: log.NewHelper(logger)}
}

func (s *AccountService) Authorizations(ctx context.Context, req *pb.AuthorizationsReq) (*pb.AuthorizationsResp, error) {
	s.log.WithContext(ctx).Infof("Authorizations request: mobile=%s", req.Mobile)

	resp, err := s.uc.Authorizations(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Authorizations failed: %v", err)
		// 直接返回错误，中间件会统一处理
		return nil, err
	}

	s.log.WithContext(ctx).Infof("Authorizations success: userId=%d, mobile=%s", resp.UserId, resp.Mobile)
	return resp, nil
}
