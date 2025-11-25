package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"jike_server/internal/biz"
	"jike_server/internal/utils"
)

type articleRepo struct {
	data *Data
	log  *log.Helper
}

func NewArticleRepo(data *Data, logger log.Logger) biz.ArticleRepo {
	return &articleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *articleRepo) AddArticle(ctx context.Context, article *biz.ArticleModel) (int64, error) {
	d := Article{}
	d.Biz2Data(article)
	err := r.data.db.WithContext(ctx).Model(&Article{}).Create(&article).Error
	if err != nil {
		r.log.WithContext(ctx).Errorf("AddArticle|Create fail,data:%v,err:%v", article, err)
		return 0, err
	}
	return article.Id, nil
}

func (r *articleRepo) GetArticleList(ctx context.Context, p *biz.PageArticle) (int64, []*biz.ArticleModel, error) {
	db := r.data.db.WithContext(ctx).Model(&Article{})
	if p.Page == 0 || p.PageSize == 0 {
		p.Page = 1
		p.PageSize = 20
	}
	if p.ChannelId != 0 {
		db = db.Where("channel_id=?", p.ChannelId)
	}
	if p.Status != 0 {
		db = db.Where("status=?", p.Status)
	}
	if p.BeginPubdate != "" {
		db = db.Where("created_at >= ?", p.BeginPubdate)
	}
	if p.EndPubdate != "" {
		db = db.Where("created_at <= ?", p.EndPubdate)
	}
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		r.log.WithContext(ctx).Errorf("GetArticleList|Count fail,data:%v, err:%v", total, err)
		return 0, nil, err
	}
	var dataList []Article
	err = db.Offset(int((p.Page - 1) * p.PageSize)).Limit(int(p.PageSize)).Order("id DESC").Find(&dataList).Error
	if err != nil {
		r.log.WithContext(ctx).Errorf("GetArticleList|Offset fail,data:%v, err:%v", total, err)
		return 0, nil, err
	}
	retList := utils.Slice2Slice(dataList, func(d Article) *biz.ArticleModel {
		m := d.Data2Biz()
		return m
	})
	return total, retList, nil
}
