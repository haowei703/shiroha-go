package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// Success 生成成功响应
func Success(c *gin.Context, data interface{}) {
	result := Result{
		Code: http.StatusOK,
		Data: data,
		Msg:  "Success",
	}
	c.JSON(http.StatusOK, result)
}

// Error 生成错误响应
func Error(c *gin.Context, code int, msg string) {
	result := Result{
		Code: code,
		Data: nil,
		Msg:  msg,
	}
	c.AbortWithStatusJSON(code, result)
}
