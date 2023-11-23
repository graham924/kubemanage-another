package v1

import (
	"github.com/noovertime7/kubemanage/cmd/app/config"
	"github.com/noovertime7/kubemanage/cmd/app/options"
	"github.com/noovertime7/kubemanage/pkg/logger"
)

// CoreV1 全局的k8s核心资源服务对象
var CoreV1 CoreService

// Log 全局的日志对象
var Log Logger

// Setup 完成核心应用接口的设置
func Setup(o *options.Options) {
	// 新建一个Logger对象，作为全局日志处理器
	Log = logger.New()
	// 新建一个CoreService对象，作为全局和k8s交互的处理器
	CoreV1 = New(config.SysConfig, o.Factory)
}
