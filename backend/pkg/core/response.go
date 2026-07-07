package core

import (
	"net/http"

	"cloudflared-tunnel/pkg/errno"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// OK 成功响应
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "成功",
		Data:    data,
	})
}

// OKList 列表响应
func OKList(c *gin.Context, list any, total int64) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "成功",
		Data: gin.H{
			"list":  list,
			"total": total,
		},
	})
}

// Fail 错误响应
func Fail(c *gin.Context, err error) {
	code, message := errno.Decode(err)
	httpStatus := errno.HTTPStatus(err)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

// FailMsg 自定义错误消息
func FailMsg(c *gin.Context, err error, msg string) {
	code, _ := errno.Decode(err)
	httpStatus := errno.HTTPStatus(err)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: msg,
	})
}
