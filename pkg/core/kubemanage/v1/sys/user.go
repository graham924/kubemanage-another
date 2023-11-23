package sys

import (
	"database/sql"
	"github.com/gin-gonic/gin"

	"github.com/noovertime7/kubemanage/dao"
	"github.com/noovertime7/kubemanage/dao/model"
	"github.com/noovertime7/kubemanage/dto"
	"github.com/noovertime7/kubemanage/pkg"
	"github.com/pkg/errors"
)

// UserServiceGetter UserService对象获取器
type UserServiceGetter interface {
	// User 获取一个新的UserService对象
	User() UserService
}

// UserService User相关api操作的Service方法
type UserService interface {
	// Login 用户登陆
	Login(ctx *gin.Context, userInfo *dto.AdminLoginInput) (string, error)
	// LoginOut 用户登出
	LoginOut(ctx *gin.Context, uid int) error
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx *gin.Context, uid int, aid uint) (*dto.UserInfoOut, error)
	// SetUserAuth 更改用户权限
	SetUserAuth(ctx *gin.Context, uid int, aid uint) error
	// DeleteUser 删除用户
	DeleteUser(ctx *gin.Context, uid int) error
	// ChangePassword 更改用户密码
	ChangePassword(ctx *gin.Context, uid int, info *dto.ChangeUserPwdInput) error
	// ResetPassword 重置用户密码
	ResetPassword(ctx *gin.Context, uid int) error
}

// userService UserService接口的实现类（在user操作相关的接口中，需要用到哪个service，就在这里内置什么Service）
type userService struct {
	// Menu userService中，需要内置一个MenuService对象，用于操作Menu
	Menu MenuService
	// Casbin userService中，需要内置一个CasbinService对象，用于操作Casbin
	Casbin CasbinService
	// factory userService中，需要内置一个db工厂对象，用于操作数据库
	factory dao.ShareDaoFactory
}

// NewUserService 新建一个 *userService 对象
func NewUserService(factory dao.ShareDaoFactory) *userService {
	return &userService{
		// factory
		factory: factory,
		//
		Menu:   NewMenuService(factory),
		Casbin: NewCasbinService(factory),
	}
}

var _ UserService = &userService{}

// Login 用户登陆
func (u *userService) Login(ctx *gin.Context, userInfo *dto.AdminLoginInput) (string, error) {
	// 从userService中，取出内置的db工厂对象，先User()获取一个user.User对象，然后调用Find查询username符合条件的User
	user, err := u.factory.User().Find(ctx, &model.SysUser{UserName: userInfo.UserName})
	if err != nil {
		// 查询出错了，返回错误
		return "", err
	}

	// 查询到了user，进行密码比对
	if !pkg.CheckPassword(userInfo.Password, user.Password) {
		return "", errors.New("密码错误,请重新输入")
	}

	// 使用JWT生成token
	token, err := pkg.JWTToken.GenerateToken(pkg.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		Username:    user.UserName,
		NickName:    user.NickName,
		AuthorityId: user.AuthorityId,
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userService) LoginOut(ctx *gin.Context, uid int) error {
	// 创建一个SysUser对象，设置 id + status=0（TODO sql.NullInt64方法仔细看看）
	user := &model.SysUser{ID: uid, Status: sql.NullInt64{Int64: 0, Valid: true}}
	// 从userService中取出db factory，从中创建一个user.User，调用User表的Updates方法，完成更新
	return u.factory.User().Updates(ctx, user)
}

func (u *userService) GetUserInfo(ctx *gin.Context, uid int, aid uint) (*dto.UserInfoOut, error) {
	// 从userService中取出db factory，从中创建一个user.User，然后根据userId查询用户
	user, err := u.factory.User().Find(ctx, &model.SysUser{ID: uid})
	if err != nil {
		return nil, err
	}
	// 从userService中取出内置的MenuService对象，根据user的authorityId，查询用户的menu列表
	menus, err := u.Menu.GetMenuByAuthorityID(ctx, aid)
	if err != nil {
		return nil, err
	}
	var outRules []string
	rules := u.Casbin.GetPolicyPathByAuthorityId(aid)
	for _, rule := range rules {
		item := rule.Path + "," + rule.Method
		outRules = append(outRules, item)
	}
	return &dto.UserInfoOut{
		User:      *user,
		Menus:     menus,
		RuleNames: outRules,
	}, nil
}

func (u *userService) SetUserAuth(ctx *gin.Context, uid int, aid uint) error {
	user := &model.SysUser{ID: uid, AuthorityId: aid}
	return u.factory.User().Updates(ctx, user)
}

func (u *userService) DeleteUser(ctx *gin.Context, uid int) error {
	user := &model.SysUser{ID: uid}
	return u.factory.User().Delete(ctx, user)
}

func (u *userService) ChangePassword(ctx *gin.Context, uid int, info *dto.ChangeUserPwdInput) error {
	userDB := &model.SysUser{ID: uid}
	user, err := u.factory.User().Find(ctx, userDB)
	if err != nil {
		return err
	}

	if !pkg.CheckPassword(info.OldPwd, user.Password) {
		return errors.New("原密码错误,请重新输入")
	}

	//生成新密码
	user.Password, err = pkg.GenSaltPassword(info.NewPwd)
	if err != nil {
		return err
	}
	return u.factory.User().Updates(ctx, user)
}

func (u *userService) ResetPassword(ctx *gin.Context, uid int) error {
	newPwd, err := pkg.GenSaltPassword("kubemanage")
	if err != nil {
		return err
	}
	user := &model.SysUser{ID: uid, Password: newPwd}
	return u.factory.User().Updates(ctx, user)
}
