package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/middlewares"
	"groundTurn/src/services"
	"groundTurn/src/types"
)

type BookingController struct{ Svc *services.Services }

func NewBookingController(svc *services.Services) *BookingController {
	return &BookingController{Svc: svc}
}

func (ctl *BookingController) List(c *gin.Context) {
	rows, err := ctl.Svc.Bookings.List()
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, rows)
}

func (ctl *BookingController) Adjust(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	var req types.AdjustBookingRequest
	if !MustBindJSON(c, &req) {
		return
	}
	row, err := ctl.Svc.Bookings.Adjust(middlewares.ActorFrom(c), id, req)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, row)
}

func (ctl *BookingController) Release(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	row, err := ctl.Svc.Bookings.Release(middlewares.ActorFrom(c), id)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, row)
}
