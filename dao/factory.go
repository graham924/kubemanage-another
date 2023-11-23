package dao

import (
	"gorm.io/gorm"

	"github.com/noovertime7/kubemanage/dao/api"
	"github.com/noovertime7/kubemanage/dao/authority"
	"github.com/noovertime7/kubemanage/dao/menu"
	"github.com/noovertime7/kubemanage/dao/operation"
	"github.com/noovertime7/kubemanage/dao/user"
	"github.com/noovertime7/kubemanage/dao/workflow"
)

// ShareDaoFactory 数据库抽象工厂 包含所有数据操作接口
type ShareDaoFactory interface {
	GetDB() *gorm.DB
	WorkFlow() workflow.WorkFlowInterface
	// User 创建一个 db的User 对象
	User() user.User
	Api() api.APi
	Authority() authority.Authority
	AuthorityMenu() authority.AuthorityMenu
	BaseMenu() menu.BaseMenu
	Opera() operation.Operation
}

func NewShareDaoFactory(db *gorm.DB) ShareDaoFactory {
	return &shareDaoFactory{db: db}
}

var _ ShareDaoFactory = &shareDaoFactory{}

// shareDaoFactory ShareDaoFactory接口实现类
type shareDaoFactory struct {
	// TODO 每个接口的实现类都有*gorm.DB，这是什么
	db *gorm.DB
}

func (s *shareDaoFactory) GetDB() *gorm.DB {
	return s.db
}

func (s *shareDaoFactory) WorkFlow() workflow.WorkFlowInterface {
	return workflow.NewWorkFlow(s.db)
}

// User 创建一个 user.User 对象
func (s *shareDaoFactory) User() user.User {
	return user.NewUser(s.db)
}

func (s *shareDaoFactory) Api() api.APi {
	return api.NewApi(s.db)
}

func (s *shareDaoFactory) Authority() authority.Authority {
	return authority.NewAuthority(s.db)
}

func (s *shareDaoFactory) AuthorityMenu() authority.AuthorityMenu {
	return authority.NewAuthorityMenu(s.db)
}

func (s *shareDaoFactory) BaseMenu() menu.BaseMenu {
	return menu.NewBaseMenu(s.db)
}

func (s *shareDaoFactory) Opera() operation.Operation {
	return operation.NewOperation(s.db)
}
