package data

import (
	"gorm.io/gorm"
	"jike_server/internal/biz"
	"time"
)

type Account struct {
	Id          int64     `gorm:"id"`            // 用户ID
	Mobile      string    `gorm:"mobile"`        // 手机号
	Email       string    `gorm:"email"`         // 邮箱
	Password    string    `gorm:"password"`      // 密码
	Nickname    string    `gorm:"nickname"`      // 昵称
	Avatar      string    `gorm:"avatar"`        // 头像URL
	Gender      int8      `gorm:"gender"`        // 性别: 0-未知 1-男 2-女
	Birthday    time.Time `gorm:"birthday"`      // 生日
	Bio         string    `gorm:"bio"`           // 个人简介
	Status      int8      `gorm:"status"`        // 状态: 0-禁用 1-正常 2-冻结
	LastLoginAt time.Time `gorm:"last_login_at"` // 最后登录时间
	CreatedAt   time.Time `gorm:"created_at"`    // 创建时间
	UpdatedAt   time.Time `gorm:"updated_at"`    // 更新时间
}

// TableName 表名称
func (Account) TableName() string {
	return "account"
}

// BeforeCreate 创建时的钩子，设置默认值
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	// 如果生日是零值，设置默认值
	if a.Birthday.IsZero() {
		a.Birthday = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// 如果最后登录时间是零值，设置默认值
	if a.LastLoginAt.IsZero() {
		a.LastLoginAt = time.Now()
	}
	return nil
}

// 在 data 层实现转换逻辑
func (d *Account) Data2Biz() *biz.AccountModel {
	return &biz.AccountModel{
		Id:          d.Id,
		Mobile:      d.Mobile,
		Email:       d.Email,
		Password:    d.Password,
		Nickname:    d.Nickname,
		Avatar:      d.Avatar,
		Gender:      d.Gender,
		Birthday:    d.Birthday,
		Bio:         d.Bio,
		Status:      d.Status,
		LastLoginAt: d.LastLoginAt,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func (d *Account) Biz2Data(b *biz.AccountModel) {
	if d == nil {
		return
	}
	*d = Account{
		Id:          b.Id,
		Mobile:      b.Mobile,
		Email:       b.Email,
		Password:    b.Password,
		Nickname:    b.Nickname,
		Avatar:      b.Avatar,
		Gender:      b.Gender,
		Birthday:    b.Birthday,
		Bio:         b.Bio,
		Status:      b.Status,
		LastLoginAt: b.LastLoginAt,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}
