package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreaterPost(p *models.Post) (err error) {
	// 1.生成postid
	ID, _ := snowflake.GenID()
	p.ID = int64(ID)
	// 2.保存数据到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID)
	if err != nil {
		return err
	}
	return
	// 3.返回
}

// GetPostById 根据id获取帖子详情
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们接口想用的数据
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Error(err))
		return
	}

	// 根据作者id查询作者信息
	user, err := mysql.GetUserBuId(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserBuId(post.AuthorID) failed", zap.Error(err))
		return
	}

	// 根据社区id查询社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailById(post.CommunityID) failed", zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		// 根据作者id查询作者信息
		user, errs := mysql.GetUserBuId(post.AuthorID)
		if errs != nil {
			zap.L().Error("mysql.GetUserBuId(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(errs))
			continue
		}

		// 根据社区id查询社区信息
		community, errs := mysql.GetCommunityDetailByID(post.CommunityID)
		if errs != nil {
			zap.L().Error("mysql.GetCommunityDetailById(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(errs))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postDetail)
	}

	return
}

// GetPostList2 升级版根据id获取帖子数据
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2.去redis查询id列表
	ids, err := redis.GetPostIDsInorder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInorder(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 3.根据id去mysql数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id顺序返回
	posts, err := mysql.GetPostByIDs(ids)
	if err != nil {
		return
	}

	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将贴子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, errs := mysql.GetUserBuId(post.AuthorID)
		if errs != nil {
			zap.L().Error("mysql.GetUserBuId(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(errs))
			continue
		}

		// 根据社区id查询社区信息
		community, errs := mysql.GetCommunityDetailByID(post.CommunityID)
		if errs != nil {
			zap.L().Error("mysql.GetCommunityDetailById(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(errs))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VotesNum:        voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postDetail)
	}
	return
}

// deepseek修复redis 解决帖子投票问题出错问题
// logic/post.go
func InitPostsToRedis() error {
	// 获取所有帖子（包括创建时间）
	posts, err := mysql.GetAllPostsWithTime()
	if err != nil {
		zap.L().Error("mysql.GetAllPostsWithTime() failed", zap.Error(err))
		return err
	}

	zap.L().Info("Initializing posts to Redis", zap.Int("count", len(posts)))

	for _, post := range posts {
		// 添加到 Redis，使用实际的创建时间
		err := redis.CreatePostWithTime(post.ID, post.CreateTime)
		if err != nil {
			zap.L().Error("redis.CreatePostWithTime failed",
				zap.Int64("post_id", post.ID),
				zap.Error(err))
		} else {
			zap.L().Info("Post added to Redis",
				zap.Int64("post_id", post.ID))
		}
	}

	return nil
}
