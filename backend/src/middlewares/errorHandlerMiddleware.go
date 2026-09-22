package middlewares

import (
	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
)

// ErrorHandlerMiddleware converts panics into the standard error envelope so
// services and controllers can each wrap their own AppError while unhandled
// failures still return a stable shape.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				abortEnvelope(c, 500, constants.CodeInternal, constants.MsgInternal)
			}
		}()
		c.Next()
	}
}
