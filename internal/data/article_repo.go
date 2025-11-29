package data

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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
	err := r.data.db.WithContext(ctx).Model(&Article{}).Create(&d).Error
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

func (r *articleRepo) UpdateArticleStatus(ctx context.Context, idList []int64, status int) (int64, error) {
	result := r.data.db.WithContext(ctx).Model(&Article{}).
		Where("id IN (?)", idList).
		Update("status", status)

	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("UpdateArticleStatus|Update fail, idList:%v, err:%v", idList, result.Error)
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (r *articleRepo) DeleteArticleById(ctx context.Context, id int64) error {
	err := r.data.db.WithContext(ctx).Model(&Article{}).Delete(&Article{Id: id}).Error
	if err != nil {
		r.log.WithContext(ctx).Errorf("DeleteArticleById|Delete fail, id:%v, err:%v", id, err)
		return err
	}
	return nil
}

func (r *articleRepo) GetArticleById(ctx context.Context, id int64) (bool, *biz.ArticleModel, error) {
	data := Article{}
	err := r.data.db.WithContext(ctx).Model(&Article{}).Where("id=?", id).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil, nil
	}
	if err != nil {
		r.log.WithContext(ctx).Errorf("GetArticleById|Get fail, id:%v, err:%v", id, err)
		return false, nil, err
	}
	m := data.Data2Biz()
	return true, m, nil
}

func (r *articleRepo) UpdateArticleById(ctx context.Context, article *biz.ArticleModel) error {
	data := Article{}
	data.Biz2Data(article)
	err := r.data.db.WithContext(ctx).Model(&Article{}).Where("id = ?", article.Id).Updates(&data).Error
	if err != nil {
		r.log.WithContext(ctx).Errorf("UpdateArticleById|Update fail, id:%v, err:%v", data.Id, err)
		return err
	}
	return nil

}
