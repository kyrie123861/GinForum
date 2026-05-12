package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 投票功能

// 投票分数
// 投一票就加432分 86400/200 --》200赞成票就给帖子续一天
/* 投票的几种情况：
direction=1 时，有两种情况：
	1. 之前没有投过票，现在投赞成票 --> 更新分数和投票记录 差值的绝对值：1 +432
	2. 之前投反对票，现在改投赞成票 --> 更新分数和投票记录 差值的绝对值：2 +432*2
direction=0 时，有两种情况：
	1. 之前投过反对票，现在要取消投票 --> 更新分数和投票记录 差值的绝对值：1 +432
	2. 之前投过赞成票，现在要取消投票 --> 更新分数和投票记录 差值的绝对值：1 -432
direction=-1 时，有两种情况：
	1. 之前没有投过票，现在投反对票 --> 更新分数和投票记录 差值的绝对值：1 -432
	2. 之前投赞成票，现在改投反对票 --> 更新分数和投票记录 差值的绝对值：2 -432*2
投票的限制：
	1. 每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	2. 到期之后将 redis 中保存的赞成票数及反对票数存储到 mysql 表中
	3. 到期之后删除那个 KeyPostVotedPrefix 相关的记录
*/

// const (
// 	oneWeekInSeconds = 7 * 24 * 3600 // 一周的秒数
// 	scorePerVote     = 432           // 每票的分数
// )

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
)
var ctx = context.Background()

// func CreatePost(postID int64) error {
// 	pipeline := rdb.TxPipeline()
// 	// 帖子时间
// 	rdb.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), redis.Z{
// 		Score:  float64(time.Now().Unix()),
// 		Member: postID,
// 	})

// 	// 帖子分数
// 	rdb.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), redis.Z{
// 		Score:  float64(time.Now().Unix()),
// 		Member: postID,
// 	})
// 	_, err := pipeline.Exec(ctx)
// 	return err
// }

// func VoteForPost(userID, postID string, value float64) error {
// 	// 1.判断投票限制
// 	// 取redis发布帖子的时间
// 	postTime := rdb.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
// 	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
// 		return ErrorVoteTimeExpire
// 	}
// 	// 2.更新帖子的分数
// 	// 查询之前投票纪录
// 	ov := rdb.ZScore(ctx, getRedisKey(KeyPostVotedPrefix+postID), userID).Val()
// 	var op float64
// 	if value > ov {
// 		op = 1
// 	} else {
// 		op = -1
// 	}
// 	diff := math.Abs(ov - value) // 计算本次投票和之前投票的差值
// 	err := rdb.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID).Err()
// 	if err != nil {
// 		return err
// 	}
// 	// 3.记录用户为该帖子投票的数据
// 	if value == 0 {
// 		rdb.ZRem(ctx, getRedisKey(KeyPostVotedPrefix+postID), userID).Result()
// 	} else {
// 		_, err = rdb.ZAdd(ctx, getRedisKey(KeyPostVotedPrefix+postID), redis.Z{
// 			Score:  value,
// 			Member: userID,
// 		}).Result()
// 	}
// 	return err
// }

// deepseek修复redis 解决帖子投票问题出错问题
// dao/redis/post.go
func VoteForPost(userID, postID string, value float64) error {
	zap.L().Debug("VoteForPost开始",
		zap.String("userID", userID),
		zap.String("postID", postID),
		zap.Float64("value", value))

	// 1. 判断投票限制
	postTime := rdb.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	zap.L().Debug("帖子时间",
		zap.Float64("postTime", postTime),
		zap.Float64("currentTime", float64(time.Now().Unix())),
		zap.Float64("timeDiff", float64(time.Now().Unix())-postTime),
		zap.Float64("oneWeek", oneWeekInSeconds))

	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		zap.L().Warn("投票时间已过期")
		return ErrorVoteTimeExpire
	}

	// 2. 获取之前的投票记录
	votedKey := getRedisKey(KeyPostVotedPrefix + postID)
	zap.L().Debug("投票记录键", zap.String("votedKey", votedKey))

	ov := rdb.ZScore(ctx, votedKey, userID).Val()
	zap.L().Debug("之前的投票值",
		zap.Float64("ov", ov),
		zap.String("userID", userID))

	// 3. 计算分数变化
	var scoreDiff float64
	var op string

	// 情况分析：
	switch {
	case value == 0 && ov == 1: // 取消赞成票
		scoreDiff = -scorePerVote
		op = "取消赞成票"
	case value == 0 && ov == -1: // 取消反对票
		scoreDiff = scorePerVote
		op = "取消反对票"
	case ov == 0 && value == 1: // 新投赞成票
		scoreDiff = scorePerVote
		op = "新投赞成票"
	case ov == 0 && value == -1: // 新投反对票
		scoreDiff = -scorePerVote
		op = "新投反对票"
	case ov == -1 && value == 1: // 反对改赞成
		scoreDiff = 2 * scorePerVote
		op = "反对改赞成"
	case ov == 1 && value == -1: // 赞成改反对
		scoreDiff = -2 * scorePerVote
		op = "赞成改反对"
	default:
		// 其他情况（如重复投票）不改变分数
		scoreDiff = 0
		op = "无变化"
	}

	zap.L().Debug("投票操作",
		zap.String("op", op),
		zap.Float64("scoreDiff", scoreDiff),
		zap.Float64("scorePerVote", scorePerVote))

	// 4. 更新帖子分数
	if scoreDiff != 0 {
		scoreKey := getRedisKey(KeyPostScoreZSet)
		zap.L().Debug("更新分数",
			zap.String("scoreKey", scoreKey),
			zap.Float64("scoreDiff", scoreDiff))

		// 先获取当前分数
		currentScore := rdb.ZScore(ctx, scoreKey, postID).Val()
		zap.L().Debug("当前分数", zap.Float64("currentScore", currentScore))

		err := rdb.ZIncrBy(ctx, scoreKey, scoreDiff, postID).Err()
		if err != nil {
			zap.L().Error("更新分数失败", zap.Error(err))
			return err
		}

		// 获取更新后的分数
		newScore := rdb.ZScore(ctx, scoreKey, postID).Val()
		zap.L().Debug("更新后分数", zap.Float64("newScore", newScore))
	} else {
		zap.L().Debug("分数无变化，跳过更新")
	}

	// 5. 更新投票记录
	if value == 0 {
		zap.L().Debug("删除投票记录")
		result, err := rdb.ZRem(ctx, votedKey, userID).Result()
		if err != nil {
			zap.L().Error("删除投票记录失败", zap.Error(err))
			return err
		}
		zap.L().Debug("删除结果", zap.Int64("result", result))
	} else {
		zap.L().Debug("添加/更新投票记录",
			zap.Float64("value", value))

		result, err := rdb.ZAdd(ctx, votedKey, redis.Z{
			Score:  value,
			Member: userID,
		}).Result()

		if err != nil {
			zap.L().Error("添加投票记录失败", zap.Error(err))
			return err
		}

		zap.L().Debug("添加结果",
			zap.Int64("result", result),
			zap.String("votedKey", votedKey))

		// 验证添加结果
		members, err := rdb.ZRangeWithScores(ctx, votedKey, 0, -1).Result()
		if err != nil {
			zap.L().Error("验证投票记录失败", zap.Error(err))
		} else {
			zap.L().Debug("投票记录内容",
				zap.Int("count", len(members)))
			for _, member := range members {
				zap.L().Debug("投票成员",
					zap.String("member", member.Member.(string)),
					zap.Float64("score", member.Score))
			}
		}
	}

	zap.L().Debug("VoteForPost结束")
	return nil
}
