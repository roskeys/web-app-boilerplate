package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
	"github.com/roskeys/app/db"
	"github.com/roskeys/app/utils"
)

var rateLimitCtx = context.Background()

func IPRateLimiter(rate redis_rate.Limit) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		RateLimit(c, key, rate)
	}
}

func RateLimit(c *gin.Context, key string, rate redis_rate.Limit) {
	res, err := db.RedisLimiter.Allow(rateLimitCtx, key, rate)
	if err != nil || res.Allowed == 0 {
		utils.SendErrorResponse(c, utils.TO_MANY_REQUESTS)
	}
	c.Next()
}
