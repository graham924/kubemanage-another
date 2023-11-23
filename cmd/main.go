package main

import (
	"github.com/gin-gonic/gin"
	"github.com/noovertime7/kubemanage/cmd/app"
	"os"
)

func main() {
	// 设置gin框架的运行模式：release
	gin.SetMode(gin.ReleaseMode)
	// 创建程序启动命令
	cmd := app.NewServerCommand()
	// 执行命令，启动服务器
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
