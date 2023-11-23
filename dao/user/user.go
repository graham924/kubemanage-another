package user

import (
	"context"
	"github.com/noovertime7/kubemanage/dao/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// User User的db操作接口
type User interface {
	// Find 根据条件查找SysUser
	Find(ctx context.Context, userInfo *model.SysUser) (*model.SysUser, error)
	// Save 保存user
	Save(ctx context.Context, userInfo *model.SysUser) error
	// Updates 更改user信息
	Updates(ctx context.Context, userInfo *model.SysUser) error
	// Delete 删除user
	Delete(ctx context.Context, userInfo *model.SysUser) error
}

var _ User = &user{}

// user User接口的实现类
type user struct {
	db *gorm.DB
}

// NewUser 新建一个User对象，user是User接口的实现类
func NewUser(db *gorm.DB) *user {
	// 新建一个user对象
	return &user{db: db}
}

// Find 根据条件查找user
func (u *user) Find(ctx context.Context, userInfo *model.SysUser) (*model.SysUser, error) {
	user := &model.SysUser{}
	if err := u.db.WithContext(ctx).Preload("Authorities").Preload("Authority").Where(userInfo).Find(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return user, nil
}

// Save 保存user
func (u *user) Save(ctx context.Context, userInfo *model.SysUser) error {
	return u.db.WithContext(ctx).Save(userInfo).Error
}

// Updates 更改user信息
func (u *user) Updates(ctx context.Context, userInfo *model.SysUser) error {
	if userInfo.ID == 0 {
		return errors.New("id not set")
	}
	return u.db.WithContext(ctx).Updates(userInfo).Error
}

// Delete 删除user
func (u *user) Delete(ctx context.Context, userInfo *model.SysUser) error {
	return u.db.WithContext(ctx).Delete(userInfo).Error
}
