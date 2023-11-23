package api

import "github.com/gin-gonic/gin"

type apiController struct{}

// NewApiRouter 初始化 Api相关的路由
func NewApiRouter(ginGroup *gin.RouterGroup) {
	api := &apiController{}
	// 初始化 Api相关的路由
	api.initRoutes(ginGroup)
}

// initRoutes 初始化 Api相关的路由
func (a *apiController) initRoutes(ginGroup *gin.RouterGroup) {
	// 创建一个 api组，前缀为：/sysApi
	apiRoute := ginGroup.Group("/sysApi")
	// 注册一个接口 /getAPiList
	apiRoute.GET("/getAPiList", a.GetApiList)
}
