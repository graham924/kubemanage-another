package v1

import (
	"github.com/noovertime7/kubemanage/dao"
	"github.com/noovertime7/kubemanage/pkg/core/kubemanage/v1/sys"
)

// SystemGetter SystemInterface对象获取器
type SystemGetter interface {
	// System 创建一个新的SystemInterface对象
	System() SystemInterface
}

// SystemInterface 顶层抽象 组合了系统操作相关的多个接口
type SystemInterface interface {
	// UserServiceGetter 用户操作相关接口
	sys.UserServiceGetter
	// MenuGetter 菜单相关接口
	sys.MenuGetter
	// CasbinServiceGetter Casbin相关接口
	sys.CasbinServiceGetter
	// AuthorityGetter 权限相关接口
	sys.AuthorityGetter
	// OperationServiceGetter 操作记录相关接口
	sys.OperationServiceGetter
	// APIServiceGetter API服务相关接口
	sys.APIServiceGetter
}

var _ SystemInterface = &system{}

// NewSystem 新建一个SystemInterface对象
func NewSystem(app *KubeManage) SystemInterface {
	return &system{app: app, factory: app.Factory}
}

// system SystemInterface接口的实现类
type system struct {
	app     *KubeManage
	factory dao.ShareDaoFactory
}

// User 获取一个 sys.UserService 对象，此方法为实现 sys.UserServiceGetter 接口方法
func (s *system) User() sys.UserService {
	// sys.NewUserService 返回的是 *userService 对象，userService 是 UserService 的实现类，所以返回userService对象也可以
	return sys.NewUserService(s.factory)
}

// Menu 获取一个 sys.MenuGetter 对象，此方法为实现 sys.MenuGetter 接口方法
func (s *system) Menu() sys.MenuService {
	return sys.NewMenuService(s.factory)
}

// CasbinService 获取一个 sys.CasbinService 对象，此方法为实现 sys.CasbinServiceGetter 接口方法
func (s *system) CasbinService() sys.CasbinService {
	return sys.NewCasbinService(s.factory)
}

// Authority 获取一个 sys.Authority 对象，此方法为实现 sys.AuthorityGetter 接口方法
func (s *system) Authority() sys.Authority {
	return sys.NewAuthority(s.factory)
}

// Operation 获取一个 sys.OperationService 对象，此方法为实现 sys.OperationServiceGetter 接口方法
func (s *system) Operation() sys.OperationService {
	return sys.NewOperationService(s.factory)
}

// Api 获取一个 sys.APIService 对象，此方法为实现 sys.APIServiceGetter 接口方法
func (s *system) Api() sys.APIService {
	return sys.NewApiService(s.factory)
}
