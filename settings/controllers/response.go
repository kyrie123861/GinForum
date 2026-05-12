package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code": 10000   //程序中的错误码
	"msg": xx       //提示信息
	"data": {}      //数据
}
*/

type ResponseData struct {
	Code MyCode      `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// ResponseError 返回标准错误响应
func ResponseError(c *gin.Context, code MyCode) {
	rd := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

// ResponseErrorWithMsg 返回带自定义消息的错误
func ResponseErrorWithMsg(c *gin.Context, code MyCode, errMsg interface{}) {
	rd := &ResponseData{
		Code: code,
		Msg:  errMsg,
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

// ResponseSuccess 返回标准正确的响应
func ResponseSuccess(c *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}
