package middlewares

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
)

// RateLimitMiddleware is a tiny per-IP token bucket guard (60 req/min).
func RateLimitMiddleware() gin.HandlerFunc {
	const limit = 60
	const window = time.Minute
	var mu sync.Mutex
	hits := map[string]int{}
	reset := map[string]time.Time{}

	return func(c *gin.Context) {
		key := c.ClientIP()
		mu.Lock()
		now := time.Now()
		if deadline, ok := reset[key]; !ok || now.After(deadline) {
			hits[key] = 0
			reset[key] = now.Add(window)
		}
		hits[key]++
		exceeded := hits[key] > limit
		mu.Unlock()
		if exceeded {
			abortEnvelope(c, 429, constants.CodeRateLimited, constants.MsgRateLimited)
			return
		}
		c.Next()
	}
}
