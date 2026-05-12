package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"
)

// 投票功能

// 投票分数
// 投一票就加432分 86400/200 --》200赞成票就给帖子续一天
/* 投票的几种情况：
direction=1 时，有两种情况：
	1. 之前没有投过票，现在投赞成票
	2. 之前投反对票，现在改投赞成票
direction=0 时，有两种情况：
	1. 之前投过赞成票，现在要取消投票
	2. 之前投过反对票，现在要取消投票
direction=-1 时，有两种情况：
	1. 之前没有投过票，现在投反对票
	2. 之前投赞成票，现在改投反对票
投票的限制
	每个帖子自发布那一天起一周内可以投票，超过一个星期就不允许投票了
	1. 到期之后将redis中的赞成票和反对票保存到mysql中2
	2. 到期之后删除 KeyPostVotedPrefix 相关的记录
*/

// VoteForPost 投票功能
// func VoteForPost(userID int64, p *models.ParamVoteData) error {
// 	zap.L().Debug("VoteForPost",
// 		zap.Int64("userID", userID),
// 		zap.String("postID", p.PostID),
// 		zap.Int8("direction", p.Direction))
// 	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
// }

// deepseek修复redis 解决帖子投票问题出错问题
// logic/vote.go
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Info("VoteForPost开始",
		zap.Int64("userID", userID),
		zap.Int64("postID", p.PostID),
		zap.Int8("direction", p.Direction))

	// 转换为字符串
	userIDStr := strconv.FormatInt(userID, 10)
	postIDStr := strconv.FormatInt(p.PostID, 10)

	zap.L().Debug("转换后的ID",
		zap.String("userIDStr", userIDStr),
		zap.String("postIDStr", postIDStr))

	err := redis.VoteForPost(userIDStr, postIDStr, float64(p.Direction))
	if err != nil {
		zap.L().Error("redis.VoteForPost失败", zap.Error(err))
		return err
	}

	zap.L().Info("VoteForPost成功")
	return nil
}
