package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "jike_server/api/article/v1"
	"jike_server/internal/utils"
	"time"
)

type ArticleChannelModel struct {
	Id        int64     `json:"id"`         // 频道Id
	Name      string    `json:"name"`       // 频道名称
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

type ArticleChannelRepo interface {
	GetArticleChannelList(ctx context.Context) ([]*ArticleChannelModel, error)
}

type ArticleChannelUsecase struct {
	repo ArticleChannelRepo
	log  *log.Helper
}

func NewArticleChannelUsecase(repo ArticleChannelRepo, logger log.Logger) *ArticleChannelUsecase {
	return &ArticleChannelUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (u *ArticleChannelUsecase) GetArticleChannelList(ctx context.Context) (*pb.GetArticleChannelResp, error) {

	list, err := u.repo.GetArticleChannelList(ctx)
	if err != nil {
		u.log.WithContext(ctx).Errorf("GetArticleChannelList|GetArticleChannelList fail,data:%v,err:%v", list, err)
		return nil, err
	}

	retList := utils.Slice2Slice(list, func(m *ArticleChannelModel) *pb.ArticleChannel {
		return &pb.ArticleChannel{
			Id:   m.Id,
			Name: m.Name,
		}
	})
	return &pb.GetArticleChannelResp{
		Channels: retList,
	}, nil

}
