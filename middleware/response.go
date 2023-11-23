package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/noovertime7/kubemanage/pkg/globalError"
	"github.com/pkg/errors"
	"net/http"
)

// ResponseCode 创建一个int类型的变量
type ResponseCode int

// Response 自定义的响应体结构
type Response struct {
	// Code 响应状态码（这是响应体中的code，用于代表我们的处理结果，并非请求的状态码）
	Code ResponseCode `json:"code"`
	// Msg 响应信息
	Msg string `json:"msg"`
	// RealErr 当结果为错误时，存储错误信息
	RealErr string `json:"real_err"`
	// Data 响应数据
	Data interface{} `json:"data"`
}

// ResponseSuccess 响应成功的公共处理，内部会封装好响应体，设置给context
func ResponseSuccess(c *gin.Context, data interface{}) {
	resp := &Response{Code: http.StatusOK, Msg: "", Data: data}
	// data如果是string，则直接将data写入msg
	tempMsg, ok := data.(string)
	if ok && tempMsg == "" {
		resp.Msg = "操作成功"
	}
	// 请求的response中展示业务处理的code，请求本身的响应状态码也应该设置为200
	c.JSON(200, resp)
	// 将resp对象，序列化成json格式
	response, _ := json.Marshal(resp)
	// 将封装好的响应体response，设置给 上下文对象Context，后面就可以直接返回给前端了
	c.Set("response", string(response))
}

// ResponseError 响应失败的公共处理，内部会封装好响应体，设置给context
func ResponseError(c *gin.Context, err error) {
	//判断错误类型
	// As - 获取错误的具体实现
	var code ResponseCode
	var myError = new(globalError.GlobalError)
	if errors.As(err, &myError) {
		code = ResponseCode(myError.Code)
	}
	resp := &Response{Code: code, Msg: err.Error(), RealErr: myError.RealErrorMessage, Data: ""}
	// 即使发生错误，也只是在请求的response中展示，请求本身的响应状态码还应该设置为200
	c.JSON(200, resp)
	// 将resp对象，序列化成json格式
	response, _ := json.Marshal(resp)
	// 将封装好的响应体response，设置给 上下文对象Context，后面就可以直接返回给前端了
	c.Set("response", string(response))
}
