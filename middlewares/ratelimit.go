package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
// 	bucket := ratelimit.NewBucket(fillInterval, cap)
// 	return func(c *gin.Context) {
// 		// 如果取不到令牌就中断本次请求返回 rate limit...
// 		if bucket.TakeAvailable(1) < 1 {
// 			c.String(http.StatusOK, "rate limit...")
// 			c.Abort()
// 			return
// 		}
// 		c.Next()
// 	}
// }

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	// 创建令牌桶限流器
	limiter := rate.NewLimiter(rate.Every(fillInterval), int(cap))
	return func(c *gin.Context) {
		// 尝试获取1个令牌，失败则返回限流
		if !limiter.Allow() {
			c.String(http.StatusTooManyRequests, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}
