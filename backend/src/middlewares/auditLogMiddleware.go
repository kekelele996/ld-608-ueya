package middlewares

import "github.com/gin-gonic/gin"

// AuditEntry is recorded by services for every write operation; this
// middleware only exposes a request-scoped recorder hook. Persistence is
// intentionally done inside services so business context is available.
func AuditContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("auditEnabled", true)
		c.Next()
	}
}
