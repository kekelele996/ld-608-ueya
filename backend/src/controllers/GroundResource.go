package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/middlewares"
	"groundTurn/src/services"
	"groundTurn/src/types"
)

type ResourceController struct{ Svc *services.Services }

func NewResourceController(svc *services.Services) *ResourceController {
	return &ResourceController{Svc: svc}
}

func (ctl *ResourceController) List(c *gin.Context) {
	rows, err := ctl.Svc.Resources.List()
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, rows)
}

func (ctl *ResourceController) Bookings(c *gin.Context) {
	rows, err := ctl.Svc.Resources.AllBookings()
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, rows)
}

func (ctl *ResourceController) PendingBookings(c *gin.Context) {
	rows, err := ctl.Svc.Resources.PendingBookings()
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, rows)
}

func (ctl *ResourceController) SetMaintenance(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	var req types.MaintenanceRequest
	if !MustBindJSON(c, &req) {
		return
	}
	row, err := ctl.Svc.Resources.SetMaintenance(middlewares.ActorFrom(c), id, req)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, row)
}
