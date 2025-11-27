package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "jike_server/api/article/v1"
	"jike_server/internal/middleware"
	"jike_server/internal/pkg/ecode"
	"jike_server/internal/utils"
	"time"
)

var (
	MaxFileSize  = 10 * 1024 * 1024 // 10MB
	AllowedTypes = map[string]bool{"image/jpeg": true, "image/png": true, "image/gif": true, "image/webp": true}
	UploadDir    = "./uploads"
)

type ArticleModel struct {
	Id           int64     `json:"id"`            // 文章Id
	Title        string    `json:"title"`         // 文章标题
	Content      string    `json:"content"`       // 文章内容
	CoverType    int8      `json:"cover_type"`    // 封面类型：1-无图，2-单图，3-多图
	CoverImages  []string  `json:"cover_images"`  // 封面图片地址数组
	ChannelId    int64     `json:"channel_id"`    // 频道Id
	Status       int8      `json:"status"`        // 文章状态：0-草稿，1-已发布，2-已删除
	AuthorId     int64     `json:"author_id"`     // 作者Id
	ViewCount    int64     `json:"view_count"`    // 阅读量
	LikeCount    int64     `json:"like_count"`    // 点赞数
	CommentCount int64     `json:"comment_count"` // 评论数
	IsTop        int8      `json:"is_top"`        // 是否置顶：0-否，1-是
	IsRecommend  int8      `json:"is_recommend"`  // 是否推荐：0-否，1-是
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`    // 更新时间
}

type PageArticle struct {
	Status       int64  `json:"status"`
	ChannelId    int64  `json:"channel_id"`
	BeginPubdate string `json:"begin_pubdate"`
	EndPubdate   string `json:"end_pubdate"`
	Page         int64  `json:"page"`
	PageSize     int64  `json:"page_size"`
}

type ArticleRepo interface {
	AddArticle(ctx context.Context, article *ArticleModel) (int64, error)
	GetArticleList(ctx context.Context, p *PageArticle) (int64, []*ArticleModel, error)
	UpdateArticleStatus(ctx context.Context, idList []int64, status int) (int64, error)
}

type ArticleUsecase struct {
	repo ArticleRepo
	log  *log.Helper
}

func NewArticleUsecase(repo ArticleRepo, logger log.Logger) *ArticleUsecase {
	return &ArticleUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (u *ArticleUsecase) AddArticle(ctx context.Context, req *pb.AddArticleReq) (*pb.AddArticleResp, error) {
	info, isExist := middleware.FromContext(ctx)
	if !isExist {
		return nil, ecode.ErrUserNotFound
	}
	id, err := u.repo.AddArticle(ctx, &ArticleModel{
		Title:        req.GetTitle(),
		Content:      req.GetContent(),
		CoverType:    int8(req.GetCover().GetType()),
		CoverImages:  req.GetCover().GetImages(),
		ChannelId:    req.GetChannelId(),
		Status:       int8(utils.If(req.GetDraft() == "true", 1, 2)),
		AuthorId:     info.UserId,
		ViewCount:    0,
		LikeCount:    0,
		CommentCount: 0,
		IsTop:        0,
		IsRecommend:  0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		u.log.WithContext(ctx).Errorf("AddArticle|AddArticle fail,req:%v,err:%v", req, err)
		return nil, err
	}
	return &pb.AddArticleResp{
		Id: id,
	}, nil
}

func (u *ArticleUsecase) GetArticleList(ctx context.Context, req *pb.GetArticleListReq) (*pb.GetArticleListResp, error) {
	p := PageArticle{
		Status:       req.GetStatus(),
		ChannelId:    req.GetChannelId(),
		BeginPubdate: req.GetBeginPubdate(),
		EndPubdate:   req.GetEndPubdate(),
		Page:         req.GetPage(),
		PageSize:     req.GetPageSize(),
	}
	total, list, err := u.repo.GetArticleList(ctx, &p)
	if err != nil {
		u.log.WithContext(ctx).Errorf("GetArticleList|GetArticleList fail,req:%v,err:%v", req, err)
		return nil, err
	}
	retList := utils.Slice2Slice(list, func(m *ArticleModel) *pb.ArticleInfo {
		return &pb.ArticleInfo{
			Id:           m.Id,
			Title:        m.Title,
			Content:      m.Content,
			CoverType:    int64(m.CoverType),
			CoverImages:  m.CoverImages,
			ChannelId:    m.ChannelId,
			Status:       int64(m.Status),
			AuthorId:     m.AuthorId,
			ViewCount:    m.ViewCount,
			LikeCount:    m.LikeCount,
			CommentCount: m.CommentCount,
			IsTop:        int64(m.IsTop),
			IsRecommend:  int64(m.IsRecommend),
			CreatedAt:    timestamppb.New(m.CreatedAt),
			UpdatedAt:    timestamppb.New(m.UpdatedAt),
		}
	})
	return &pb.GetArticleListResp{
		Total: total,
		List:  retList,
	}, nil
}

func (u *ArticleUsecase) DoReviewPassed(ctx context.Context) error {
	_, list, err := u.repo.GetArticleList(ctx, &PageArticle{
		Status: 1,
	})
	if err != nil {
		u.log.WithContext(ctx).Errorf("doReviewPassed|GetArticleList fail,err:%v", err)
		return err
	}
	idList := utils.Slice2Slice(list, func(m *ArticleModel) int64 {
		return m.Id
	})
	_, err = u.repo.UpdateArticleStatus(ctx, idList, 2)
	if err != nil {
		u.log.WithContext(ctx).Errorf("doReviewPassed|UpdateArticleStatus fail,err:%v", err)
		return err
	}
	return nil
}
