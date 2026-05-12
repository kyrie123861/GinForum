package redis

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600 // 一周的秒数
	scorePerVote     = 432           // 每票的分数
)

func GetPostIDsInorder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取ID
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 查询索引地点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// ZRevRange 按分数从大到小的顺序查询指定数量的元素
	return rdb.ZRevRange(context.Background(), key, start, end).Result()
}

// GetPostVoteData 根据ids查询每篇帖子的投票赞成数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	ctx := context.Background() // 添加上下文
	// for _,id range ids{
	// 	key := getRedisKey(KeyPostVotedPrefix + id)
	// 	// 查找key中分数是1的元素数量--> 统计每篇帖子赞成票的数量
	// 	v := rdb.ZCount(key, 1, "1").Val()
	// 	data := append(data,v)
	// }
	// 使用pipeline一次发送多条命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedPrefix + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// deepseek修复redis 解决帖子投票问题出错问题
// CreatePost 创建帖子到 Redis（带检查，避免重复创建）
func CreatePost(postID int64) error {
	return CreatePostWithTime(postID, time.Now())
}

// CreatePostWithTime 使用指定时间创建帖子到 Redis
func CreatePostWithTime(postID int64, createTime time.Time) error {
	// 转换为字符串格式的 postID
	postIDStr := strconv.FormatInt(postID, 10)

	// 1. 首先检查帖子是否已在 Redis 中存在
	// 检查帖子时间是否已存在
	exists, err := rdb.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postIDStr).Result()
	if err == nil && exists > 0 {
		// 帖子已存在，不重复创建
		zap.L().Debug("帖子已存在于Redis中，跳过创建",
			zap.Int64("postID", postID),
			zap.Float64("existingTime", exists))
		return nil
	}

	// 2. 如果不存在，则创建
	pipeline := rdb.TxPipeline()

	createTimestamp := float64(createTime.Unix())

	// 帖子时间
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  createTimestamp,
		Member: postIDStr,
	})

	// 帖子分数 - 初始化为创建时间
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  createTimestamp,
		Member: postIDStr,
	})

	_, err = pipeline.Exec(ctx)
	if err != nil {
		zap.L().Error("创建帖子到Redis失败",
			zap.Int64("postID", postID),
			zap.Error(err))
		return err
	}

	zap.L().Debug("成功创建帖子到Redis",
		zap.Int64("postID", postID),
		zap.Float64("createTime", createTimestamp))

	return nil
}

// InitExistingPosts 初始化现有帖子到 Redis（用于修复数据）
func InitExistingPosts() error {
	// 这个函数应该在应用启动时调用一次
	// 或者通过API手动触发

	// 从数据库获取所有帖子
	// 注意：这里需要实现 mysql.GetAllPosts() 函数
	posts, err := mysql.GetAllPosts()
	if err != nil {
		return err
	}

	count := 0
	for _, post := range posts {
		// 检查是否已存在
		postIDStr := strconv.FormatInt(post.ID, 10)
		exists, _ := rdb.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postIDStr).Result()

		if exists == 0 {
			// 不存在，创建
			err := CreatePostWithTime(post.ID, post.CreateTime)
			if err != nil {
				zap.L().Error("初始化帖子到Redis失败",
					zap.Int64("postID", post.ID),
					zap.Error(err))
			} else {
				count++
				zap.L().Info("初始化帖子到Redis",
					zap.Int64("postID", post.ID))
			}
		}
	}

	zap.L().Info("初始化完成",
		zap.Int("total", len(posts)),
		zap.Int("created", count))

	return nil
}
