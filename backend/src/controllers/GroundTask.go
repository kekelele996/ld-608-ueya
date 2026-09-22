package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/middlewares"
	"groundTurn/src/services"
	"groundTurn/src/types"
)

type TaskController struct{ Svc *services.Services }

func NewTaskController(svc *services.Services) *TaskController {
	return &TaskController{Svc: svc}
}

func (ctl *TaskController) List(c *gin.Context) {
	rows, err := ctl.Svc.Tasks.ListAll()
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, rows)
}

func (ctl *TaskController) Sign(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	var req types.SignTaskRequest
	_ = c.ShouldBindJSON(&req)
	row, err := ctl.Svc.Tasks.Sign(middlewares.ActorFrom(c), id, req)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, row)
}

func (ctl *TaskController) Finish(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	var req types.FinishTaskRequest
	_ = c.ShouldBindJSON(&req)
	row, err := ctl.Svc.Tasks.Finish(middlewares.ActorFrom(c), id, req)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, row)
}

func (ctl *TaskController) Block(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	var req types.BlockTaskRequest
	if !MustBindJSON(c, &req) {
		return
	}
	row, err := ctl.Svc.Tasks.Block(middlewares.ActorFrom(c), id, req)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, row)
}
