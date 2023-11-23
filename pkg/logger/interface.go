package logger

import (
	"go.uber.org/zap"

	"github.com/noovertime7/kubemanage/pkg/globalError"
)

// Logger 自定义的日志接口，接口的实现里，实际是使用zap库的 *zap.Logger 对象，进行的日志输出
type Logger interface {
	// Info info级别日志
	Info(msg interface{})
	// Infof info级别字符串日志
	Infof(template string, args ...interface{})
	// Warn warning级别日志
	Warn(msg interface{})
	// Warnf warning级别字符串日志
	Warnf(template string, args ...interface{})
	// Error 错误级别日志
	Error(msg interface{})
	// ErrorWithCode 指定错误代码+error对象，输出错误日志
	ErrorWithCode(code int, err error)
	// ErrorWithErr 指定错误信息+error对象，输出错误日志
	ErrorWithErr(msg string, err error)
}

// New 新建一个Logger对象
func New() Logger {
	return logger{}
}

// logger Logger接口的一个实现
type logger struct{}

func (logger) Info(msg interface{}) {
	// TODO Sugar什么意思，需要仔细学习一下 zap 库
	LG.Sugar().Info(msg)
}

func (logger) Infof(template string, args ...interface{}) {
	LG.Sugar().Infof(template, args)
}
func (logger) Warn(msg interface{}) {
	LG.Sugar().Warn(msg)
}

func (logger) Warnf(template string, args ...interface{}) {
	LG.Sugar().Warnf(template, args)
}

func (logger) ErrorWithCode(code int, err error) {
	// 根据code，获取错误信息
	msg := globalError.GetErrorMsg(code)
	LG.Error(msg, zap.Error(err))
}

func (logger) Error(msg interface{}) {
	LG.Sugar().Error(msg)
}

func (logger) ErrorWithErr(msg string, err error) {
	LG.Error(msg, zap.Error(err))
}
