package middlewares

import (
	"groundTurn/src/constants"

	"github.com/gin-gonic/gin"
)

// RequireRole enforces RBAC: only listed roles may pass. Reaches constants,
// routes, frontend route guard and button visibility through /me.
func RequireRole(allowed ...string) gin.HandlerFunc {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, role := range allowed {
		allowedSet[role] = struct{}{}
	}
	return func(c *gin.Context) {
		role, _ := c.Get(CtxRole)
		roleStr, _ := role.(string)
		if _, ok := allowedSet[roleStr]; !ok {
			c.AbortWithStatusJSON(403, gin.H{
				"code":    constants.CodeRBACDenied,
				"message": constants.MsgRBACDenied,
				"details": gin.H{"required_roles": allowed, "actual_role": roleStr},
			})
			return
		}
		c.Next()
	}
}
