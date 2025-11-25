package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "jike_server/api/article/v1"
	"jike_server/internal/biz"
)

type ArticleService struct {
	pb.UnimplementedArticleServer
	log       *log.Helper
	channelUc *biz.ArticleChannelUsecase
	uc        *biz.ArticleUsecase
}

func NewArticleService(uc *biz.ArticleUsecase, channelUc *biz.ArticleChannelUsecase, logger log.Logger) *ArticleService {
	return &ArticleService{uc: uc, channelUc: channelUc, log: log.NewHelper(logger)}
}

func (s *ArticleService) GetArticleChannel(ctx context.Context, req *pb.GetArticleChannelReq) (*pb.GetArticleChannelResp, error) {
	resp, err := s.channelUc.GetArticleChannelList(ctx)
	if err != nil {
		s.log.WithContext(ctx).Errorf("get article channel list failed: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *ArticleService) AddArticle(ctx context.Context, req *pb.AddArticleReq) (*pb.AddArticleResp, error) {
	resp, err := s.uc.AddArticle(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("add article failed: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *ArticleService) GetArticleList(ctx context.Context, req *pb.GetArticleListReq) (*pb.GetArticleListResp, error) {
	resp, err := s.uc.GetArticleList(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("get article list failed: %v", err)
		return nil, err
	}
	return resp, nil
}
