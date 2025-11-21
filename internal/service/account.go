package service

import (
	"context"

	pb "jike_server/api/account/v1"
)

type AccountService struct {
	pb.UnimplementedAccountServer
}

func NewAccountService() *AccountService {
	return &AccountService{}
}

func (s *AccountService) Authorizations(ctx context.Context, req *pb.AuthorizationsReq) (*pb.AuthorizationsResp, error) {
	return &pb.AuthorizationsResp{}, nil
}
