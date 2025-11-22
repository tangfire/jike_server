package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jike_server/internal/conf"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewAccountRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db  *gorm.DB
	rdb *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	db, err := NewDB(c, logger)
	if err != nil {
		return nil, cleanup, err
	}
	rdb, err := NewRedis(c, logger)
	if err != nil {
		return nil, cleanup, err
	}
	return &Data{db: db, rdb: rdb}, cleanup, nil
}

func NewDB(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		log.NewHelper(logger).Errorf("mysql connect failed: %v", err)
		return nil, err
	}
	return db, nil
}

func NewRedis(c *conf.Data, logger log.Logger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.Db),
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.NewHelper(logger).Errorf("redis connect failed: %v", err)
		_ = rdb.Close()
		return nil, nil
	}
	return rdb, nil
}
