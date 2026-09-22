package controllers

import (
	"strconv"
	"time"

	"groundTurn/src/constructors"
	"groundTurn/src/models"
	"groundTurn/src/services"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// ResourceController lists resources and updates availability / maintenance.
type ResourceController struct {
	Resource *services.ResourceService
	Query    *services.QueryService
}

func NewResourceController(resource *services.ResourceService, query *services.QueryService) *ResourceController {
	return &ResourceController{Resource: resource, Query: query}
}

func (r *ResourceController) List(c *gin.Context) {
	page, size := utilsPagination(c)
	rows, total, err := r.Query.Resource.List(page, size, c.Query("status"), c.Query("type"))
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, gin.H{"items": constructors.NewResourceViews(rows), "total": total, "page": page, "page_size": size})
}

func (r *ResourceController) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req constructors.ResourceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": "VALIDATION_FAILED", "message": err.Error()})
		return
	}
	var start, end *time.Time
	if req.MaintenanceStart != "" {
		if t, err := utils.ParseTime(req.MaintenanceStart); err == nil {
			start = &t
		}
	}
	if req.MaintenanceEnd != "" {
		if t, err := utils.ParseTime(req.MaintenanceEnd); err == nil {
			end = &t
		}
	}
	username, role, _ := actor(c)
	resource, err := r.Resource.UpdateStatus(uint(id), req.Status, start, end, username, role)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, constructors.NewResourceView(*resource))
}

// Calendar returns bookings inside ?from=&to= for ResourceCalendar.
func (r *ResourceController) Calendar(c *gin.Context) {
	from, err1 := utils.ParseTime(c.Query("from"))
	to, err2 := utils.ParseTime(c.Query("to"))
	if err1 != nil || err2 != nil {
		now := time.Now()
		from = now.Add(-2 * time.Hour)
		to = now.Add(12 * time.Hour)
	}
	bookings, err := r.Query.CalendarWindow(from, to)
	if err != nil {
		fail(c, err)
		return
	}
	resources, _ := r.Query.AllResources()
	resourceMap := map[uint]models.GroundResource{}
	for _, resource := range resources {
		resourceMap[resource.ID] = resource
	}
	c.JSON(200, gin.H{
		"items":     constructors.NewBookingViews(bookings, resourceMap),
		"resources": constructors.NewResourceViews(resources),
	})
}
