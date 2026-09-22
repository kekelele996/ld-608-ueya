package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"groundTurn/src/constructors"
)

// AuditLogMiddleware records mutating API calls as operation logs. Business
// state changes are additionally logged inside services with field details.
func AuditLogMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		mutating := c.Request.Method == "POST" || c.Request.Method == "PUT" ||
			c.Request.Method == "PATCH" || c.Request.Method == "DELETE"
		if !mutating || c.Request.URL.Path == "/api/auth/login" {
			return
		}
		actor := ActorFrom(c)
		row := constructors.NewAuditLog(
			actor.Username, actor.Role,
			fmt.Sprintf("HTTP_%s", c.Request.Method),
			c.Request.URL.Path, "",
			fmt.Sprintf("%s %s -> %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status()),
		)
		row.CreatedAt = time.Now()
		_ = db.Create(row).Error
	}
}
