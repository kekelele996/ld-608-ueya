package controllers

import (
	"strconv"

	"groundTurn/src/constructors"
	"groundTurn/src/services"

	"github.com/gin-gonic/gin"
)

// TaskController lists tasks and drives accept/complete/block actions.
type TaskController struct {
	Task  *services.TaskService
	Query *services.QueryService
}

func NewTaskController(task *services.TaskService, query *services.QueryService) *TaskController {
	return &TaskController{Task: task, Query: query}
}

func (t *TaskController) List(c *gin.Context) {
	page, size := utilsPagination(c)
	status := c.Query("status")
	rows, total, err := t.Query.Task.List(page, size, status)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, gin.H{"items": constructors.NewTaskViews(rows), "total": total, "page": page, "page_size": size})
}

func (t *TaskController) Accept(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	username, role, teamID := actor(c)
	task, err := t.Task.Accept(uint(id), username, role, teamID)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, constructors.NewTaskView(*task))
}

func (t *TaskController) Complete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	username, role, _ := actor(c)
	task, err := t.Task.Complete(uint(id), username, role)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, constructors.NewTaskView(*task))
}

func (t *TaskController) Block(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req constructors.BlockTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": "VALIDATION_FAILED", "message": err.Error()})
		return
	}
	username, role, _ := actor(c)
	task, err := t.Task.Block(uint(id), req.BlockerNote, username, role)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, constructors.NewTaskView(*task))
}
