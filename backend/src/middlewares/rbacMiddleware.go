package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
)

// RBACMiddleware enforces the action -> roles matrix from constants/rbac.go.
func RBACMiddleware(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		actor := ActorFrom(c)
		if !constants.Can(constants.UserRole(actor.Role), action) {
			abortEnvelope(c, http.StatusForbidden, constants.CodeRBACDenied, constants.MsgRBACDenied)
			return
		}
		c.Next()
	}
}
