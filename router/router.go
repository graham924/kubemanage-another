package router

import (
	"github.com/gin-gonic/gin"
	"github.com/noovertime7/kubemanage/controller/api"
	"github.com/noovertime7/kubemanage/controller/operation"
	"github.com/noovertime7/kubemanage/controller/user"

	"github.com/noovertime7/kubemanage/cmd/app/options"
	"github.com/noovertime7/kubemanage/controller/authority"
	"github.com/noovertime7/kubemanage/controller/kubeController"
	"github.com/noovertime7/kubemanage/controller/menu"
	"github.com/noovertime7/kubemanage/controller/other"
	"github.com/noovertime7/kubemanage/middleware"
)

// InstallRouters 初始化 APIs 路由
func InstallRouters(opt *options.Options) {
	// 创建一个apiGroup，前缀是/api
	apiGroup := opt.GinEngine.Group("/api")
	// 安装中间件组件
	middleware.InstallMiddlewares(apiGroup)
	// 安装 不需要记录操作历史 的路由
	installUnOperationRouters(apiGroup)
	// 安装 需要记录操作历史 的路由，路由前都加上了一个中间件
	installOperationRouters(apiGroup)
}

// installUnOperationRouters 本方法里的api路由，不需要记录操作历史
func installUnOperationRouters(apiGroup *gin.RouterGroup) {
	// 安装 与api相关 的路由（包括获取系统api、获取api列表等）
	api.NewApiRouter(apiGroup)
	// 安装 操作历史记录相关 的路由
	operation.NewOperationRouter(apiGroup)
	// 安装 用户相关 的路由
	user.NewUserRouter(apiGroup)
}

// installOperationRouters 本方法里的api路由，都需要记录操作历史
func installOperationRouters(apiGroup *gin.RouterGroup) {
	// 需要操作记录，所以先加上一个 操作记录的中间件
	apiGroup.Use(middleware.OperationRecord())
	{
		// 安装swagger路由
		other.NewSwaggarRoute(apiGroup)
		// 安装k8s资源操作相关的路由
		kubeController.NewKubeRouter(apiGroup)
		// 安装菜单相关的路由
		menu.NewMenuRouter(apiGroup)
		// 安装csbin权限安全相关的路由
		authority.NewCasbinRouter(apiGroup)
	}
}
