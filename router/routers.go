package router

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	// Default是创建一个引擎，并默认注册了 Logger 和 Recovery 两个中间件
	//r := gin.Default()
	// 使用New() 创建的是一个空白的引擎，不带任何中间件
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	//注册业务路由
	v1.POST("/signup", controllers.SignUpHandler)

	//登录功能得实现
	v1.POST("/login", controllers.LoginHandler)

	//运用JWT定义中间件
	v1.Use(middlewares.JWTAuthMiddleware())

	{
		//注册社区
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)

		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.GetPostDetailHandler)

		// 两者都是获取帖子列表，但是实现逻辑不一样
		v1.GET("/posts/", controllers.GetPostListHandler)
		// 根据时间或者分数获得帖子列表
		v1.GET("/posts2", controllers.GetPostListHandler2)
		//投票
		v1.POST("/vote", controllers.PostVoteCotroller)

		// deepseek修复redis 解决帖子投票问题出错问题
		r.GET("/init/redis/posts", controllers.InitRedisPostsHandler)
	}

	// pprof 性能优化指标
	pprof.Register(r)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": 404,
		})
	})
	return r
}
