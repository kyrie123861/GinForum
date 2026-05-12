package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	// 1.获取参数以及参数的校验
	// c.ShouldBindJSON // validator --> binding tag
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create psot with invalid param")
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 从 c 取到当前用户的ID并赋值给 post 的 AuthorID 字段
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	p.AuthorID = userID

	// 2.创建帖子
	if err := logic.CreaterPost(p); err != nil {
		zap.L().Error("logic.CreaterPost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	// 1.获取参数（从URL中获取帖子的id）
	pidstr := c.Param("id")
	pid, err := strconv.ParseInt(pidstr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2.根据id去数据库查询帖子详情
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3.返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)

	// 1.获取参数
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2.返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版获取帖子列表接口
// 根据前端传来的参数动态获取帖子的列表
// 按照创建时间顺序或者按照分数排序
// 1.获取参数
// 2.去redis查询id列表
// 3.根据id去数据库查询帖子详细信息

func GetPostListHandler2(c *gin.Context) {
	// GET请求参数 /api/v1/posts2/?page=2&size=10&order=time
	// 初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("c.ShouldBindQuery(p) error", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	// c.ShouldBing() 根据请求的数据类型选择相应的方法获取数据
	// c.ShouldBindJSON() 如果请求中携带的是JSON格式的数据，才使用他
	// 获取分页参数

	// 1.获取参数
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2.返回响应
	ResponseSuccess(c, data)
}

// deepseek修复redis 解决帖子投票问题出错问题
// InitRedisPostsHandler 初始化帖子到 Redis（仅供测试/初始化使用）
func InitRedisPostsHandler(c *gin.Context) {
	// 注意：这个接口应该有权限控制，只在开发或测试环境使用
	go func() {
		err := logic.InitPostsToRedis()
		if err != nil {
			zap.L().Error("InitPostsToRedis failed", zap.Error(err))
		}
	}()

	ResponseSuccess(c, gin.H{"message": "开始初始化帖子到Redis"})
}
