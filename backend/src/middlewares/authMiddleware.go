package middlewares

import (
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// Identity keys stored in gin.Context and read by services/controllers.
const (
	CtxUsername = "ctxUsername"
	CtxRole     = "ctxRole"
	CtxTeamID   = "ctxTeamID"
)

// AuthRequired verifies the bearer JWT and injects identity.
func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"code": "AUTH_REQUIRED", "message": "缺少登录令牌，请先登录",
			})
			return
		}
		claims, err := utils.ParseToken(secret, token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code": "AUTH_REQUIRED", "message": "登录令牌无效或已过期",
			})
			return
		}
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxRole, claims.Role)
		c.Set(CtxTeamID, claims.TeamID)
		c.Next()
	}
}

// CurrentUser reads identity injected by AuthRequired.
func CurrentUser(c *gin.Context) (username, role, teamID string) {
	if v, ok := c.Get(CtxUsername); ok {
		username, _ = v.(string)
	}
	if v, ok := c.Get(CtxRole); ok {
		role, _ = v.(string)
	}
	if v, ok := c.Get(CtxTeamID); ok {
		teamID, _ = v.(string)
	}
	return
}
