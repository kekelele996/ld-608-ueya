package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/middlewares"
	"groundTurn/src/services"
)

type ReportController struct{ Svc *services.Services }

func NewReportController(svc *services.Services) *ReportController {
	return &ReportController{Svc: svc}
}

func (ctl *ReportController) Dashboard(c *gin.Context) {
	stats, err := ctl.Svc.Reports.Dashboard()
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, stats)
}

func (ctl *ReportController) AuditLogs(c *gin.Context) {
	rows, err := ctl.Svc.Reports.AuditLogs(100)
	if err != nil {
		fail(c, err)
		return
	}
	middlewares.OK(c, http.StatusOK, rows)
}
