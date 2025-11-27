package service

import (
	"context"
	"jike_server/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

type JobService struct {
	cron           *cron.Cron
	logger         *log.Helper
	articleUsecase *biz.ArticleUsecase // 直接依赖 biz 层
}

func NewJobService(logger log.Logger, articleUsecase *biz.ArticleUsecase) *JobService {
	return &JobService{
		cron:           cron.New(cron.WithSeconds()), // 支持秒级精度
		logger:         log.NewHelper(logger),
		articleUsecase: articleUsecase,
	}
}

func (s *JobService) Start() {
	// 每5分钟执行一次审核通过任务
	_, err := s.cron.AddFunc("0 */5 * * * *", s.reviewPassedJob)
	if err != nil {
		s.logger.Errorf("add review passed job failed: %v", err)
		return
	}

	// 每30秒执行一次
	//_, err := s.cron.AddFunc("*/30 * * * * *", s.syncDataJob)
	//if err != nil {
	//	s.logger.Errorf("add sync data job failed: %v", err)
	//	return
	//}

	// 每天凌晨1点执行
	//_, err = s.cron.AddFunc("0 0 1 * * *", s.cleanupJob)
	//if err != nil {
	//	s.logger.Errorf("add cleanup job failed: %v", err)
	//	return
	//}

	s.cron.Start()
	s.logger.Info("cron jobs started")
}

func (s *JobService) Stop() {
	if s.cron != nil {
		s.cron.Stop()
		s.logger.Info("cron jobs stopped")
	}
}

func (s *JobService) syncDataJob() {
	ctx := context.Background()
	s.logger.Infof("start sync data job at %v", time.Now())

	// 执行具体的业务逻辑
	err := s.doSyncData(ctx)
	if err != nil {
		s.logger.Errorf("sync data job failed: %v", err)
	}

	s.logger.Infof("finish sync data job at %v", time.Now())
}

func (s *JobService) cleanupJob() {
	ctx := context.Background()
	s.logger.Infof("start cleanup job at %v", time.Now())

	// 执行清理逻辑
	err := s.doCleanup(ctx)
	if err != nil {
		s.logger.Errorf("cleanup job failed: %v", err)
	}

	s.logger.Infof("finish cleanup job at %v", time.Now())
}

func (s *JobService) doSyncData(ctx context.Context) error {
	// 实现数据同步逻辑
	return nil
}

func (s *JobService) doCleanup(ctx context.Context) error {
	// 实现清理逻辑
	return nil
}

// 新增：每5分钟执行的审核通过任务
func (s *JobService) reviewPassedJob() {
	ctx := context.Background()
	s.logger.Infof("start review passed job at %v", time.Now())

	err := s.doReviewPassed(ctx)
	if err != nil {
		s.logger.Errorf("review passed job failed: %v", err)
	}

	s.logger.Infof("finish review passed job at %v", time.Now())
}

func (s *JobService) doReviewPassed(ctx context.Context) error {
	err := s.articleUsecase.DoReviewPassed(ctx)
	if err != nil {
		s.logger.Errorf("do review passed failed: %v", err)
		return err
	}
	s.logger.Infof("do review passed at %v", time.Now())
	return nil
}
