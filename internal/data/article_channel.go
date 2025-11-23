package data

import (
	"jike_server/internal/biz"
	"time"
)

// ArticleChannel 文章频道表
type ArticleChannel struct {
	Id        int64     `gorm:"id"`         // 频道Id
	Name      string    `gorm:"name"`       // 频道名称
	CreatedAt time.Time `gorm:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"updated_at"` // 更新时间
}

// TableName 表名称
func (ArticleChannel) TableName() string {
	return "article_channel"
}

func (d *ArticleChannel) Data2Biz() *biz.ArticleChannelModel {
	if d == nil {
		return nil
	}
	return &biz.ArticleChannelModel{
		Id:        d.Id,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (d *ArticleChannel) Biz2Data(b *biz.ArticleChannelModel) {
	if d == nil {
		return
	}
	*d = ArticleChannel{
		Id:        b.Id,
		Name:      b.Name,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
