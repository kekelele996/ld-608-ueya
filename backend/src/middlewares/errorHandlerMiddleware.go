package middlewares

import (
	"net/http"

	"groundTurn/src/constants"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// ErrorHandler renders *utils.APIError. Controllers and services each wrap
// their own errors rather than collapsing everything here.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		if apiErr, ok := err.(*utils.APIError); ok {
			c.JSON(apiErr.HTTPCode, gin.H{
				"code":    apiErr.Code,
				"message": apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    constants.CodeInternal,
			"message": constants.MsgInternal,
		})
	}
}
