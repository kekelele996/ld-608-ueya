package middlewares

import (
	"sync"
	"time"

	"groundTurn/src/constants"

	"github.com/gin-gonic/gin"
)

type rateBucket struct {
	tokens float64
	seen   time.Time
}

// RateLimit is a tiny per-IP token bucket (no third-party limiter).
func RateLimit(qps int) gin.HandlerFunc {
	var mu sync.Mutex
	buckets := map[string]*rateBucket{}
	return func(c *gin.Context) {
		key := c.ClientIP()
		mu.Lock()
		bucket, ok := buckets[key]
		now := time.Now()
		if !ok {
			bucket = &rateBucket{tokens: float64(qps), seen: now}
			buckets[key] = bucket
		}
		elapsed := now.Sub(bucket.seen).Seconds()
		bucket.tokens = minF(float64(qps), bucket.tokens+elapsed*float64(qps))
		bucket.seen = now
		allowed := bucket.tokens >= 1
		if allowed {
			bucket.tokens--
		}
		mu.Unlock()
		if !allowed {
			c.AbortWithStatusJSON(429, gin.H{
				"code": constants.CodeRateLimited, "message": constants.MsgRateLimited,
			})
			return
		}
		c.Next()
	}
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
