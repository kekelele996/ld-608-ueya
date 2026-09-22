package controllers

import (
	"strconv"

	"groundTurn/src/constructors"
	"groundTurn/src/services"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// BookingController lists bookings and resolves PENDING ones.
type BookingController struct {
	Booking *services.BookingService
	Query   *services.QueryService
}

func NewBookingController(booking *services.BookingService, query *services.QueryService) *BookingController {
	return &BookingController{Booking: booking, Query: query}
}

func (b *BookingController) List(c *gin.Context) {
	page, size := utilsPagination(c)
	rows, total, err := b.Query.Booking.List(page, size, c.Query("status"))
	if err != nil {
		fail(c, err)
		return
	}
	resources, _ := b.Query.AllResources()
	rm := map[uint]resourceEntry{}
	for _, r := range resources {
		rm[r.ID] = resourceEntry{r.ID, r.ResourceCode, r.ResourceType}
	}
	views := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		entry := rm[row.ResourceID]
		view := constructors.NewBookingView(row, nil)
		view.ResourceCode = entry.code
		view.ResourceType = entry.resourceType
		views = append(views, gin.H{
			"id": view.ID, "resource_id": view.ResourceID,
			"resource_code": view.ResourceCode, "resource_type": view.ResourceType,
			"turnaround_id": view.TurnaroundID, "task_id": view.TaskID,
			"start_time": view.StartTime, "end_time": view.EndTime,
			"booking_status": view.BookingStatus, "status_text": view.StatusText,
			"conflict_code": view.ConflictCode, "conflict_reason": view.ConflictReason,
		})
	}
	c.JSON(200, gin.H{"items": views, "total": total, "page": page, "page_size": size})
}

type resourceEntry struct {
	id           uint
	code         string
	resourceType string
}

// Adjust POST /api/bookings/:id/adjust
func (b *BookingController) Adjust(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req constructors.AdjustBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": "VALIDATION_FAILED", "message": err.Error()})
		return
	}
	start, err := utils.ParseTime(req.StartTime)
	if err != nil {
		c.JSON(400, gin.H{"code": "VALIDATION_FAILED", "message": "开始时间格式错误"})
		return
	}
	end, err := utils.ParseTime(req.EndTime)
	if err != nil {
		c.JSON(400, gin.H{"code": "VALIDATION_FAILED", "message": "结束时间格式错误"})
		return
	}
	username, role, _ := actor(c)
	booking, conflicts, err := b.Booking.Adjust(uint(id), start, end, username, role)
	if err != nil {
		if apiErr, ok := err.(*utils.APIError); ok {
			c.JSON(apiErr.HTTPCode, gin.H{
				"code": apiErr.Code, "message": apiErr.Message,
				"details":   conflicts,
				"booking_id": func() uint {
					if booking != nil {
						return booking.ID
					}
					return 0
				}(),
			})
			return
		}
		fail(c, err)
		return
	}
	c.JSON(200, gin.H{"saved": true, "booking_id": booking.ID, "conflicts": []struct{}{}})
}

// Release POST /api/bookings/:id/release
func (b *BookingController) Release(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	username, role, _ := actor(c)
	booking, err := b.Booking.Release(uint(id), username, role)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, gin.H{"id": booking.ID, "booking_status": booking.BookingStatus})
}
