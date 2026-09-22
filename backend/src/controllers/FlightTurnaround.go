package controllers

import (
	"strconv"

	"groundTurn/src/constructors"
	"groundTurn/src/services"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// TurnaroundController covers registration, list, detail and plan generation.
type TurnaroundController struct {
	Plan  *services.PlanService
	Query *services.QueryService
}

func NewTurnaroundController(plan *services.PlanService, query *services.QueryService) *TurnaroundController {
	return &TurnaroundController{Plan: plan, Query: query}
}

func (t *TurnaroundController) Register(c *gin.Context) {
	var req constructors.RegisterTurnaroundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": "VALIDATION_FAILED", "message": err.Error()})
		return
	}
	arrival, err := utils.ParseTime(req.ArrivalTime)
	if err != nil {
		c.JSON(400, gin.H{"code": "VALIDATION_FAILED", "message": "到站时间格式错误"})
		return
	}
	departure, err := utils.ParseTime(req.DepartureTime)
	if err != nil {
		c.JSON(400, gin.H{"code": "VALIDATION_FAILED", "message": "离港时间格式错误"})
		return
	}
	username, role, _ := actor(c)
	turn, err := t.Plan.Register(req.FlightNo, req.AircraftReg, req.StandNo, arrival, departure, username, role)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(201, constructors.NewTurnaroundSummary(*turn))
}

func (t *TurnaroundController) List(c *gin.Context) {
	page, size := utils.Pagination(c)
	rows, total, err := t.Query.Turn.List(page, size)
	if err != nil {
		fail(c, err)
		return
	}
	out := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		out = append(out, constructors.NewTurnaroundSummary(row))
	}
	c.JSON(200, gin.H{"items": out, "total": total, "page": page, "page_size": size})
}

func (t *TurnaroundController) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	turn, tasks, bookings, resourceMap, err := t.Plan.GetDetail(uint(id))
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, constructors.NewTurnaroundView(*turn, tasks, bookings, resourceMap))
}

// GeneratePlan POST /api/turnarounds/:id/plan
func (t *TurnaroundController) GeneratePlan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req constructors.GeneratePlanRequest
	_ = c.ShouldBindJSON(&req)
	username, role, _ := actor(c)
	result, err := t.Plan.GeneratePlan(uint(id), req.TaskTypes, username, role)
	if err != nil {
		if apiErr, ok := err.(*utils.APIError); ok {
			// Return the rejected (unsaved) preview alongside the 409.
			status := apiErr.HTTPCode
			c.JSON(status, gin.H{
				"code":    apiErr.Code,
				"message": apiErr.Message,
				"details": apiErr.Details,
				"result":  result,
			})
			return
		}
		fail(c, err)
		return
	}
	c.JSON(201, result)
}
