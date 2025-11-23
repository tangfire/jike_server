package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"jike_server/internal/biz"
	"jike_server/internal/utils"
)

type articleChannelRepo struct {
	data *Data
	log  *log.Helper
}

func NewArticleChannelRepo(data *Data, logger log.Logger) biz.ArticleChannelRepo {
	return &articleChannelRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *articleChannelRepo) GetArticleChannelList(ctx context.Context) ([]*biz.ArticleChannelModel, error) {
	var list []ArticleChannel
	err := r.data.db.WithContext(ctx).Find(&list).Error
	if err != nil {
		r.log.WithContext(ctx).Errorf("GetArticleChannelList|Find fail,data:%v,err:%v", list, err)
		return nil, err
	}
	retList := utils.Slice2Slice(list, func(d ArticleChannel) *biz.ArticleChannelModel {
		m := d.Data2Biz()
		return m
	})
	return retList, nil
}
