package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/middlewares"
	"groundTurn/src/services"
	"groundTurn/src/types"
)

type TurnaroundController struct{ Svc *services.Services }

func NewTurnaroundController(svc *services.Services) *TurnaroundController {
	return &TurnaroundController{Svc: svc}
}

func (ctl *TurnaroundController) List(c *gin.Context) {
	rows, err := ctl.Svc.Turnarounds.List()
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, rows)
}

func (ctl *TurnaroundController) Detail(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	detail, err := ctl.Svc.Turnarounds.Detail(id)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, detail)
}

func (ctl *TurnaroundController) Create(c *gin.Context) {
	var req types.CreateTurnaroundRequest
	if !MustBindJSON(c, &req) {
		return
	}
	row, err := ctl.Svc.Turnarounds.Create(middlewares.ActorFrom(c), req)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusCreated, row)
}

func (ctl *TurnaroundController) Arrive(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	var req types.ArriveRequest
	_ = c.ShouldBindJSON(&req)
	row, err := ctl.Svc.Turnarounds.Arrive(middlewares.ActorFrom(c), id, req)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, row)
}

func (ctl *TurnaroundController) Release(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	row, err := ctl.Svc.Turnarounds.Release(middlewares.ActorFrom(c), id)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, row)
}

// GeneratePlan returns 409 with itemised conflicts when the batch is rejected;
// on success the whole batch (tasks + bookings) is committed atomically.
func (ctl *TurnaroundController) GeneratePlan(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	result, err := ctl.Svc.Turnarounds.GeneratePlan(middlewares.ActorFrom(c), id)
	if err != nil {
		fail(c, err)
		return
	}
	if result.Saved {
		middlewares.OK(c, http.StatusCreated, result)
		return
	}
	middlewares.Conflict(c, http.StatusConflict, "PLAN_CONFLICT",
		"保障计划整批未保存，请逐项处理下列冲突", result)
}
