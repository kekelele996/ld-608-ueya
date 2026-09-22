package controllers

import (
	"groundTurn/src/constructors"
	"groundTurn/src/services"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// AuthController exposes local login (JWT) and identity inspection.
type AuthController struct {
	Auth *services.AuthService
}

func NewAuthController(auth *services.AuthService) *AuthController {
	return &AuthController{Auth: auth}
}

// Login accepts {"username":"dispatcher"} and returns a JWT.
func (a *AuthController) Login(c *gin.Context) {
	var req constructors.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": "VALIDATION_FAILED", "message": err.Error()})
		return
	}
	token, user, err := a.Auth.Login(req.Username)
	if err != nil {
		if apiErr, ok := err.(*utils.APIError); ok {
			c.JSON(apiErr.HTTPCode, gin.H{"code": apiErr.Code, "message": apiErr.Message})
			return
		}
		c.JSON(500, gin.H{"code": "INTERNAL_ERROR", "message": err.Error()})
		return
	}
	c.JSON(200, constructors.LoginResponse{
		Token: token, Username: user.Username, DisplayName: user.DisplayName,
		Role: user.Role, TeamID: user.TeamID,
	})
}

// Me returns identity for the frontend route guard / button visibility.
func (a *AuthController) Me(c *gin.Context) {
	c.JSON(200, gin.H{
		"username": c.GetString("ctxUsername"),
		"role":     c.GetString("ctxRole"),
		"team_id":  c.GetString("ctxTeamID"),
	})
}
