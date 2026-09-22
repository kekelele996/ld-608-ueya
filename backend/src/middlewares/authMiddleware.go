package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"groundTurn/src/config"
	"groundTurn/src/constants"
	"groundTurn/src/services"
	"groundTurn/src/utils"
)

const actorContextKey = "actor"

// AuthMiddleware validates the Bearer JWT and injects the actor.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" || token == header {
			abortEnvelope(c, http.StatusUnauthorized, constants.CodeAuthRequired, constants.MsgAuthRequired)
			return
		}
		claims, err := utils.ParseToken(cfg.JWTSecret, token)
		if err != nil {
			abortEnvelope(c, http.StatusUnauthorized, constants.CodeAuthInvalid, constants.MsgAuthInvalid)
			return
		}
		c.Set(actorContextKey, services.Actor{
			Username: claims.Username,
			Role:     claims.Role,
			TeamID:   claims.TeamID,
			Name:     claims.Name,
		})
		c.Next()
	}
}

// ActorFrom extracts the authenticated actor inside controllers.
func ActorFrom(c *gin.Context) services.Actor {
	value, ok := c.Get(actorContextKey)
	if !ok {
		return services.Actor{}
	}
	actor, _ := value.(services.Actor)
	return actor
}
