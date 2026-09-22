package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/middlewares"
	"groundTurn/src/services"
	"groundTurn/src/types"
)

type AuthController struct{ Svc *services.Services }

func NewAuthController(svc *services.Services) *AuthController {
	return &AuthController{Svc: svc}
}

func (ctl *AuthController) Login(c *gin.Context) {
	var req types.LoginRequest
	if !MustBindJSON(c, &req) {
		return
	}
	result, err := ctl.Svc.Auth.Login(req)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, result)
}

func (ctl *AuthController) Me(c *gin.Context) {
	middlewares.OK(c, http.StatusOK, middlewares.ActorFrom(c))
}
