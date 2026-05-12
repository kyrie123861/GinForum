package controllers

import (
	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 投票

// type ParamVoteData struct {
// 	// UserID 从请求中获取当前用户
// 	PostID    int64 `json:"post_id,string"`   //帖子ID
// 	Direction int   `json:"direction,string"` // 赞成票(1) 反对票(-1)
// }

func PostVoteCotroller(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		errData := RemoveTopStruct(errs.Translate(trans)) //翻译并去除错误提示中的结构体标识
		ResponseErrorWithMsg(c, CodeInvalidParams, errData)
		return
	}

	// 获取用户ID
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost(userID, p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}
