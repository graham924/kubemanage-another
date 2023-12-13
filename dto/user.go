package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/noovertime7/kubemanage/dao/model"
	"github.com/noovertime7/kubemanage/pkg"
)

// AdminLoginInput admin登陆接口入参结构
type AdminLoginInput struct {
	// UserName 用户名
	UserName string `form:"username" json:"username" comment:"用户名"  validate:"required" example:"用户名"`
	// Password 密码
	Password string `form:"password" json:"password" comment:"密码"   validate:"required" example:"密码"`
}

// AdminLoginOut admin登陆接口出参
type AdminLoginOut struct {
	// Token token值
	Token string `form:"token" json:"token" comment:"token"  example:"token"`
}

type UserInfoOut struct {
	User      model.SysUser   `json:"user"`
	Menus     []model.SysMenu `json:"menus"`
	RuleNames []string        `json:"ruleNames"`
}

type SetUserAuth struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
}

type ChangeUserPwdInput struct {
	OldPwd string `json:"old_pwd" form:"old_pwd" comment:"原密码" validate:"required"`
	NewPwd string `json:"new_pwd" form:"new_pwd" comment:"new_pwd" validate:"required"`
}

// BindingValidParams 绑定并校验参数
func (a *ChangeUserPwdInput) BindingValidParams(ctx *gin.Context) error {
	return pkg.DefaultGetValidParams(ctx, a)
}

// BindingValidParams 绑定并校验参数
func (a *SetUserAuth) BindingValidParams(ctx *gin.Context) error {
	return pkg.DefaultGetValidParams(ctx, a)
}

// BindingValidParams AdminLoginInput绑定并校验参数
func (a *AdminLoginInput) BindingValidParams(ctx *gin.Context) error {
	return pkg.DefaultGetValidParams(ctx, a)
}
