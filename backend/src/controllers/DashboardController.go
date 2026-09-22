package controllers

import (
	"groundTurn/src/constructors"
	"groundTurn/src/services"

	"github.com/gin-gonic/gin"
)

// DashboardController aggregates stats, pending bookings and audit logs.
type DashboardController struct {
	Query *services.QueryService
}

func NewDashboardController(query *services.QueryService) *DashboardController {
	return &DashboardController{Query: query}
}

func (d *DashboardController) Stats(c *gin.Context) {
	stats, err := d.Query.DashboardStats()
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, stats)
}

func (d *DashboardController) PendingBookings(c *gin.Context) {
	rows, err := d.Query.PendingBookings()
	if err != nil {
		fail(c, err)
		return
	}
	resources, _ := d.Query.AllResources()
	rm := map[uint]string{}
	for _, r := range resources {
		rm[r.ID] = r.ResourceCode
	}
	out := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		view := constructors.NewBookingView(row, nil)
		view.ResourceCode = rm[row.ResourceID]
		out = append(out, gin.H{
			"id": view.ID, "resource_id": view.ResourceID, "resource_code": view.ResourceCode,
			"turnaround_id": view.TurnaroundID, "task_id": view.TaskID,
			"start_time": view.StartTime, "end_time": view.EndTime,
			"status_text": view.StatusText, "conflict_code": view.ConflictCode,
			"conflict_reason": view.ConflictReason,
		})
	}
	c.JSON(200, gin.H{"items": out, "total": len(out)})
}

func (d *DashboardController) AuditLogs(c *gin.Context) {
	page, size := utilsPagination(c)
	rows, total, err := d.Query.Audit.List(page, size)
	if err != nil {
		fail(c, err)
		return
	}
	logs := make([]constructors.AuditView, 0, len(rows))
	for _, row := range rows {
		logs = append(logs, constructors.NewAuditView(row))
	}
	c.JSON(200, gin.H{"items": logs, "total": total, "page": page, "page_size": size})
}
