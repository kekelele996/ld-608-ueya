package controllers

import (
	"groundTurn/src/middlewares"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// fail renders a service-layer *utils.APIError or a generic 500. Controllers
// wrap errors here; services wrap business errors before returning them.
func fail(c *gin.Context, err error) {
	if apiErr, ok := err.(*utils.APIError); ok {
		c.JSON(apiErr.HTTPCode, gin.H{
			"code": apiErr.Code, "message": apiErr.Message, "details": apiErr.Details,
		})
		return
	}
	c.JSON(500, gin.H{"code": "INTERNAL_ERROR", "message": err.Error()})
}

// actor reads the JWT identity injected by authMiddleware.
func actor(c *gin.Context) (username, role, teamID string) {
	return middlewares.CurrentUser(c)
}

func utilsPagination(c *gin.Context) (int, int) {
	return utils.Pagination(c)
}
