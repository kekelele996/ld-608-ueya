package controllers

import (
	"strconv"

	"groundTurn/src/constructors"
	"groundTurn/src/services"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// DelayController registers delay minutes and lists/resolves delay events.
type DelayController struct {
	Delay *services.DelayService
	Query *services.QueryService
}

func NewDelayController(delay *services.DelayService, query *services.QueryService) *DelayController {
	return &DelayController{Delay: delay, Query: query}
}

// Register POST /api/turnarounds/:id/delays
func (d *DelayController) Register(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req constructors.RegisterDelayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": "VALIDATION_FAILED", "message": err.Error()})
		return
	}
	username, role, _ := actor(c)
	result, err := d.Delay.RegisterDelay(uint(id), req.DelayType, req.Minutes,
		req.RootCause, req.ResponsibilityTeam, username, role)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(201, result)
}

func (d *DelayController) List(c *gin.Context) {
	page, size := utilsPagination(c)
	rows, total, err := d.Query.Delay.List(page, size)
	if err != nil {
		fail(c, err)
		return
	}
	turns, _, _ := d.Query.Turn.List(1, 200)
	flightMap := map[uint]string{}
	for _, t := range turns {
		flightMap[t.ID] = t.FlightNo
	}
	c.JSON(200, gin.H{
		"items": constructors.NewDelayViews(rows, flightMap),
		"total": total, "page": page, "page_size": size,
	})
}

func (d *DelayController) Resolve(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	username, role, _ := actor(c)
	event, err := d.Delay.Resolve(uint(id), username, role)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, gin.H{"id": event.ID, "resolved_at": utils.FormatTimePtr(event.ResolvedAt)})
}
