package data

import (
	"jike_server/internal/biz"
	"jike_server/internal/utils"
	"time"
)

// Article 文章表
type Article struct {
	Id           int64     `gorm:"id"`            // 文章Id
	Title        string    `gorm:"title"`         // 文章标题
	Content      string    `gorm:"content"`       // 文章内容
	CoverType    int8      `gorm:"cover_type"`    // 封面类型：1-无图，2-单图，3-多图
	CoverImages  string    `gorm:"cover_images"`  // 封面图片地址数组
	ChannelId    int64     `gorm:"channel_id"`    // 频道Id
	Status       int8      `gorm:"status"`        // 文章状态：0-草稿，1-已发布，2-已删除
	AuthorId     int64     `gorm:"author_id"`     // 作者Id
	ViewCount    int64     `gorm:"view_count"`    // 阅读量
	LikeCount    int64     `gorm:"like_count"`    // 点赞数
	CommentCount int64     `gorm:"comment_count"` // 评论数
	IsTop        int8      `gorm:"is_top"`        // 是否置顶：0-否，1-是
	IsRecommend  int8      `gorm:"is_recommend"`  // 是否推荐：0-否，1-是
	CreatedAt    time.Time `gorm:"created_at"`    // 创建时间
	UpdatedAt    time.Time `gorm:"updated_at"`    // 更新时间
}

// TableName 表名称
func (Article) TableName() string {
	return "article"
}

func (d *Article) Data2Biz() *biz.ArticleModel {
	if d == nil {
		return nil
	}
	return &biz.ArticleModel{
		Id:           d.Id,
		Title:        d.Title,
		Content:      d.Content,
		CoverType:    d.CoverType,
		CoverImages:  utils.ParseToArray(d.CoverImages),
		ChannelId:    d.ChannelId,
		Status:       d.Status,
		AuthorId:     d.AuthorId,
		ViewCount:    d.ViewCount,
		LikeCount:    d.LikeCount,
		CommentCount: d.CommentCount,
		IsTop:        d.IsTop,
		IsRecommend:  d.IsRecommend,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

func (d *Article) Biz2Data(b *biz.ArticleModel) {
	if d == nil {
		return
	}
	*d = Article{
		Id:           b.Id,
		Title:        b.Title,
		Content:      b.Content,
		CoverType:    b.CoverType,
		CoverImages:  utils.JoinToString(b.CoverImages),
		ChannelId:    b.ChannelId,
		Status:       b.Status,
		AuthorId:     b.AuthorId,
		ViewCount:    b.ViewCount,
		LikeCount:    b.LikeCount,
		CommentCount: b.CommentCount,
		IsTop:        b.IsTop,
		IsRecommend:  b.IsRecommend,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
}
