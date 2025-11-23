package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"jike_server/internal/biz"
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

func (r articleRepo) AddArticle(ctx context.Context, article *biz.ArticleModel) (int64, error) {
	d := Article{}
	d.Biz2Data(article)
	err := r.data.db.WithContext(ctx).Model(&Article{}).Create(&article).Error
	if err != nil {
		r.log.WithContext(ctx).Errorf("AddArticle|Create fail,data:%v,err:%v", article, err)
		return 0, err
	}
	return article.Id, nil
}
