package controllers

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// -----社区相关-----

func CommunityHandler(c *gin.Context) {
	// 查询所有社区（community_id,community_name）以列表的形式返回
	data, err := logic.GetCommunity()
	if err != nil {
		zap.L().Error("logic.GetCommunity()	failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不暴露错误给服务器外
		return
	}
	ResponseSuccess(c, data)
}

func CommunityDetailHandler(c *gin.Context) {
	// 1.获取社区id
	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2.根据id获取社区详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunity()	failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不暴露错误给服务器外
		return
	}
	ResponseSuccess(c, data)
}
