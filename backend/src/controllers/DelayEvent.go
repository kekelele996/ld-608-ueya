package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/middlewares"
	"groundTurn/src/services"
	"groundTurn/src/types"
)

type DelayController struct{ Svc *services.Services }

func NewDelayController(svc *services.Services) *DelayController {
	return &DelayController{Svc: svc}
}

func (ctl *DelayController) List(c *gin.Context) {
	rows, err := ctl.Svc.Delays.List()
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, rows)
}

// Register adds delay minutes and returns the same replanned result object
// the detail page refreshes to: moved tasks/kept tasks/bumped bookings.
func (ctl *DelayController) Register(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	var req types.RegisterDelayRequest
	if !MustBindJSON(c, &req) {
		return
	}
	result, err := ctl.Svc.Delays.Register(middlewares.ActorFrom(c), id, req)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, result)
}

func (ctl *DelayController) Resolve(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	row, err := ctl.Svc.Delays.Resolve(middlewares.ActorFrom(c), id)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, row)
}

func (ctl *DelayController) Impacts(c *gin.Context) {
	rows, err := ctl.Svc.Reports.DelayImpacts()
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, rows)
}
